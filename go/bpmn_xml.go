package nativebpm

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math"
)

// XML Schema structures for OMG BPMN 2.0 and BPMN-DI

type XMLDefinitions struct {
	XMLName         xml.Name       `xml:"definitions"`
	XmlnsXsi        string         `xml:"xmlns:xsi,attr"`
	Xmlns           string         `xml:"xmlns,attr"`
	XmlnsBpmndi     string         `xml:"xmlns:bpmndi,attr"`
	XmlnsOmgdc      string         `xml:"xmlns:omgdc,attr"`
	XmlnsOmgdi      string         `xml:"xmlns:omgdi,attr"`
	TargetNamespace string         `xml:"targetNamespace,attr"`
	ID              string         `xml:"id,attr"`
	Process         XMLProcess     `xml:"process"`
	BPMNDiagram     XMLBPMNDiagram `xml:"bpmndi:BPMNDiagram"`
}

type XMLProcess struct {
	ID                 string                 `xml:"id,attr"`
	Name               string                 `xml:"name,attr"`
	IsExecutable       bool                   `xml:"isExecutable,attr"`
	StartEvents        []XMLStartEvent        `xml:"startEvent,omitempty"`
	EndEvents          []XMLEndEvent          `xml:"endEvent,omitempty"`
	ServiceTasks       []XMLServiceTask       `xml:"serviceTask,omitempty"`
	UserTasks          []XMLUserTask          `xml:"userTask,omitempty"`
	ExclusiveGateways  []XMLExclusiveGateway  `xml:"exclusiveGateway,omitempty"`
	ParallelGateways   []XMLParallelGateway   `xml:"parallelGateway,omitempty"`
	EventBasedGateways []XMLEventBasedGateway `xml:"eventBasedGateway,omitempty"`
	SequenceFlows      []XMLSequenceFlow      `xml:"sequenceFlow,omitempty"`
}

type XMLStartEvent struct {
	ID       string   `xml:"id,attr"`
	Name     string   `xml:"name,attr"`
	Outgoing []string `xml:"outgoing,omitempty"`
}

type XMLEndEvent struct {
	ID       string   `xml:"id,attr"`
	Name     string   `xml:"name,attr"`
	Incoming []string `xml:"incoming,omitempty"`
}

type XMLServiceTask struct {
	ID       string   `xml:"id,attr"`
	Name     string   `xml:"name,attr"`
	Topic    string   `xml:"topic,attr,omitempty"`
	Incoming []string `xml:"incoming,omitempty"`
	Outgoing []string `xml:"outgoing,omitempty"`
}

type XMLUserTask struct {
	ID       string   `xml:"id,attr"`
	Name     string   `xml:"name,attr"`
	Incoming []string `xml:"incoming,omitempty"`
	Outgoing []string `xml:"outgoing,omitempty"`
}

type XMLExclusiveGateway struct {
	ID       string   `xml:"id,attr"`
	Name     string   `xml:"name,attr"`
	Incoming []string `xml:"incoming,omitempty"`
	Outgoing []string `xml:"outgoing,omitempty"`
}

type XMLParallelGateway struct {
	ID       string   `xml:"id,attr"`
	Name     string   `xml:"name,attr"`
	Incoming []string `xml:"incoming,omitempty"`
	Outgoing []string `xml:"outgoing,omitempty"`
}

type XMLEventBasedGateway struct {
	ID       string   `xml:"id,attr"`
	Name     string   `xml:"name,attr"`
	Incoming []string `xml:"incoming,omitempty"`
	Outgoing []string `xml:"outgoing,omitempty"`
}

type XMLSequenceFlow struct {
	ID                  string                  `xml:"id,attr"`
	SourceRef           string                  `xml:"sourceRef,attr"`
	TargetRef           string                  `xml:"targetRef,attr"`
	ConditionExpression *XMLConditionExpression `xml:"conditionExpression,omitempty"`
}

type XMLConditionExpression struct {
	Type string `xml:"xsi:type,attr"`
	Text string `xml:",chardata"`
}

// Diagram Interchange (BPMN-DI)

type XMLBPMNDiagram struct {
	ID        string       `xml:"id,attr"`
	BPMNPlane XMLBPMNPlane `xml:"bpmndi:BPMNPlane"`
}

type XMLBPMNPlane struct {
	ID          string          `xml:"id,attr"`
	BPMNElement string          `xml:"bpmnElement,attr"`
	Shapes      []XMLBPMNShape  `xml:"bpmndi:BPMNShape"`
	Edges       []XMLBPMNEdge   `xml:"bpmndi:BPMNEdge"`
}

type XMLBPMNShape struct {
	ID          string    `xml:"id,attr"`
	BPMNElement string    `xml:"bpmnElement,attr"`
	Bounds      XMLBounds `xml:"omgdc:Bounds"`
}

type XMLBounds struct {
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Width  float64 `xml:"width,attr"`
	Height float64 `xml:"height,attr"`
}

type XMLBPMNEdge struct {
	ID          string        `xml:"id,attr"`
	BPMNElement string        `xml:"bpmnElement,attr"`
	Waypoints   []XMLWaypoint `xml:"omgdi:waypoint"`
}

type XMLWaypoint struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
}

// Coords represents computed bounds of a node.
type Coords struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// ToBPMNXML compiles the Workflow into a standard BPMN 2.0 XML with complete BPMNDI layout.
func (w *Workflow) ToBPMNXML() ([]byte, error) {
	if w.err != nil {
		return nil, w.err
	}

	// 1. Collect all nodes and flows, ensuring auto-connected end events
	sourceIDs := make(map[string]bool)
	for _, f := range w.Flows {
		if src, ok := f["source"].(string); ok {
			sourceIDs[src] = true
		}
	}

	nodes := make([]map[string]interface{}, len(w.Nodes))
	copy(nodes, w.Nodes)

	flows := make([]map[string]interface{}, len(w.Flows))
	copy(flows, w.Flows)

	for _, node := range w.Nodes {
		nodeType, _ := node["type"].(string)
		nodeID, _ := node["id"].(string)
		if nodeType == "endEvent" || nodeType == "startEvent" {
			continue
		}
		if !sourceIDs[nodeID] {
			endID := fmt.Sprintf("end_%s", nodeID)
			nodes = append(nodes, map[string]interface{}{
				"type": "endEvent",
				"id":   endID,
				"name": "Process Finished",
			})
			flows = append(flows, map[string]interface{}{
				"id":        fmt.Sprintf("flow-%s-%s", nodeID, endID),
				"source":    nodeID,
				"target":    endID,
				"condition": "",
			})
		}
	}

	// 2. Compute Layout (BFS Level-based auto layout)
	coords := computeWorkflowLayout(nodes, flows)

	// 3. Build XML Process
	proc := XMLProcess{
		ID:           w.ID,
		Name:         w.Name,
		IsExecutable: true,
	}

	getInAndOut := func(id string) ([]string, []string) {
		var in, out []string
		for _, f := range flows {
			fid, _ := f["id"].(string)
			src, _ := f["source"].(string)
			tgt, _ := f["target"].(string)
			if src == id {
				out = append(out, fid)
			}
			if tgt == id {
				in = append(in, fid)
			}
		}
		return in, out
	}

	for _, n := range nodes {
		id, _ := n["id"].(string)
		name, _ := n["name"].(string)
		nType, _ := n["type"].(string)
		in, out := getInAndOut(id)

		switch nType {
		case "startEvent":
			proc.StartEvents = append(proc.StartEvents, XMLStartEvent{
				ID:       id,
				Name:     name,
				Outgoing: out,
			})
		case "endEvent":
			proc.EndEvents = append(proc.EndEvents, XMLEndEvent{
				ID:       id,
				Name:     name,
				Incoming: in,
			})
		case "serviceTask", "aiTask", "aiServiceTask":
			topic, _ := n["topic"].(string)
			proc.ServiceTasks = append(proc.ServiceTasks, XMLServiceTask{
				ID:       id,
				Name:     name,
				Topic:    topic,
				Incoming: in,
				Outgoing: out,
			})
		case "userTask":
			proc.UserTasks = append(proc.UserTasks, XMLUserTask{
				ID:       id,
				Name:     name,
				Incoming: in,
				Outgoing: out,
			})
		case "exclusiveGateway":
			proc.ExclusiveGateways = append(proc.ExclusiveGateways, XMLExclusiveGateway{
				ID:       id,
				Name:     name,
				Incoming: in,
				Outgoing: out,
			})
		case "parallelGateway":
			proc.ParallelGateways = append(proc.ParallelGateways, XMLParallelGateway{
				ID:       id,
				Name:     name,
				Incoming: in,
				Outgoing: out,
			})
		case "eventBasedGateway":
			proc.EventBasedGateways = append(proc.EventBasedGateways, XMLEventBasedGateway{
				ID:       id,
				Name:     name,
				Incoming: in,
				Outgoing: out,
			})
		}
	}

	// 4. Build Sequence Flows
	for _, f := range flows {
		fid, _ := f["id"].(string)
		src, _ := f["source"].(string)
		tgt, _ := f["target"].(string)
		cond, _ := f["condition"].(string)

		var condExpr *XMLConditionExpression
		if cond != "" {
			condExpr = &XMLConditionExpression{
				Type: "xsi:tFormalExpression",
				Text: fmt.Sprintf("${%s}", cond),
			}
		}
		proc.SequenceFlows = append(proc.SequenceFlows, XMLSequenceFlow{
			ID:                  fid,
			SourceRef:           src,
			TargetRef:           tgt,
			ConditionExpression: condExpr,
		})
	}

	// 5. Build Diagram Interchange (Shapes and Edges)
	plane := XMLBPMNPlane{
		ID:          "BPMNPlane_1",
		BPMNElement: w.ID,
	}

	for _, n := range nodes {
		id, _ := n["id"].(string)
		c, ok := coords[id]
		if !ok {
			continue
		}
		plane.Shapes = append(plane.Shapes, XMLBPMNShape{
			ID:          fmt.Sprintf("%s_di", id),
			BPMNElement: id,
			Bounds: XMLBounds{
				X:      c.X,
				Y:      c.Y,
				Width:  c.Width,
				Height: c.Height,
			},
		})
	}

	for _, f := range flows {
		fid, _ := f["id"].(string)
		src, _ := f["source"].(string)
		tgt, _ := f["target"].(string)

		srcCoord, hasSrc := coords[src]
		dstCoord, hasDst := coords[tgt]
		if !hasSrc || !hasDst {
			continue
		}

		isLeftToRight := dstCoord.X > (srcCoord.X + srcCoord.Width + 20.0)
		isRightToLeft := (srcCoord.X - 20.0) > (dstCoord.X + dstCoord.Width)
		isTopToBottom := dstCoord.Y > (srcCoord.Y + srcCoord.Height + 20.0)
		isBottomToTop := (srcCoord.Y - 20.0) > (dstCoord.Y + dstCoord.Height)

		var startWp, endWp XMLWaypoint
		if isLeftToRight {
			startWp = XMLWaypoint{
				X: srcCoord.X + srcCoord.Width,
				Y: srcCoord.Y + srcCoord.Height/2.0,
			}
			endWp = XMLWaypoint{
				X: dstCoord.X,
				Y: dstCoord.Y + dstCoord.Height/2.0,
			}
		} else if isRightToLeft {
			startWp = XMLWaypoint{
				X: srcCoord.X,
				Y: srcCoord.Y + srcCoord.Height/2.0,
			}
			endWp = XMLWaypoint{
				X: dstCoord.X + dstCoord.Width,
				Y: dstCoord.Y + dstCoord.Height/2.0,
			}
		} else if isTopToBottom {
			startWp = XMLWaypoint{
				X: srcCoord.X + srcCoord.Width/2.0,
				Y: srcCoord.Y + srcCoord.Height,
			}
			endWp = XMLWaypoint{
				X: dstCoord.X + dstCoord.Width/2.0,
				Y: dstCoord.Y,
			}
		} else if isBottomToTop {
			startWp = XMLWaypoint{
				X: srcCoord.X + srcCoord.Width/2.0,
				Y: srcCoord.Y,
			}
			endWp = XMLWaypoint{
				X: dstCoord.X + dstCoord.Width/2.0,
				Y: dstCoord.Y + dstCoord.Height,
			}
		} else {
			startWp = XMLWaypoint{
				X: srcCoord.X + srcCoord.Width,
				Y: srcCoord.Y + srcCoord.Height/2.0,
			}
			endWp = XMLWaypoint{
				X: dstCoord.X,
				Y: dstCoord.Y + dstCoord.Height/2.0,
			}
		}

		edge := XMLBPMNEdge{
			ID:          fmt.Sprintf("%s_di", fid),
			BPMNElement: fid,
		}

		if srcCoord.Y != dstCoord.Y {
			midY := (startWp.Y + endWp.Y) / 2.0
			if math.Abs(dstCoord.X-srcCoord.X) > 250.0 {
				maxY := srcCoord.Y
				if dstCoord.Y > maxY {
					maxY = dstCoord.Y
				}
				midY = maxY + srcCoord.Height/2.0 + 90.0
			}
			var midX1, midX2 float64
			if isRightToLeft {
				midX1 = startWp.X - 30.0
				midX2 = endWp.X + 30.0
			} else {
				midX1 = startWp.X + 30.0
				midX2 = endWp.X - 30.0
			}

			if !isLeftToRight && !isRightToLeft {
				edge.Waypoints = []XMLWaypoint{startWp, endWp}
			} else {
				edge.Waypoints = []XMLWaypoint{
					startWp,
					{X: midX1, Y: startWp.Y},
					{X: midX1, Y: midY},
					{X: midX2, Y: midY},
					{X: midX2, Y: endWp.Y},
					endWp,
				}
			}
		} else {
			edge.Waypoints = []XMLWaypoint{startWp, endWp}
		}

		plane.Edges = append(plane.Edges, edge)
	}

	defs := XMLDefinitions{
		XmlnsXsi:        "http://www.w3.org/2001/XMLSchema-instance",
		Xmlns:           "http://www.omg.org/spec/BPMN/20100524/MODEL",
		XmlnsBpmndi:     "http://www.omg.org/spec/BPMN/20100524/DI",
		XmlnsOmgdc:      "http://www.omg.org/spec/DD/20100524/DC",
		XmlnsOmgdi:      "http://www.omg.org/spec/DD/20100524/DI",
		TargetNamespace: "http://bpmn.io/schema/bpmn",
		ID:              "Definitions_1",
		Process:         proc,
		BPMNDiagram: XMLBPMNDiagram{
			ID:        "BPMNDiagram_1",
			BPMNPlane: plane,
		},
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(defs); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// computeWorkflowLayout calculates coordinates using BFS level layout.
func computeWorkflowLayout(nodes []map[string]interface{}, flows []map[string]interface{}) map[string]Coords {
	coords := make(map[string]Coords)
	if len(nodes) == 0 {
		return coords
	}

	outgoing := make(map[string][]string)
	incomingCount := make(map[string]int)
	for _, f := range flows {
		src, _ := f["source"].(string)
		tgt, _ := f["target"].(string)
		if src != "" && tgt != "" {
			outgoing[src] = append(outgoing[src], tgt)
			incomingCount[tgt]++
		}
	}

	// Detect back-edges using DFS
	state := make(map[string]int) // 0: unvisited, 1: visiting, 2: visited
	backEdges := make(map[string]map[string]bool)

	var dfs func(string)
	dfs = func(u string) {
		state[u] = 1
		for _, v := range outgoing[u] {
			if state[v] == 1 {
				if backEdges[u] == nil {
					backEdges[u] = make(map[string]bool)
				}
				backEdges[u][v] = true
			} else if state[v] == 0 {
				dfs(v)
			}
		}
		state[u] = 2
	}

	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != "" && incomingCount[id] == 0 {
			dfs(id)
		}
	}
	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != "" && state[id] == 0 {
			dfs(id)
		}
	}

	// BFS Level assignment
	var queue []string
	nodeLevels := make(map[string]int)
	queued := make(map[string]bool)

	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != "" && incomingCount[id] == 0 {
			queue = append(queue, id)
			queued[id] = true
			nodeLevels[id] = 0
		}
	}

	if len(queue) == 0 && len(nodes) > 0 {
		id, _ := nodes[0]["id"].(string)
		if id != "" {
			queue = append(queue, id)
			queued[id] = true
			nodeLevels[id] = 0
		}
	}

	maxDepth := 0
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		queued[curr] = false

		currLevel := nodeLevels[curr]
		if currLevel > maxDepth {
			maxDepth = currLevel
		}

		for _, child := range outgoing[curr] {
			if backEdges[curr] != nil && backEdges[curr][child] {
				continue
			}
			childLevel := currLevel + 1
			if childLevel > nodeLevels[child] {
				nodeLevels[child] = childLevel
				if !queued[child] {
					queue = append(queue, child)
					queued[child] = true
				}
			}
		}
	}

	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != "" {
			if _, ok := nodeLevels[id]; !ok {
				nodeLevels[id] = maxDepth + 1
			}
		}
	}

	levelGroups := make(map[int][]string)
	for _, n := range nodes {
		id, _ := n["id"].(string)
		if id != "" {
			lvl := nodeLevels[id]
			levelGroups[lvl] = append(levelGroups[lvl], id)
		}
	}

	const maxCols = 5
	for lvl, nIDs := range levelGroups {
		row := lvl / maxCols
		col := lvl % maxCols
		if row%2 != 0 {
			col = (maxCols - 1) - col
		}
		x := 150.0 + float64(col)*240.0
		count := len(nIDs)
		for i, nID := range nIDs {
			y := 200.0 + float64(row)*220.0 + (float64(i)-float64(count-1)/2.0)*150.0

			var w, h float64 = 100.0, 80.0
			for _, n := range nodes {
				if n["id"] == nID {
					nType, _ := n["type"].(string)
					switch nType {
					case "startEvent", "endEvent":
						w, h = 36.0, 36.0
						y = y + (80.0-36.0)/2.0
					case "exclusiveGateway", "parallelGateway", "eventBasedGateway":
						w, h = 50.0, 50.0
						y = y + (80.0-50.0)/2.0
					}
					break
				}
			}

			coords[nID] = Coords{
				X:      x,
				Y:      y,
				Width:  w,
				Height: h,
			}
		}
	}

	return coords
}
