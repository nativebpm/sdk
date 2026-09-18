package nativebpm

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorker_UnixDomainSocket(t *testing.T) {
	sockDir, err := os.MkdirTemp("", "nbpm-worker-uds-*")
	require.NoError(t, err)
	defer os.RemoveAll(sockDir)

	sockPath := filepath.Join(sockDir, "test.sock")

	listener, err := net.Listen("unix", sockPath)
	require.NoError(t, err)
	defer listener.Close()

	// Channel to receive the established Yamux server session
	serverMuxChan := make(chan *yamux.Session, 1)

	// Mock server listening on UDS
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			return
		}

		// 1. Read handshake request
		var req struct {
			WorkerID       string   `json:"worker_id"`
			Topics         []string `json:"topics"`
			MaxConcurrency int      `json:"max_concurrency"`
		}
		if err := json.NewDecoder(conn).Decode(&req); err != nil {
			conn.Close()
			return
		}

		// 2. Write handshake response
		resp := map[string]string{"status": "ok", "session_id": "sess-uds-1"}
		if err := json.NewEncoder(conn).Encode(resp); err != nil {
			conn.Close()
			return
		}

		// 3. Upgrade to Yamux Server
		smux, err := yamux.Server(conn, nil)
		if err != nil {
			conn.Close()
			return
		}
		serverMuxChan <- smux
	}()

	// Worker client
	worker := NewWorker("unix://"+sockPath, "token-123").
		WithWorkerID("w-calc-1").
		WithTopic("calculate_tax", func(ctx context.Context, task *TaskContext) (map[string]interface{}, error) {
			amount := task.Variables["amount"].(float64)
			tax := amount * 0.20
			return map[string]interface{}{
				"tax":   tax,
				"total": amount + tax,
			}, nil
		})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = worker.Start(ctx)
	}()
	defer worker.Close()

	// Await server session
	var serverMux *yamux.Session
	select {
	case serverMux = <-serverMuxChan:
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for server Yamux session")
	}
	defer serverMux.Close()

	// Engine opens virtual stream down to worker
	stream, err := serverMux.OpenStream()
	require.NoError(t, err)
	defer stream.Close()

	taskMsg := map[string]interface{}{
		"id":        "task-101",
		"topic":     "calculate_tax",
		"variables": map[string]interface{}{"amount": 100.0},
	}
	err = json.NewEncoder(stream).Encode(taskMsg)
	require.NoError(t, err)

	var result struct {
		TaskID    string                 `json:"task_id"`
		Status    string                 `json:"status"`
		Variables map[string]interface{} `json:"variables"`
		Error     string                 `json:"error"`
	}
	err = json.NewDecoder(stream).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "task-101", result.TaskID)
	assert.Equal(t, "completed", result.Status)
	assert.Equal(t, 20.0, result.Variables["tax"])
	assert.Equal(t, 120.0, result.Variables["total"])
}

func TestWorker_HTTPUpgrade(t *testing.T) {
	serverMuxChan := make(chan *yamux.Session, 1)

	// Mock server supporting HTTP Upgrade: yamux
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/tunnels/connect", r.URL.Path)
		assert.Equal(t, "Upgrade", r.Header.Get("Connection"))
		assert.Equal(t, "yamux", r.Header.Get("Upgrade"))

		hj, ok := w.(http.Hijacker)
		require.True(t, ok)

		conn, bufrw, err := hj.Hijack()
		require.NoError(t, err)

		_, _ = bufrw.WriteString("HTTP/1.1 101 Switching Protocols\r\nConnection: Upgrade\r\nUpgrade: yamux\r\n\r\n")
		_ = bufrw.Flush()

		smux, err := yamux.Server(&bufferedConn{Conn: conn, r: bufrw}, nil)
		require.NoError(t, err)
		serverMuxChan <- smux
	}))
	defer ts.Close()

	worker := NewWorker(ts.URL, "secret-token").
		WithTopic("send_email", func(ctx context.Context, task *TaskContext) (map[string]interface{}, error) {
			to := task.Variables["to"].(string)
			return map[string]interface{}{
				"sent_to": to,
				"success": true,
			}, nil
		})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = worker.Start(ctx)
	}()
	defer worker.Close()

	var serverMux *yamux.Session
	select {
	case serverMux = <-serverMuxChan:
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for server Yamux session")
	}
	defer serverMux.Close()

	stream, err := serverMux.OpenStream()
	require.NoError(t, err)
	defer stream.Close()

	taskMsg := map[string]interface{}{
		"id":        "email-1",
		"topic":     "send_email",
		"variables": map[string]interface{}{"to": "user@example.com"},
	}
	err = json.NewEncoder(stream).Encode(taskMsg)
	require.NoError(t, err)

	var result struct {
		TaskID    string                 `json:"task_id"`
		Status    string                 `json:"status"`
		Variables map[string]interface{} `json:"variables"`
		Error     string                 `json:"error"`
	}
	err = json.NewDecoder(stream).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "email-1", result.TaskID)
	assert.Equal(t, "completed", result.Status)
	assert.Equal(t, "user@example.com", result.Variables["sent_to"])
	assert.Equal(t, true, result.Variables["success"])
}

func TestWorker_HandlerError(t *testing.T) {
	serverMuxChan := make(chan *yamux.Session, 1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		require.True(t, ok)

		conn, bufrw, err := hj.Hijack()
		require.NoError(t, err)

		_, _ = bufrw.WriteString("HTTP/1.1 101 Switching Protocols\r\nConnection: Upgrade\r\nUpgrade: yamux\r\n\r\n")
		_ = bufrw.Flush()

		smux, err := yamux.Server(&bufferedConn{Conn: conn, r: bufrw}, nil)
		require.NoError(t, err)
		serverMuxChan <- smux
	}))
	defer ts.Close()

	worker := NewWorker(ts.URL, "").
		WithTopic("failing_topic", func(ctx context.Context, task *TaskContext) (map[string]interface{}, error) {
			return nil, errors.New("database unavailable")
		})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = worker.Start(ctx)
	}()
	defer worker.Close()

	var serverMux *yamux.Session
	select {
	case serverMux = <-serverMuxChan:
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for server Yamux session")
	}
	defer serverMux.Close()

	stream, err := serverMux.OpenStream()
	require.NoError(t, err)
	defer stream.Close()

	taskMsg := map[string]interface{}{
		"id":    "err-1",
		"topic": "failing_topic",
	}
	err = json.NewEncoder(stream).Encode(taskMsg)
	require.NoError(t, err)

	var res struct {
		TaskID string `json:"task_id"`
		Status string `json:"status"`
		Error  string `json:"error"`
	}
	err = json.NewDecoder(stream).Decode(&res)
	require.NoError(t, err)

	assert.Equal(t, "err-1", res.TaskID)
	assert.Equal(t, "failed", res.Status)
	assert.Contains(t, res.Error, "database unavailable")
}
