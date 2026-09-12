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
	XmlnsBioc       string         `xml:"xmlns:bioc,attr,omitempty"`
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
	BiocStroke  string    `xml:"bioc:stroke,attr,omitempty"`
	BiocFill    string    `xml:"bioc:fill,attr,omitempty"`
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

	// 2. Compute Layout with LayoutOptions and Zero-Overlap Engine
	layoutOpts := w.LayoutOpts
	if layoutOpts.WorkflowID == "" {
		layoutOpts.WorkflowID = w.ID
	}
	coords := computeWorkflowLayout(nodes, flows, layoutOpts)

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
		var stroke, fill string
		if colors, ok := w.NodeColors[id]; ok {
			stroke = colors[0]
			fill = colors[1]
		}
		plane.Shapes = append(plane.Shapes, XMLBPMNShape{
			ID:          fmt.Sprintf("%s_di", id),
			BPMNElement: id,
			BiocStroke:  stroke,
			BiocFill:    fill,
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

	var xmlnsBioc string
	if len(w.NodeColors) > 0 {
		xmlnsBioc = "http://bpmn.io/schema/bpmn/biocolor/1.0"
	}

	defs := XMLDefinitions{
		XmlnsXsi:        "http://www.w3.org/2001/XMLSchema-instance",
		Xmlns:           "http://www.omg.org/spec/BPMN/20100524/MODEL",
		XmlnsBpmndi:     "http://www.omg.org/spec/BPMN/20100524/DI",
		XmlnsOmgdc:      "http://www.omg.org/spec/DD/20100524/DC",
		XmlnsOmgdi:      "http://www.omg.org/spec/DD/20100524/DI",
		XmlnsBioc:       xmlnsBioc,
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

	rawXML := buf.Bytes()
	if w.LayoutOpts.Orientation == OrientationVertical {
		return TransposeBPMNXML(rawXML)
	}

	return rawXML, nil
}

// TransposeBPMNXML transposes a standard horizontal BPMN 2.0 XML diagram into a vertical waterfall layout.
// It swaps shape center positions (Bounds) and re-routes sequence flow waypoints to top/bottom boundary ports.
func TransposeBPMNXML(xmlBytes []byte) ([]byte, error) {
	var defs XMLDefinitions
	if err := xml.Unmarshal(xmlBytes, &defs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal BPMN XML for transposition: %w", err)
	}

	origBounds := make(map[string]XMLBounds)
	for _, shape := range defs.BPMNDiagram.BPMNPlane.Shapes {
		origBounds[shape.BPMNElement] = shape.Bounds
	}

	newBounds := make(map[string]XMLBounds)
	for i := range defs.BPMNDiagram.BPMNPlane.Shapes {
		shape := &defs.BPMNDiagram.BPMNPlane.Shapes[i]
		orig := origBounds[shape.BPMNElement]

		cx := orig.X + orig.Width/2.0
		cy := orig.Y + orig.Height/2.0

		newCx := cy
		newCy := cx

		shape.Bounds.X = newCx - orig.Width/2.0
		shape.Bounds.Y = newCy - orig.Height/2.0
		shape.Bounds.Width = orig.Width
		shape.Bounds.Height = orig.Height

		newBounds[shape.BPMNElement] = shape.Bounds
	}

	flowMap := make(map[string]XMLSequenceFlow)
	for _, f := range defs.Process.SequenceFlows {
		flowMap[f.ID] = f
	}

	getTransposedEndpoint := func(ox, oy float64, orig, newEl XMLBounds) XMLWaypoint {
		distRight := math.Abs(ox - (orig.X + orig.Width))
		distLeft := math.Abs(ox - orig.X)
		distBottom := math.Abs(oy - (orig.Y + orig.Height))
		distTop := math.Abs(oy - orig.Y)

		minDist := math.Min(math.Min(distRight, distLeft), math.Min(distBottom, distTop))

		if minDist == distRight {
			return XMLWaypoint{
				X: newEl.X + newEl.Width/2.0,
				Y: newEl.Y + newEl.Height,
			}
		}
		if minDist == distLeft {
			return XMLWaypoint{
				X: newEl.X + newEl.Width/2.0,
				Y: newEl.Y,
			}
		}
		if minDist == distBottom {
			return XMLWaypoint{
				X: newEl.X + newEl.Width,
				Y: newEl.Y + newEl.Height/2.0,
			}
		}
		return XMLWaypoint{
			X: newEl.X,
			Y: newEl.Y + newEl.Height/2.0,
		}
	}

	for i := range defs.BPMNDiagram.BPMNPlane.Edges {
		edge := &defs.BPMNDiagram.BPMNPlane.Edges[i]
		flow, exists := flowMap[edge.BPMNElement]
		if !exists || len(edge.Waypoints) < 2 {
			for j := range edge.Waypoints {
				wp := &edge.Waypoints[j]
				wp.X, wp.Y = wp.Y, wp.X
			}
			continue
		}

		origSrc, hasSrc := origBounds[flow.SourceRef]
		origDst, hasDst := origBounds[flow.TargetRef]
		newSrc := newBounds[flow.SourceRef]
		newDst := newBounds[flow.TargetRef]

		if !hasSrc || !hasDst {
			for j := range edge.Waypoints {
				wp := &edge.Waypoints[j]
				wp.X, wp.Y = wp.Y, wp.X
			}
			continue
		}

		origWps := make([]XMLWaypoint, len(edge.Waypoints))
		copy(origWps, edge.Waypoints)

		startWp := getTransposedEndpoint(origWps[0].X, origWps[0].Y, origSrc, newSrc)
		endWp := getTransposedEndpoint(origWps[len(origWps)-1].X, origWps[len(origWps)-1].Y, origDst, newDst)

		newWps := []XMLWaypoint{startWp}

		if len(origWps) > 2 {
			midX := (startWp.X + endWp.X) / 2.0
			midY1 := startWp.Y + 30.0
			if origWps[1].X < origWps[0].X {
				midY1 = startWp.Y - 30.0
			}

			midY2 := endWp.Y - 30.0
			if origWps[len(origWps)-2].X > origWps[len(origWps)-1].X {
				midY2 = endWp.Y + 30.0
			}

			isRightToLeft := (newDst.X + newDst.Width) <= newSrc.X
			isLeftToRight := (newSrc.X + newSrc.Width) <= newDst.X

			if isLeftToRight || isRightToLeft {
				newWps = append(newWps,
					XMLWaypoint{X: startWp.X, Y: midY1},
					XMLWaypoint{X: midX, Y: midY1},
					XMLWaypoint{X: midX, Y: midY2},
					XMLWaypoint{X: endWp.X, Y: midY2},
				)
			} else {
				newWps = append(newWps,
					XMLWaypoint{X: startWp.X, Y: (startWp.Y + endWp.Y) / 2.0},
					XMLWaypoint{X: endWp.X, Y: (startWp.Y + endWp.Y) / 2.0},
				)
			}
		}

		newWps = append(newWps, endWp)
		edge.Waypoints = newWps
	}

	var outBuf bytes.Buffer
	outBuf.WriteString(xml.Header)
	outEnc := xml.NewEncoder(&outBuf)
	outEnc.Indent("", "  ")
	if err := outEnc.Encode(defs); err != nil {
		return nil, err
	}

	return outBuf.Bytes(), nil
}

// computeWorkflowLayout calculates coordinates using LayoutOptions and Zero-Overlap collision resolution.
func computeWorkflowLayout(nodes []map[string]interface{}, flows []map[string]interface{}, opts LayoutOptions) map[string]Coords {
	if opts.Preset == "" || opts.Preset == LayoutAuto {
		metrics := AnalyzeGraphComplexity(nodes, flows)
		opts.Preset = AutoSelectLayoutPreset(nodes, flows)
		if opts.CenterHubID == "" && opts.Preset == LayoutCenterHub {
			opts.CenterHubID = metrics.DetectedHubID
		}
		LogAutoSelection(opts.WorkflowID, opts.Preset, metrics)
	}

	if customStrategy, ok := GetLayoutStrategy(opts.Preset); ok && customStrategy != nil {
		return customStrategy(nodes, flows, opts)
	}

	return StrategyTiered(nodes, flows, opts)
}
