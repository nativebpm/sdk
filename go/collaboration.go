package nativebpm

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"math"
)

// Collaboration represents an OMG BPMN 2.0 multi-process collaboration diagram with pools and swimlanes.
type Collaboration struct {
	ID         string
	Name       string
	Pools      []*PoolParticipant
	LayoutOpts LayoutOptions
	NodeColors map[string][2]string
	err        error
}

// PoolParticipant defines metadata and the workflow pipeline for a participant pool.
type PoolParticipant struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	ProcessRef string    `json:"processRef"`
	Workflow   *Workflow `json:"workflow"`
}

// NewCollaboration initializes a new BPMN 2.0 Collaboration builder.
func NewCollaboration(id, name string) *Collaboration {
	return &Collaboration{
		ID:         id,
		Name:       name,
		Pools:      make([]*PoolParticipant, 0),
		LayoutOpts: DefaultLayoutOptions(),
		NodeColors: make(map[string][2]string),
	}
}

// AddPool registers a new pool/participant within the collaboration and returns
// its underlying *Workflow for fluent node and flow construction.
func (c *Collaboration) AddPool(id, name, processRef string) *Workflow {
	wf := NewWorkflow(processRef, name)
	wf.LayoutOpts = c.LayoutOpts
	c.Pools = append(c.Pools, &PoolParticipant{
		ID:         id,
		Name:       name,
		ProcessRef: processRef,
		Workflow:   wf,
	})
	return wf
}

// SetLayoutPreset configures the default layout preset across all pools.
func (c *Collaboration) SetLayoutPreset(preset LayoutPreset) *Collaboration {
	c.LayoutOpts.Preset = preset
	for _, p := range c.Pools {
		p.Workflow.SetLayoutPreset(preset)
	}
	return c
}

// SetOrientation configures layout orientation (horizontal or vertical) across all pools.
func (c *Collaboration) SetOrientation(orientation LayoutOrientation) *Collaboration {
	c.LayoutOpts.Orientation = orientation
	for _, p := range c.Pools {
		p.Workflow.SetOrientation(orientation)
	}
	return c
}

// SetCenterHub designates a central hub node for the LayoutCenterHub preset.
func (c *Collaboration) SetCenterHub(nodeID string) *Collaboration {
	c.LayoutOpts.CenterHubID = nodeID
	for _, p := range c.Pools {
		p.Workflow.SetCenterHub(nodeID)
	}
	return c
}

// SetNodePosition sets explicit pixel coordinates for a node across any pool.
func (c *Collaboration) SetNodePosition(nodeID string, x, y float64) *Collaboration {
	if c.LayoutOpts.CustomCoords == nil {
		c.LayoutOpts.CustomCoords = make(map[string]Coords)
	}
	c.LayoutOpts.CustomCoords[nodeID] = Coords{X: x, Y: y}
	for _, p := range c.Pools {
		p.Workflow.SetNodePosition(nodeID, x, y)
	}
	return c
}

// SetNodeColor sets bioc stroke and fill colors for an element in the collaboration.
func (c *Collaboration) SetNodeColor(nodeID, stroke, fill string) *Collaboration {
	c.NodeColors[nodeID] = [2]string{stroke, fill}
	for _, p := range c.Pools {
		p.Workflow.SetNodeColor(nodeID, stroke, fill)
	}
	return c
}

// ToBPMNXML compiles the entire collaboration into a standard BPMN 2.0 XML document
// containing <bpmn:collaboration>, <bpmn:participant> entries, multiple <bpmn:process> definitions,
// and compliant BPMNDI pool boundaries (swimlanes).
func (c *Collaboration) ToBPMNXML() ([]byte, error) {
	if c.err != nil {
		return nil, c.err
	}
	for _, p := range c.Pools {
		if p.Workflow.err != nil {
			return nil, p.Workflow.err
		}
	}
	if len(c.Pools) == 0 {
		return nil, fmt.Errorf("collaboration %s has no pools defined", c.ID)
	}

	startX := c.LayoutOpts.StartX
	if startX <= 0 {
		startX = 150.0
	}
	currentY := c.LayoutOpts.StartY
	if currentY <= 0 {
		currentY = 120.0
	}
	poolSpacing := 80.0

	type poolLayoutResult struct {
		poolID      string
		bounds      XMLBounds
		nodes       []map[string]interface{}
		flows       []map[string]interface{}
		coords      map[string]Coords
		colors      map[string][2]string
		orientation LayoutOrientation
	}

	var poolResults []poolLayoutResult
	maxPoolWidth := 1200.0

	for _, p := range c.Pools {
		nodes, flows := prepareAutoEndEvents(p.Workflow.Nodes, p.Workflow.Flows)

		poolOpts := p.Workflow.LayoutOpts
		if poolOpts.Preset == "" || poolOpts.Preset == LayoutAuto {
			poolOpts.Preset = c.LayoutOpts.Preset
		}
		if poolOpts.Orientation == "" {
			poolOpts.Orientation = c.LayoutOpts.Orientation
		}
		poolOpts.StartX = startX + 80.0 // internal padding inside the pool swimlane
		poolOpts.StartY = currentY
		if poolOpts.WorkflowID == "" {
			poolOpts.WorkflowID = p.ProcessRef
		}

		coords := computeWorkflowLayout(nodes, flows, poolOpts)

		if poolOpts.Orientation == OrientationVertical {
			for id, coord := range coords {
				cx := coord.X + coord.Width/2.0
				cy := coord.Y + coord.Height/2.0
				newCx := cy
				newCy := cx
				coord.X = newCx - coord.Width/2.0
				coord.Y = newCy - coord.Height/2.0
				coords[id] = coord
			}
		}

		// Calculate bounds of nodes in this pool
		minX := math.MaxFloat64
		maxX := -math.MaxFloat64
		minY := math.MaxFloat64
		maxY := -math.MaxFloat64

		if len(coords) == 0 {
			minX = startX
			maxX = startX + 600.0
			minY = currentY
			maxY = currentY + 120.0
		} else {
			for _, coord := range coords {
				if coord.X < minX {
					minX = coord.X
				}
				if coord.X+coord.Width > maxX {
					maxX = coord.X + coord.Width
				}
				if coord.Y < minY {
					minY = coord.Y
				}
				if coord.Y+coord.Height > maxY {
					maxY = coord.Y + coord.Height
				}
			}
		}

		poolLeft := startX - 50.0
		if minX-80.0 < poolLeft {
			poolLeft = minX - 80.0
		}
		poolTop := minY - 50.0
		poolHeight := (maxY - poolTop) + 50.0
		if poolHeight < 220.0 {
			poolHeight = 220.0
		}
		poolWidth := (maxX - poolLeft) + 100.0
		if poolWidth > maxPoolWidth {
			maxPoolWidth = poolWidth
		}

		bounds := XMLBounds{
			X:      poolLeft,
			Y:      poolTop,
			Width:  poolWidth,
			Height: poolHeight,
		}

		mergedColors := make(map[string][2]string)
		for k, v := range c.NodeColors {
			mergedColors[k] = v
		}
		for k, v := range p.Workflow.NodeColors {
			mergedColors[k] = v
		}

		poolResults = append(poolResults, poolLayoutResult{
			poolID:      p.ID,
			bounds:      bounds,
			nodes:       nodes,
			flows:       flows,
			coords:      coords,
			colors:      mergedColors,
			orientation: poolOpts.Orientation,
		})

		currentY = poolTop + poolHeight + poolSpacing
	}

	// Harmonize all pool widths so swimlanes align evenly
	for i := range poolResults {
		poolResults[i].bounds.Width = maxPoolWidth
	}

	// Build Collaboration Participants & Plane shapes
	collab := &XMLCollaboration{
		ID:   c.ID,
		Name: c.Name,
	}

	var allShapes []XMLBPMNShape
	var allEdges []XMLBPMNEdge
	var processes []XMLProcess

	for i, p := range c.Pools {
		collab.Participants = append(collab.Participants, XMLParticipant{
			ID:         p.ID,
			Name:       p.Name,
			ProcessRef: p.ProcessRef,
		})

		isHoriz := true
		allShapes = append(allShapes, XMLBPMNShape{
			ID:           fmt.Sprintf("%s_di", p.ID),
			BPMNElement:  p.ID,
			IsHorizontal: &isHoriz,
			Bounds:       poolResults[i].bounds,
		})
	}

	// Build Processes and Element Shapes/Edges
	hasColors := len(c.NodeColors) > 0
	for i, p := range c.Pools {
		res := poolResults[i]
		proc := buildXMLProcess(p.ProcessRef, p.Name, res.nodes, res.flows)
		processes = append(processes, proc)

		shapes, edges := buildPlaneShapesAndEdges(res.nodes, res.flows, res.coords, res.colors, res.orientation)
		allShapes = append(allShapes, shapes...)
		allEdges = append(allEdges, edges...)
		if len(res.colors) > 0 {
			hasColors = true
		}
	}

	plane := XMLBPMNPlane{
		ID:          "BPMNPlane_1",
		BPMNElement: c.ID, // References Collaboration ID
		Shapes:      allShapes,
		Edges:       allEdges,
	}

	var xmlnsBioc string
	if hasColors {
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
		Collaboration:   collab,
		Processes:       processes,
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
