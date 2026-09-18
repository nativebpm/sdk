package nativebpm

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/yamux"
)

// TaskContext provides contextual information about the workflow step being executed.
type TaskContext struct {
	ID           string                 `json:"id"`
	Topic        string                 `json:"topic"`
	DefinitionID string                 `json:"definition_id"`
	InstanceID   string                 `json:"instance_id"`
	StepID       string                 `json:"step_id"`
	Variables    map[string]interface{} `json:"variables"`
	TimeoutMs    int64                  `json:"timeout_ms,omitempty"`
}

// TaskHandler is the user-provided function executed for a task topic.
// It receives the TaskContext and returns output variables (or an error).
type TaskHandler func(ctx context.Context, task *TaskContext) (map[string]interface{}, error)

// Worker runs a Yamux reverse-tunnel client over UDS (unix://) or TCP (http://, https://).
type Worker struct {
	serverURL      string
	apiToken       string
	workerID       string
	maxConcurrency int
	handlers       map[string]TaskHandler
	tcpTargets     map[string]string
	mu             sync.RWMutex
	activeMux      *yamux.Session
	closed         int32
	backoffMin     time.Duration
	backoffMax     time.Duration
	sem            chan struct{}
}

// NewWorker initializes a new Worker configured to connect to the specified serverURL.
// Supports:
//   - "unix:///path/to/worker.sock" for zero-TCP UDS IPC (< 0.2ms latency)
//   - "http://localhost:8080" or "https://..." for Cloudflared-style Yamux over TCP
func NewWorker(serverURL, apiToken string) *Worker {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "worker"
	}
	return &Worker{
		serverURL:      serverURL,
		apiToken:       apiToken,
		workerID:       fmt.Sprintf("%s-%d", hostname, time.Now().UnixNano()%100000),
		maxConcurrency: 50,
		handlers:       make(map[string]TaskHandler),
		tcpTargets:     make(map[string]string),
		backoffMin:     50 * time.Millisecond,
		backoffMax:     2 * time.Second,
	}
}

// WithWorkerID assigns a custom unique worker identifier.
func (w *Worker) WithWorkerID(id string) *Worker {
	w.workerID = id
	return w
}

// WithMaxConcurrency configures the maximum concurrent streams handled simultaneously.
func (w *Worker) WithMaxConcurrency(n int) *Worker {
	if n > 0 {
		w.maxConcurrency = n
	}
	return w
}

// WithTopic registers a handler function for a specific BPMN serviceTask topic.
func (w *Worker) WithTopic(topic string, handler TaskHandler) *Worker {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.handlers[topic] = handler
	return w
}

// WithTCPForwarding registers allowed L4 proxy destinations (e.g. "ssh" -> "sshd:22", "rdp" -> "xrdp:3389").
func (w *Worker) WithTCPForwarding(targets map[string]string) *Worker {
	w.mu.Lock()
	defer w.mu.Unlock()
	for k, v := range targets {
		w.tcpTargets[k] = v
	}
	return w
}

// Start initiates the reverse tunnel connection and enters the task execution loop.
// It automatically reconnects with exponential backoff if the connection drops.
func (w *Worker) Start(ctx context.Context) error {
	w.sem = make(chan struct{}, w.maxConcurrency)
	backoff := w.backoffMin

	for {
		if atomic.LoadInt32(&w.closed) == 1 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		mux, err := w.connect(ctx)
		if err != nil {
			if atomic.LoadInt32(&w.closed) == 1 {
				return nil
			}
			slog.Debug("Tunnel worker connection failed, retrying...", "error", err, "backoff", backoff)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff *= 2
				if backoff > w.backoffMax {
					backoff = w.backoffMax
				}
				continue
			}
		}

		backoff = w.backoffMin
		w.mu.Lock()
		w.activeMux = mux
		w.mu.Unlock()

		slog.Info("Tunnel worker connected and ready for tasks", "worker_id", w.workerID, "topics", w.getTopicList())

		// Accept loop: wait for incoming task streams dispatched by NativeBPM
		w.serveMux(ctx, mux)

		w.mu.Lock()
		w.activeMux = nil
		w.mu.Unlock()

		if atomic.LoadInt32(&w.closed) == 1 {
			return nil
		}
	}
}

// Close gracefully disconnects the worker and terminates active streams.
func (w *Worker) Close() error {
	atomic.StoreInt32(&w.closed, 1)
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.activeMux != nil {
		return w.activeMux.Close()
	}
	return nil
}

func (w *Worker) getTopicList() []string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	topics := make([]string, 0, len(w.handlers))
	for t := range w.handlers {
		topics = append(topics, t)
	}
	return topics
}

func (w *Worker) connect(ctx context.Context) (*yamux.Session, error) {
	topics := w.getTopicList()

	// 1. Unix Domain Socket connection
	if strings.HasPrefix(w.serverURL, "unix://") || strings.HasPrefix(w.serverURL, "/") {
		sockPath := strings.TrimPrefix(w.serverURL, "unix://")
		var d net.Dialer
		conn, err := d.DialContext(ctx, "unix", sockPath)
		if err != nil {
			return nil, fmt.Errorf("failed to dial unix domain socket %s: %w", sockPath, err)
		}

		// Handshake
		req := struct {
			WorkerID       string   `json:"worker_id"`
			Topics         []string `json:"topics"`
			MaxConcurrency int      `json:"max_concurrency"`
		}{
			WorkerID:       w.workerID,
			Topics:         topics,
			MaxConcurrency: w.maxConcurrency,
		}

		if err := json.NewEncoder(conn).Encode(req); err != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("handshake request failed: %w", err)
		}

		var resp struct {
			Status    string `json:"status"`
			SessionID string `json:"session_id,omitempty"`
			Error     string `json:"error,omitempty"`
		}
		if err := json.NewDecoder(conn).Decode(&resp); err != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("handshake response decode failed: %w", err)
		}

		if resp.Status != "ok" {
			_ = conn.Close()
			return nil, fmt.Errorf("handshake rejected: %s", resp.Error)
		}

		return yamux.Client(conn, nil)
	}

	// 2. HTTP Upgrade connection over TCP/TLS
	u, err := url.Parse(w.serverURL)
	if err != nil {
		return nil, fmt.Errorf("invalid server URL %s: %w", w.serverURL, err)
	}

	host := u.Host
	if !strings.Contains(host, ":") {
		if u.Scheme == "https" {
			host += ":443"
		} else {
			host += ":80"
		}
	}

	var conn net.Conn
	var d net.Dialer
	if u.Scheme == "https" {
		tlsConfig := &tls.Config{
			ServerName: u.Hostname(),
		}
		conn, err = tls.DialWithDialer(&d, "tcp", host, tlsConfig)
	} else {
		conn, err = d.DialContext(ctx, "tcp", host)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", host, err)
	}

	// Build query parameters
	q := url.Values{}
	q.Set("worker_id", w.workerID)
	q.Set("topics", strings.Join(topics, ","))
	q.Set("concurrency", strconv.Itoa(w.maxConcurrency))

	targetPath := "/api/tunnels/connect"
	if u.Path != "" && u.Path != "/" {
		targetPath = strings.TrimSuffix(u.Path, "/") + "/api/tunnels/connect"
	}

	reqStr := fmt.Sprintf("GET %s?%s HTTP/1.1\r\n"+
		"Host: %s\r\n"+
		"Connection: Upgrade\r\n"+
		"Upgrade: yamux\r\n", targetPath, q.Encode(), u.Host)

	if w.apiToken != "" {
		reqStr += fmt.Sprintf("Authorization: Bearer %s\r\n", w.apiToken)
	}
	reqStr += "\r\n"

	if _, err := conn.Write([]byte(reqStr)); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to send HTTP upgrade request: %w", err)
	}

	bufr := bufio.NewReader(conn)
	resp, err := http.ReadResponse(bufr, nil)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to read HTTP upgrade response: %w", err)
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		_ = conn.Close()
		return nil, fmt.Errorf("unexpected HTTP upgrade status: %d %s", resp.StatusCode, resp.Status)
	}

	return yamux.Client(&bufferedConn{Conn: conn, r: bufr}, nil)
}

func (w *Worker) serveMux(ctx context.Context, mux *yamux.Session) {
	for {
		if atomic.LoadInt32(&w.closed) == 1 {
			return
		}
		stream, err := mux.AcceptStream()
		if err != nil {
			return
		}

		w.sem <- struct{}{}
		go func(s net.Conn) {
			defer func() {
				<-w.sem
				_ = s.Close()
			}()
			w.handleStream(ctx, s)
		}(stream)
	}
}

func (w *Worker) handleStream(ctx context.Context, stream net.Conn) {
	var head [1]byte
	if _, err := io.ReadFull(stream, head[:]); err != nil {
		return
	}

	// 1. Raw L4 TCP Stream Proxying (SSH, RDP, DB)
	if head[0] == 0x02 {
		w.handleRawTCPStream(ctx, stream)
		return
	}

	// 2. BPMN ServiceTask RPC (JSON payload)
	var reader io.Reader = stream
	if head[0] == '{' {
		reader = io.MultiReader(bytes.NewReader(head[:]), stream)
	}

	var msg TaskContext
	if err := json.NewDecoder(reader).Decode(&msg); err != nil {
		return
	}

	w.mu.RLock()
	handler, ok := w.handlers[msg.Topic]
	w.mu.RUnlock()

	type taskResult struct {
		TaskID    string                 `json:"task_id"`
		Status    string                 `json:"status"`
		Variables map[string]interface{} `json:"variables,omitempty"`
		Error     string                 `json:"error,omitempty"`
	}

	res := taskResult{
		TaskID: msg.ID,
		Status: "completed",
	}

	if !ok {
		res.Status = "failed"
		res.Error = fmt.Sprintf("no handler registered for topic: %s", msg.Topic)
	} else {
		taskCtx := ctx
		if msg.TimeoutMs > 0 {
			var cancel context.CancelFunc
			taskCtx, cancel = context.WithTimeout(ctx, time.Duration(msg.TimeoutMs)*time.Millisecond)
			defer cancel()
		}

		outVars, err := handler(taskCtx, &msg)
		if err != nil {
			res.Status = "failed"
			res.Error = err.Error()
		} else {
			res.Variables = outVars
		}
	}

	_ = json.NewEncoder(stream).Encode(res)
}

func (w *Worker) handleRawTCPStream(ctx context.Context, stream net.Conn) {
	var lenBuf [2]byte
	if _, err := io.ReadFull(stream, lenBuf[:]); err != nil {
		_ = writeRawResponse(stream, 0x01, "failed to read target length")
		return
	}
	targetLen := binary.BigEndian.Uint16(lenBuf[:])
	targetBuf := make([]byte, targetLen)
	if _, err := io.ReadFull(stream, targetBuf); err != nil {
		_ = writeRawResponse(stream, 0x01, "failed to read target")
		return
	}
	target := string(targetBuf)

	w.mu.RLock()
	dialAddr, allowed := w.tcpTargets[target]
	if !allowed && strings.Contains(target, ":") {
		if _, wildcard := w.tcpTargets["*"]; wildcard {
			dialAddr = target
			allowed = true
		}
	}
	w.mu.RUnlock()

	if !allowed || dialAddr == "" {
		_ = writeRawResponse(stream, 0x01, fmt.Sprintf("target %q not allowed by worker policy", target))
		return
	}

	targetConn, err := net.DialTimeout("tcp", dialAddr, 5*time.Second)
	if err != nil {
		_ = writeRawResponse(stream, 0x01, fmt.Sprintf("dial failed: %v", err))
		return
	}
	defer targetConn.Close()

	// Acknowledge connection success (status 0x00, empty error)
	if err := writeRawResponse(stream, 0x00, ""); err != nil {
		return
	}

	slog.Info("Worker established L4 proxy bridge", "target", target, "dialAddr", dialAddr)

	errChan := make(chan error, 2)
	go func() {
		_, err := io.Copy(targetConn, stream)
		errChan <- err
	}()
	go func() {
		_, err := io.Copy(stream, targetConn)
		errChan <- err
	}()

	<-errChan
}

func writeRawResponse(w io.Writer, status byte, errMsg string) error {
	errBytes := []byte(errMsg)
	buf := make([]byte, 1+2+len(errBytes))
	buf[0] = status
	binary.BigEndian.PutUint16(buf[1:3], uint16(len(errBytes)))
	copy(buf[3:], errBytes)
	_, err := w.Write(buf)
	return err
}

type bufferedConn struct {
	net.Conn
	r io.Reader
}

func (bc *bufferedConn) Read(p []byte) (int, error) {
	return bc.r.Read(p)
}
