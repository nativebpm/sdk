package nativebpm

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// XML Schema structures for OMG BPMN 2.0 and BPMN-DI

type XMLDefinitions struct {
	XMLName         xml.Name          `xml:"definitions"`
	XmlnsXsi        string            `xml:"xmlns:xsi,attr"`
	Xmlns           string            `xml:"xmlns,attr"`
	XmlnsBpmndi     string            `xml:"xmlns:bpmndi,attr"`
	XmlnsOmgdc      string            `xml:"xmlns:omgdc,attr"`
	XmlnsOmgdi      string            `xml:"xmlns:omgdi,attr"`
	XmlnsBioc       string            `xml:"xmlns:bioc,attr,omitempty"`
	TargetNamespace string            `xml:"targetNamespace,attr"`
	ID              string            `xml:"id,attr"`
	Collaboration   *XMLCollaboration `xml:"collaboration,omitempty"`
	Processes       []XMLProcess      `xml:"process"`
	BPMNDiagram     XMLBPMNDiagram    `xml:"bpmndi:BPMNDiagram"`
}

type XMLCollaboration struct {
	ID           string           `xml:"id,attr"`
	Name         string           `xml:"name,attr,omitempty"`
	Participants []XMLParticipant `xml:"participant"`
}

type XMLParticipant struct {
	ID         string `xml:"id,attr"`
	Name       string `xml:"name,attr"`
	ProcessRef string `xml:"processRef,attr"`
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
	ID           string    `xml:"id,attr"`
	BPMNElement  string    `xml:"bpmnElement,attr"`
	IsHorizontal *bool     `xml:"isHorizontal,attr,omitempty"`
	BiocStroke   string    `xml:"bioc:stroke,attr,omitempty"`
	BiocFill     string    `xml:"bioc:fill,attr,omitempty"`
	Bounds       XMLBounds `xml:"omgdc:Bounds"`
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

// prepareAutoEndEvents ensures all non-event sink nodes have an automatic end event.
func prepareAutoEndEvents(inputNodes []map[string]interface{}, inputFlows []map[string]interface{}) ([]map[string]interface{}, []map[string]interface{}) {
	sourceIDs := make(map[string]bool)
	for _, f := range inputFlows {
		if src, ok := f["source"].(string); ok {
			sourceIDs[src] = true
		}
	}

	nodes := make([]map[string]interface{}, len(inputNodes))
	copy(nodes, inputNodes)

	flows := make([]map[string]interface{}, len(inputFlows))
	copy(flows, inputFlows)

	for _, node := range inputNodes {
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
	return nodes, flows
}

// buildXMLProcess compiles nodes and flows into an XMLProcess struct.
func buildXMLProcess(id, name string, nodes []map[string]interface{}, flows []map[string]interface{}) XMLProcess {
	proc := XMLProcess{
		ID:           id,
		Name:         name,
		IsExecutable: true,
	}

	getInAndOut := func(nodeID string) ([]string, []string) {
		var in, out []string
		for _, f := range flows {
			fid, _ := f["id"].(string)
			src, _ := f["source"].(string)
			tgt, _ := f["target"].(string)
			if src == nodeID {
				out = append(out, fid)
			}
			if tgt == nodeID {
				in = append(in, fid)
			}
		}
		return in, out
	}

	for _, n := range nodes {
		nodeID, _ := n["id"].(string)
		nodeName, _ := n["name"].(string)
		nType, _ := n["type"].(string)
		in, out := getInAndOut(nodeID)

		switch nType {
		case "startEvent":
			proc.StartEvents = append(proc.StartEvents, XMLStartEvent{
				ID:       nodeID,
				Name:     nodeName,
				Outgoing: out,
			})
		case "endEvent":
			proc.EndEvents = append(proc.EndEvents, XMLEndEvent{
				ID:       nodeID,
				Name:     nodeName,
				Incoming: in,
			})
		case "serviceTask", "aiTask", "aiServiceTask":
			topic, _ := n["topic"].(string)
			proc.ServiceTasks = append(proc.ServiceTasks, XMLServiceTask{
				ID:       nodeID,
				Name:     nodeName,
				Topic:    topic,
				Incoming: in,
				Outgoing: out,
			})
		case "userTask":
			proc.UserTasks = append(proc.UserTasks, XMLUserTask{
				ID:       nodeID,
				Name:     nodeName,
				Incoming: in,
				Outgoing: out,
			})
		case "exclusiveGateway":
			proc.ExclusiveGateways = append(proc.ExclusiveGateways, XMLExclusiveGateway{
				ID:       nodeID,
				Name:     nodeName,
				Incoming: in,
				Outgoing: out,
			})
		case "parallelGateway":
			proc.ParallelGateways = append(proc.ParallelGateways, XMLParallelGateway{
				ID:       nodeID,
				Name:     nodeName,
				Incoming: in,
				Outgoing: out,
			})
		case "eventBasedGateway":
			proc.EventBasedGateways = append(proc.EventBasedGateways, XMLEventBasedGateway{
				ID:       nodeID,
				Name:     nodeName,
				Incoming: in,
				Outgoing: out,
			})
		}
	}

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

	return proc
}

// buildPlaneShapesAndEdges generates BPMNDI shapes and sequence flow edges.
func buildPlaneShapesAndEdges(nodes []map[string]interface{}, flows []map[string]interface{}, coords map[string]Coords, nodeColors map[string][2]string, orientation LayoutOrientation) ([]XMLBPMNShape, []XMLBPMNEdge) {
	var shapes []XMLBPMNShape
	var edges []XMLBPMNEdge

	for _, n := range nodes {
		id, _ := n["id"].(string)
		c, ok := coords[id]
		if !ok {
			continue
		}
		var stroke, fill string
		if colors, ok := nodeColors[id]; ok {
			stroke = colors[0]
			fill = colors[1]
		}
		shapes = append(shapes, XMLBPMNShape{
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

		if orientation == OrientationVertical {
			if math.Abs(startWp.X-endWp.X) > 5.0 {
				midY := (startWp.Y + endWp.Y) / 2.0
				edge.Waypoints = []XMLWaypoint{
					startWp,
					{X: startWp.X, Y: midY},
					{X: endWp.X, Y: midY},
					endWp,
				}
			} else {
				edge.Waypoints = []XMLWaypoint{startWp, endWp}
			}
		} else if srcCoord.Y != dstCoord.Y {
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

		edges = append(edges, edge)
	}

	return shapes, edges
}

// ToBPMNXML compiles the Workflow into a standard BPMN 2.0 XML with complete BPMNDI layout.
func (w *Workflow) ToBPMNXML() ([]byte, error) {
	if w.err != nil {
		return nil, w.err
	}

	nodes, flows := prepareAutoEndEvents(w.Nodes, w.Flows)

	layoutOpts := w.LayoutOpts
	if layoutOpts.WorkflowID == "" {
		layoutOpts.WorkflowID = w.ID
	}
	coords := computeWorkflowLayout(nodes, flows, layoutOpts)

	if layoutOpts.Orientation == OrientationVertical {
		for id, c := range coords {
			cx := c.X + c.Width/2.0
			cy := c.Y + c.Height/2.0
			newCx := cy
			newCy := cx
			c.X = newCx - c.Width/2.0
			c.Y = newCy - c.Height/2.0
			coords[id] = c
		}
	}

	proc := buildXMLProcess(w.ID, w.Name, nodes, flows)
	shapes, edges := buildPlaneShapesAndEdges(nodes, flows, coords, w.NodeColors, layoutOpts.Orientation)

	plane := XMLBPMNPlane{
		ID:          "BPMNPlane_1",
		BPMNElement: w.ID,
		Shapes:      shapes,
		Edges:       edges,
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
		Processes:       []XMLProcess{proc},
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

// TransposeBPMNXML transposes a standard horizontal BPMN 2.0 XML diagram into a vertical waterfall layout.
// It swaps shape center positions (Bounds) and re-routes sequence flow waypoints to top/bottom boundary ports.
func TransposeBPMNXML(xmlBytes []byte) ([]byte, error) {
	xmlStr := string(xmlBytes)

	// 1. Map sequence flows
	reFlow := regexp.MustCompile(`<[^:]*:?sequenceFlow\s+([^>]+)>`)
	reAttr := regexp.MustCompile(`([a-zA-Z0-9_:]+)="([^"]*)"`)

	type FlowInfo struct {
		SourceRef string
		TargetRef string
	}
	flows := make(map[string]FlowInfo)
	for _, match := range reFlow.FindAllStringSubmatch(xmlStr, -1) {
		attrs := make(map[string]string)
		for _, attr := range reAttr.FindAllStringSubmatch(match[1], -1) {
			attrs[attr[1]] = attr[2]
		}
		if id, ok := attrs["id"]; ok {
			flows[id] = FlowInfo{
				SourceRef: attrs["sourceRef"],
				TargetRef: attrs["targetRef"],
			}
		}
	}

	// 2. Map and transpose shape bounds
	type BoundsInfo struct {
		X, Y, W, H float64
	}
	origBounds := make(map[string]BoundsInfo)
	newBounds := make(map[string]BoundsInfo)

	reShape := regexp.MustCompile(`(?s)<([^:]*:?BPMNShape)\s+([^>]*bpmnElement="([^"]+)"[^>]*)>(.*?)</[^:]*:?BPMNShape>`)
	reBounds := regexp.MustCompile(`(?s)<([^:]*:?Bounds)\s+([^>]+)(/?>)`)

	for _, match := range reShape.FindAllStringSubmatch(xmlStr, -1) {
		bpmnElem := match[3]
		inner := match[4]
		bMatch := reBounds.FindStringSubmatch(inner)
		if len(bMatch) > 2 {
			attrs := make(map[string]string)
			for _, attr := range reAttr.FindAllStringSubmatch(bMatch[2], -1) {
				attrs[attr[1]] = attr[2]
			}
			x, _ := strconv.ParseFloat(attrs["x"], 64)
			y, _ := strconv.ParseFloat(attrs["y"], 64)
			w, _ := strconv.ParseFloat(attrs["width"], 64)
			h, _ := strconv.ParseFloat(attrs["height"], 64)

			orig := BoundsInfo{X: x, Y: y, W: w, H: h}
			origBounds[bpmnElem] = orig

			cx := x + w/2.0
			cy := y + h/2.0
			newCx := cy
			newCy := cx

			newX := newCx - w/2.0
			newY := newCy - h/2.0
			newBounds[bpmnElem] = BoundsInfo{X: newX, Y: newY, W: w, H: h}
		}
	}

	// Replace all shape bounds in XML
	xmlStr = reShape.ReplaceAllStringFunc(xmlStr, func(shapeXML string) string {
		match := reShape.FindStringSubmatch(shapeXML)
		if len(match) < 4 {
			return shapeXML
		}
		bpmnElem := match[3]
		nb, exists := newBounds[bpmnElem]
		if !exists {
			return shapeXML
		}
		return reBounds.ReplaceAllStringFunc(shapeXML, func(bXML string) string {
			reX := regexp.MustCompile(`\bx="[^"]*"`)
			reY := regexp.MustCompile(`\by="[^"]*"`)
			res := reX.ReplaceAllString(bXML, fmt.Sprintf(`x="%.0f"`, nb.X))
			res = reY.ReplaceAllString(res, fmt.Sprintf(`y="%.0f"`, nb.Y))
			return res
		})
	})

	// 3. Helper for transposed endpoints
	getTransposedEndpoint := func(ox, oy float64, orig, newEl BoundsInfo) (float64, float64) {
		distRight := math.Abs(ox - (orig.X + orig.W))
		distLeft := math.Abs(ox - orig.X)
		distBottom := math.Abs(oy - (orig.Y + orig.H))
		distTop := math.Abs(oy - orig.Y)
		minDist := math.Min(math.Min(distRight, distLeft), math.Min(distBottom, distTop))

		if minDist == distRight { // was right -> becomes bottom
			return newEl.X + newEl.W/2.0, newEl.Y + newEl.H
		}
		if minDist == distLeft { // was left -> becomes top
			return newEl.X + newEl.W/2.0, newEl.Y
		}
		if minDist == distBottom { // was bottom -> becomes right
			return newEl.X + newEl.W, newEl.Y + newEl.H/2.0
		}
		return newEl.X, newEl.Y + newEl.H/2.0
	}

	// 4. Replace edge waypoints
	reEdge := regexp.MustCompile(`(?s)<([^:]*:?BPMNEdge)\s+([^>]*bpmnElement="([^"]+)"[^>]*)>(.*?)</[^:]*:?BPMNEdge>`)
	reWaypoint := regexp.MustCompile(`<([^:]*:?waypoint)\s+([^>]+)(/?>)`)

	xmlStr = reEdge.ReplaceAllStringFunc(xmlStr, func(edgeXML string) string {
		match := reEdge.FindStringSubmatch(edgeXML)
		if len(match) < 4 {
			return edgeXML
		}
		tagName := match[1]
		tagAttrs := match[2]
		flowID := match[3]
		inner := match[4]

		flow, hasFlow := flows[flowID]
		wpMatches := reWaypoint.FindAllStringSubmatch(inner, -1)
		if !hasFlow || len(wpMatches) < 2 {
			newInner := reWaypoint.ReplaceAllStringFunc(inner, func(wpXML string) string {
				wMatch := reWaypoint.FindStringSubmatch(wpXML)
				attrs := make(map[string]string)
				for _, attr := range reAttr.FindAllStringSubmatch(wMatch[2], -1) {
					attrs[attr[1]] = attr[2]
				}
				x := attrs["x"]
				y := attrs["y"]
				reX := regexp.MustCompile(`\bx="[^"]*"`)
				reY := regexp.MustCompile(`\by="[^"]*"`)
				res := reX.ReplaceAllString(wpXML, fmt.Sprintf(`x="%s"`, y))
				res = reY.ReplaceAllString(res, fmt.Sprintf(`y="%s"`, x))
				return res
			})
			return fmt.Sprintf("<%s %s>%s</%s>", tagName, tagAttrs, newInner, tagName)
		}

		origSrc, hasSrc := origBounds[flow.SourceRef]
		origDst, hasDst := origBounds[flow.TargetRef]
		newSrc := newBounds[flow.SourceRef]
		newDst := newBounds[flow.TargetRef]

		if !hasSrc || !hasDst {
			return edgeXML
		}

		firstWpAttrs := make(map[string]string)
		for _, attr := range reAttr.FindAllStringSubmatch(wpMatches[0][2], -1) {
			firstWpAttrs[attr[1]] = attr[2]
		}
		fx, _ := strconv.ParseFloat(firstWpAttrs["x"], 64)
		fy, _ := strconv.ParseFloat(firstWpAttrs["y"], 64)

		lastWpAttrs := make(map[string]string)
		for _, attr := range reAttr.FindAllStringSubmatch(wpMatches[len(wpMatches)-1][2], -1) {
			lastWpAttrs[attr[1]] = attr[2]
		}
		lx, _ := strconv.ParseFloat(lastWpAttrs["x"], 64)
		ly, _ := strconv.ParseFloat(lastWpAttrs["y"], 64)

		startX, startY := getTransposedEndpoint(fx, fy, origSrc, newSrc)
		endX, endY := getTransposedEndpoint(lx, ly, origDst, newDst)

		wpPrefix := "omgdi:waypoint"
		if strings.Contains(wpMatches[0][1], ":") {
			wpPrefix = wpMatches[0][1]
		}

		var newWpsXML strings.Builder
		newWpsXML.WriteString(fmt.Sprintf("\n        <%s x=\"%.0f\" y=\"%.0f\"></%s>", wpPrefix, startX, startY, wpPrefix))

		if math.Abs(startX-endX) > 5.0 {
			midY := (startY + endY) / 2.0
			newWpsXML.WriteString(fmt.Sprintf("\n        <%s x=\"%.0f\" y=\"%.0f\"></%s>", wpPrefix, startX, midY, wpPrefix))
			newWpsXML.WriteString(fmt.Sprintf("\n        <%s x=\"%.0f\" y=\"%.0f\"></%s>", wpPrefix, endX, midY, wpPrefix))
		}
		newWpsXML.WriteString(fmt.Sprintf("\n        <%s x=\"%.0f\" y=\"%.0f\"></%s>\n      ", wpPrefix, endX, endY, wpPrefix))

		return fmt.Sprintf("<%s %s>%s</%s>", tagName, tagAttrs, newWpsXML.String(), tagName)
	})

	return []byte(xmlStr), nil
}

// computeWorkflowLayout calculates coordinates using LayoutOptions, WCC decomposition, and Zero-Overlap collision resolution.
func computeWorkflowLayout(nodes []map[string]interface{}, flows []map[string]interface{}, opts LayoutOptions) map[string]Coords {
	if len(nodes) == 0 {
		return make(map[string]Coords)
	}

	components := decomposeWCC(nodes, flows)
	if len(components) > 1 {
		return computeDisjointWorkflowLayout(components, opts)
	}

	return computeSingleComponentLayout(nodes, flows, opts)
}
