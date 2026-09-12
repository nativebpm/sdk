package nativebpm

import (
	"encoding/json"
	"fmt"
	"strings"
)

type M map[string]interface{}

type LayoutPreset string

const (
	LayoutTiered    LayoutPreset = "tiered"
	LayoutCenterHub LayoutPreset = "center_hub"
	LayoutLinear    LayoutPreset = "linear"
)

type Coords struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
}

type LayoutOrientation string

const (
	OrientationHorizontal LayoutOrientation = "horizontal"
	OrientationVertical   LayoutOrientation = "vertical"
)

type LayoutOptions struct {
	WorkflowID   string            `json:"workflowId,omitempty"`
	Preset       LayoutPreset      `json:"preset,omitempty"`
	Orientation  LayoutOrientation `json:"orientation,omitempty"`
	StartX       float64           `json:"startX,omitempty"`
	StartY       float64           `json:"startY,omitempty"`
	ColSpacing   float64           `json:"colSpacing,omitempty"`
	RowSpacing   float64           `json:"rowSpacing,omitempty"`
	CenterHubID  string            `json:"centerHubId,omitempty"`
	CustomCoords map[string]Coords `json:"customCoords,omitempty"`
	NodeTiers    map[string]int    `json:"nodeTiers,omitempty"`
}

func DefaultLayoutOptions() LayoutOptions {
	return LayoutOptions{
		Preset:       LayoutAuto,
		Orientation:  OrientationHorizontal,
		StartX:       150.0,
		StartY:       200.0,
		ColSpacing:   240.0,
		RowSpacing:   130.0,
		CustomCoords: make(map[string]Coords),
		NodeTiers:    make(map[string]int),
	}
}

type Workflow struct {
	ID            string                   `json:"id"`
	Name          string                   `json:"name"`
	Nodes         []map[string]interface{} `json:"nodes"`
	Flows         []map[string]interface{} `json:"flows"`
	LayoutOpts    LayoutOptions            `json:"layoutOpts,omitempty"`
	NodeColors    map[string][2]string     `json:"nodeColors,omitempty"`
	err           error
	currentNodeID string
	pendingMerges []string
}

func (w *Workflow) MarshalJSON() ([]byte, error) {
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

	type Alias Workflow
	return json.Marshal(&struct {
		Nodes []map[string]interface{} `json:"nodes"`
		Flows []map[string]interface{} `json:"flows"`
		*Alias
	}{
		Nodes: nodes,
		Flows: flows,
		Alias: (*Alias)(w),
	})
}

func NewWorkflow(id, name string) *Workflow {
	w := &Workflow{
		ID:         id,
		Name:       name,
		Nodes:      make([]map[string]interface{}, 0),
		Flows:      make([]map[string]interface{}, 0),
		LayoutOpts: DefaultLayoutOptions(),
		NodeColors: make(map[string][2]string),
	}
	return w
}

// SetLayoutPreset configures the automated layout engine preset.
func (w *Workflow) SetLayoutPreset(preset LayoutPreset) *Workflow {
	w.LayoutOpts.Preset = preset
	return w
}

// SetOrientation configures the layout orientation (horizontal or vertical).
func (w *Workflow) SetOrientation(orientation LayoutOrientation) *Workflow {
	w.LayoutOpts.Orientation = orientation
	return w
}

// SetCenterHub designates a central hub node for the LayoutCenterHub preset.
func (w *Workflow) SetCenterHub(nodeID string) *Workflow {
	w.LayoutOpts.CenterHubID = nodeID
	return w
}

// SetNodePosition sets pixel-perfect coordinate overrides for a specific node.
func (w *Workflow) SetNodePosition(nodeID string, x, y float64) *Workflow {
	if w.LayoutOpts.CustomCoords == nil {
		w.LayoutOpts.CustomCoords = make(map[string]Coords)
	}
	c := w.LayoutOpts.CustomCoords[nodeID]
	c.X = x
	c.Y = y
	w.LayoutOpts.CustomCoords[nodeID] = c
	return w
}

// SetNodeTier assigns an explicit column tier index for a specific node.
func (w *Workflow) SetNodeTier(nodeID string, tier int) *Workflow {
	if w.LayoutOpts.NodeTiers == nil {
		w.LayoutOpts.NodeTiers = make(map[string]int)
	}
	w.LayoutOpts.NodeTiers[nodeID] = tier
	return w
}

// SetLayoutOptions sets the complete LayoutOptions for the workflow.
func (w *Workflow) SetLayoutOptions(opts LayoutOptions) *Workflow {
	w.LayoutOpts = opts
	return w
}

// SetNodeColor configures BPMN in Color (bioc:stroke and bioc:fill) for a node.
func (w *Workflow) SetNodeColor(nodeID string, strokeHex, fillHex string) *Workflow {
	if w.NodeColors == nil {
		w.NodeColors = make(map[string][2]string)
	}
	w.NodeColors[nodeID] = [2]string{strokeHex, fillHex}
	return w
}

func (w *Workflow) Builder() *Workflow {
	return w
}

func (w *Workflow) StartEvent(id ...string) *Workflow {
	startID := "start"
	if len(id) > 0 {
		startID = id[0]
	}
	w.Nodes = append(w.Nodes, map[string]interface{}{
		"type": "startEvent",
		"id":   startID,
		"name": "Start",
	})
	w.currentNodeID = startID
	return w
}

func (w *Workflow) EndEvent(id string, name string) *Workflow {
	w.Nodes = append(w.Nodes, map[string]interface{}{
		"type": "endEvent",
		"id":   id,
		"name": name,
	})
	return w
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func toCamelCase(s string) string {
	if s == "wasm" {
		return "wasmPath"
	}
	if s == "result_variable" {
		return "resultVar"
	}
	if !strings.Contains(s, "_") {
		return s
	}
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		parts[i] = capitalize(parts[i])
	}
	return strings.Join(parts, "")
}

func populateNodeProperties(node map[string]interface{}, opts []map[string]interface{}) {
	for _, opt := range opts {
		for k, v := range opt {
			key := toCamelCase(k)
			node[key] = v
		}
	}
}

func (w *Workflow) ServiceTask(id string, name string, topic string, options ...map[string]interface{}) *Workflow {
	node := map[string]interface{}{
		"type":  "serviceTask",
		"id":    id,
		"name":  name,
		"topic": topic,
	}
	populateNodeProperties(node, options)
	w.Nodes = append(w.Nodes, node)
	return w
}

func (w *Workflow) AITask(id string, name string, options ...map[string]interface{}) *Workflow {
	node := map[string]interface{}{
		"type": "aiServiceTask",
		"id":   id,
		"name": name,
	}
	populateNodeProperties(node, options)
	w.Nodes = append(w.Nodes, node)
	return w
}

func (w *Workflow) UserTask(id string, name string, options ...map[string]interface{}) *Workflow {
	node := map[string]interface{}{
		"type": "userTask",
		"id":   id,
		"name": name,
	}
	populateNodeProperties(node, options)
	w.Nodes = append(w.Nodes, node)
	return w
}

func (w *Workflow) ExclusiveGateway(id string, name string) *Workflow {
	w.Nodes = append(w.Nodes, map[string]interface{}{
		"type": "exclusiveGateway",
		"id":   id,
		"name": name,
	})
	return w
}

func (w *Workflow) ParallelGateway(id string, name string) *Workflow {
	w.Nodes = append(w.Nodes, map[string]interface{}{
		"type": "parallelGateway",
		"id":   id,
		"name": name,
	})
	return w
}

func (w *Workflow) EventBasedGateway(id string, name string) *Workflow {
	w.Nodes = append(w.Nodes, map[string]interface{}{
		"type": "eventBasedGateway",
		"id":   id,
		"name": name,
	})
	return w
}

func (w *Workflow) CallActivity(id string, name string, calledElement string, options ...map[string]interface{}) *Workflow {
	node := map[string]interface{}{
		"type":          "callActivity",
		"id":            id,
		"name":          name,
		"calledElement": calledElement,
	}
	populateNodeProperties(node, options)
	w.Nodes = append(w.Nodes, node)
	return w
}

func (w *Workflow) BusinessRuleTask(id string, name string, decisionRef string, options ...map[string]interface{}) *Workflow {
	node := map[string]interface{}{
		"type":        "businessRuleTask",
		"id":          id,
		"name":        name,
		"decisionRef": decisionRef,
	}
	populateNodeProperties(node, options)
	w.Nodes = append(w.Nodes, node)
	return w
}

func (w *Workflow) SequenceFlow(source, target string) *Workflow {
	w.Flows = append(w.Flows, map[string]interface{}{
		"id":        fmt.Sprintf("flow-%s-%s", source, target),
		"source":    source,
		"target":    target,
		"condition": "",
	})
	return w
}

func (w *Workflow) SequenceFlowWithCondition(source, target string, condition string) *Workflow {
	w.Flows = append(w.Flows, map[string]interface{}{
		"id":        fmt.Sprintf("flow-%s-%s", source, target),
		"source":    source,
		"target":    target,
		"condition": condition,
	})
	return w
}

func (w *Workflow) findNode(id string) map[string]interface{} {
	for _, n := range w.Nodes {
		if n["id"] == id {
			return n
		}
	}
	return nil
}

// Branch is a local scope context for conditional branches.
type Branch struct {
	workflow      *Workflow
	gatewayID     string
	currentNodeID string
	isConditional bool
	condition     string
	hasEnded      bool
}

// WhenBuilder facilitates chaining Then() after When() at the workflow level.
type WhenBuilder struct {
	workflow  *Workflow
	gatewayID string
	condition string
}

// ThenBuilder facilitates chaining Else() after Then() at the workflow level.
type ThenBuilder struct {
	workflow  *Workflow
	gatewayID string
}

// WhenBranchBuilder facilitates chaining Then() after When() at the branch level.
type WhenBranchBuilder struct {
	branch    *Branch
	gatewayID string
	condition string
}

// ThenBranchBuilder facilitates chaining Else() after Then() at the branch level.
type ThenBranchBuilder struct {
	branch    *Branch
	gatewayID string
}

// IfElseBuilder facilitates canonical If-Else and If-Then-Else chaining at the workflow level.
type IfElseBuilder struct {
	workflow  *Workflow
	gatewayID string
	condition string
}

// IfElseBranchBuilder facilitates canonical If-Else chaining at the branch level.
type IfElseBranchBuilder struct {
	branch    *Branch
	gatewayID string
	condition string
}

func (w *Workflow) connectNode(id string) {
	node := w.findNode(id)
	hasStart := false
	for _, n := range w.Nodes {
		if n["type"] == "startEvent" {
			hasStart = true
			break
		}
	}
	if !hasStart && node != nil && node["type"] != "startEvent" {
		w.StartEvent("start")
		w.SequenceFlow("start", id)
		w.currentNodeID = id
		return
	}

	if len(w.pendingMerges) > 0 {
		for _, sourceID := range w.pendingMerges {
			w.SequenceFlow(sourceID, id)
		}
		w.pendingMerges = nil
	} else if w.currentNodeID != "" && w.currentNodeID != id {
		w.SequenceFlow(w.currentNodeID, id)
	}
	w.currentNodeID = id
}

// Start initiates the workflow sequential path.
func (w *Workflow) Start(id ...string) *Workflow {
	startID := "start"
	if len(id) > 0 {
		startID = id[0]
	}
	w.StartEvent(startID)
	w.connectNode(startID)
	return w
}

// End terminates the main workflow sequential path.
func (w *Workflow) End(id, name string) *Workflow {
	w.EndEvent(id, name)
	w.connectNode(id)
	w.currentNodeID = "" // terminate main path
	return w
}

// User appends a user task and links it sequentially.
func (w *Workflow) User(id, name string, options ...map[string]interface{}) *Workflow {
	w.UserTask(id, name, options...)
	w.connectNode(id)
	return w
}

// Service appends a service task and links it sequentially.
func (w *Workflow) Service(id, name, topic string, options ...map[string]interface{}) *Workflow {
	w.ServiceTask(id, name, topic, options...)
	w.connectNode(id)
	return w
}

// AI appends an AI orchestration task and links it sequentially.
func (w *Workflow) AI(id, name string, options ...map[string]interface{}) *Workflow {
	w.AITask(id, name, options...)
	w.connectNode(id)
	return w
}

func (w *Workflow) Call(id, name, calledElement string, options ...map[string]interface{}) *Workflow {
	w.CallActivity(id, name, calledElement, options...)
	w.connectNode(id)
	return w
}

func (w *Workflow) BusinessRule(id, name, decisionRef string, options ...map[string]interface{}) *Workflow {
	w.BusinessRuleTask(id, name, decisionRef, options...)
	w.connectNode(id)
	return w
}

// If defines a canonical conditional branch path on the main workflow.
// Supports both two-argument direct block form:
//   w.If("order.amount > 100", func(b *Branch) { ... }).Else(func(b *Branch) { ... })
// and chained BDD/DSL builder form:
//   w.If("order.amount > 100").Then(func(b *Branch) { ... }).Else(func(b *Branch) { ... })
func (w *Workflow) If(condition interface{}, thenFn ...func(flow *Branch)) *IfElseBuilder {
	gwID := fmt.Sprintf("gw_%s_decision", w.currentNodeID)
	w.ExclusiveGateway(gwID, "Decision Gateway")
	w.connectNode(gwID)

	var condStr string
	switch c := condition.(type) {
	case string:
		condStr = c
	case fmt.Stringer:
		condStr = c.String()
	default:
		condStr = fmt.Sprintf("%v", c)
	}

	builder := &IfElseBuilder{
		workflow:  w,
		gatewayID: gwID,
		condition: condStr,
	}

	if len(thenFn) > 0 && thenFn[0] != nil {
		thenBranch := &Branch{
			workflow:      w,
			gatewayID:     gwID,
			currentNodeID: gwID,
			isConditional: true,
			condition:     condStr,
		}
		thenFn[0](thenBranch)
		if !thenBranch.hasEnded && thenBranch.currentNodeID != gwID {
			w.pendingMerges = append(w.pendingMerges, thenBranch.currentNodeID)
		}
	}

	return builder
}

// Then defines the branch execution when the If condition evaluates to true.
func (ie *IfElseBuilder) Then(thenFn func(flow *Branch)) *IfElseBuilder {
	thenBranch := &Branch{
		workflow:      ie.workflow,
		gatewayID:     ie.gatewayID,
		currentNodeID: ie.gatewayID,
		isConditional: true,
		condition:     ie.condition,
	}
	thenFn(thenBranch)
	if !thenBranch.hasEnded && thenBranch.currentNodeID != ie.gatewayID {
		ie.workflow.pendingMerges = append(ie.workflow.pendingMerges, thenBranch.currentNodeID)
	}
	return ie
}

// ElseIf defines an alternative conditional branch from the same decision gateway.
func (ie *IfElseBuilder) ElseIf(condition interface{}, thenFn ...func(flow *Branch)) *IfElseBuilder {
	var condStr string
	switch c := condition.(type) {
	case string:
		condStr = c
	case fmt.Stringer:
		condStr = c.String()
	default:
		condStr = fmt.Sprintf("%v", c)
	}
	ie.condition = condStr

	if len(thenFn) > 0 && thenFn[0] != nil {
		branch := &Branch{
			workflow:      ie.workflow,
			gatewayID:     ie.gatewayID,
			currentNodeID: ie.gatewayID,
			isConditional: true,
			condition:     condStr,
		}
		thenFn[0](branch)
		if !branch.hasEnded && branch.currentNodeID != ie.gatewayID {
			ie.workflow.pendingMerges = append(ie.workflow.pendingMerges, branch.currentNodeID)
		}
	}
	return ie
}

// When allows seamlessly switching from If/ElseIf to When...Then semantics.
func (ie *IfElseBuilder) When(condition interface{}) *WhenBuilder {
	var condStr string
	switch c := condition.(type) {
	case string:
		condStr = c
	case fmt.Stringer:
		condStr = c.String()
	default:
		condStr = fmt.Sprintf("%v", c)
	}
	return &WhenBuilder{
		workflow:  ie.workflow,
		gatewayID: ie.gatewayID,
		condition: condStr,
	}
}

// Else defines the default fallback path when conditions evaluate to false.
func (ie *IfElseBuilder) Else(elseFn func(flow *Branch)) *Workflow {
	elseBranch := &Branch{
		workflow:      ie.workflow,
		gatewayID:     ie.gatewayID,
		currentNodeID: ie.gatewayID,
		isConditional: false,
	}
	elseFn(elseBranch)
	if !elseBranch.hasEnded && elseBranch.currentNodeID != ie.gatewayID {
		ie.workflow.pendingMerges = append(ie.workflow.pendingMerges, elseBranch.currentNodeID)
	}
	return ie.workflow
}

// Otherwise is a canonical alias for Else.
func (ie *IfElseBuilder) Otherwise(elseFn func(flow *Branch)) *Workflow {
	return ie.Else(elseFn)
}

// When defines a conditional branch path starting condition on the main workflow.
func (w *Workflow) When(condition interface{}) *WhenBuilder {
	gwID := fmt.Sprintf("gw_%s_decision", w.currentNodeID)
	w.ExclusiveGateway(gwID, "Decision Gateway")
	w.connectNode(gwID)

	var condStr string
	switch c := condition.(type) {
	case string:
		condStr = c
	case fmt.Stringer:
		condStr = c.String()
	default:
		condStr = fmt.Sprintf("%v", c)
	}

	return &WhenBuilder{
		workflow:  w,
		gatewayID: gwID,
		condition: condStr,
	}
}

// Then defines the branch steps when the condition is met.
func (wb *WhenBuilder) Then(thenFn func(flow *Branch)) *ThenBuilder {
	thenBranch := &Branch{
		workflow:      wb.workflow,
		gatewayID:     wb.gatewayID,
		currentNodeID: wb.gatewayID,
		isConditional: true,
		condition:     wb.condition,
	}

	thenFn(thenBranch)

	if !thenBranch.hasEnded && thenBranch.currentNodeID != wb.gatewayID {
		wb.workflow.pendingMerges = append(wb.workflow.pendingMerges, thenBranch.currentNodeID)
	}

	return &ThenBuilder{
		workflow:  wb.workflow,
		gatewayID: wb.gatewayID,
	}
}

// Else defines the default path on the main workflow when the condition evaluates to false.
func (tb *ThenBuilder) Else(elseFn func(flow *Branch)) *Workflow {
	elseBranch := &Branch{
		workflow:      tb.workflow,
		gatewayID:     tb.gatewayID,
		currentNodeID: tb.gatewayID,
		isConditional: false,
	}

	elseFn(elseBranch)

	if !elseBranch.hasEnded && elseBranch.currentNodeID != tb.gatewayID {
		tb.workflow.pendingMerges = append(tb.workflow.pendingMerges, elseBranch.currentNodeID)
	}

	return tb.workflow
}

// Otherwise is a canonical alias for Else on ThenBuilder.
func (tb *ThenBuilder) Otherwise(elseFn func(flow *Branch)) *Workflow {
	return tb.Else(elseFn)
}

// When allows chaining additional conditions (When...Then...When...Then...Else).
func (tb *ThenBuilder) When(condition interface{}) *WhenBuilder {
	var condStr string
	switch c := condition.(type) {
	case string:
		condStr = c
	case fmt.Stringer:
		condStr = c.String()
	default:
		condStr = fmt.Sprintf("%v", c)
	}
	return &WhenBuilder{
		workflow:  tb.workflow,
		gatewayID: tb.gatewayID,
		condition: condStr,
	}
}

// ElseIf allows chaining alternative conditions onto ThenBuilder.
func (tb *ThenBuilder) ElseIf(condition interface{}, thenFn ...func(flow *Branch)) *IfElseBuilder {
	ie := &IfElseBuilder{
		workflow:  tb.workflow,
		gatewayID: tb.gatewayID,
	}
	return ie.ElseIf(condition, thenFn...)
}



func (b *Branch) connectNode(id string) {
	if b.hasEnded {
		return
	}
	isBackEdge := false
	for i := 0; i < len(b.workflow.Nodes)-1; i++ {
		if b.workflow.Nodes[i]["id"] == id {
			isBackEdge = true
			break
		}
	}

	if len(b.workflow.pendingMerges) > 0 {
		for _, sourceID := range b.workflow.pendingMerges {
			b.workflow.SequenceFlow(sourceID, id)
		}
		b.workflow.pendingMerges = nil
		b.currentNodeID = id
		if isBackEdge {
			b.hasEnded = true
		}
		return
	}
	if b.currentNodeID == b.gatewayID {
		if b.isConditional {
			b.workflow.SequenceFlowWithCondition(b.gatewayID, id, b.condition)
		} else {
			// Connect default sequence flow
			b.workflow.SequenceFlow(b.gatewayID, id)
		}
	} else if b.currentNodeID != "" && b.currentNodeID != id {
		b.workflow.SequenceFlow(b.currentNodeID, id)
	}
	b.currentNodeID = id
	if isBackEdge {
		b.hasEnded = true
	}
}

// User appends a user task inside a branch.
func (b *Branch) User(id, name string, options ...map[string]interface{}) *Branch {
	b.workflow.UserTask(id, name, options...)
	b.connectNode(id)
	return b
}

// Service appends a service task inside a branch.
func (b *Branch) Service(id, name, topic string, options ...map[string]interface{}) *Branch {
	b.workflow.ServiceTask(id, name, topic, options...)
	b.connectNode(id)
	return b
}

// AI appends an AI orchestration task inside a branch.
func (b *Branch) AI(id, name string, options ...map[string]interface{}) *Branch {
	b.workflow.AITask(id, name, options...)
	b.connectNode(id)
	return b
}

func (b *Branch) Call(id, name, calledElement string, options ...map[string]interface{}) *Branch {
	b.workflow.CallActivity(id, name, calledElement, options...)
	b.connectNode(id)
	return b
}

func (b *Branch) BusinessRule(id, name, decisionRef string, options ...map[string]interface{}) *Branch {
	b.workflow.BusinessRuleTask(id, name, decisionRef, options...)
	b.connectNode(id)
	return b
}

// End terminates the branch.
func (b *Branch) End(id, name string) *Branch {
	b.workflow.EndEvent(id, name)
	b.connectNode(id)
	b.hasEnded = true
	return b
}

// When defines a nested conditional branch path starting condition inside a branch.
func (b *Branch) When(condition interface{}) *WhenBranchBuilder {
	gwID := fmt.Sprintf("gw_%s_decision", b.currentNodeID)
	b.workflow.ExclusiveGateway(gwID, "Decision Gateway")
	b.connectNode(gwID)

	var condStr string
	switch c := condition.(type) {
	case string:
		condStr = c
	case fmt.Stringer:
		condStr = c.String()
	default:
		condStr = fmt.Sprintf("%v", c)
	}

	return &WhenBranchBuilder{
		branch:    b,
		gatewayID: gwID,
		condition: condStr,
	}
}

// Then defines the nested branch steps when the condition is met.
func (wbb *WhenBranchBuilder) Then(thenFn func(sub *Branch)) *ThenBranchBuilder {
	thenBranch := &Branch{
		workflow:      wbb.branch.workflow,
		gatewayID:     wbb.gatewayID,
		currentNodeID: wbb.gatewayID,
		isConditional: true,
		condition:     wbb.condition,
	}

	thenFn(thenBranch)

	if !thenBranch.hasEnded && thenBranch.currentNodeID != wbb.gatewayID {
		wbb.branch.workflow.pendingMerges = append(wbb.branch.workflow.pendingMerges, thenBranch.currentNodeID)
	}

	return &ThenBranchBuilder{
		branch:    wbb.branch,
		gatewayID: wbb.gatewayID,
	}
}

// Else defines the nested default path when the condition evaluates to false.
func (tbb *ThenBranchBuilder) Else(elseFn func(sub *Branch)) *Branch {
	elseBranch := &Branch{
		workflow:      tbb.branch.workflow,
		gatewayID:     tbb.gatewayID,
		currentNodeID: tbb.gatewayID,
		isConditional: false,
	}

	elseFn(elseBranch)

	if !elseBranch.hasEnded && elseBranch.currentNodeID != tbb.gatewayID {
		tbb.branch.workflow.pendingMerges = append(tbb.branch.workflow.pendingMerges, elseBranch.currentNodeID)
	}

	return tbb.branch
}

// Otherwise is a canonical alias for Else on ThenBranchBuilder.
func (tbb *ThenBranchBuilder) Otherwise(elseFn func(sub *Branch)) *Branch {
	return tbb.Else(elseFn)
}

// If defines a nested canonical conditional branch inside a branch.
func (b *Branch) If(condition interface{}, thenFn ...func(sub *Branch)) *IfElseBranchBuilder {
	gwID := fmt.Sprintf("gw_%s_decision", b.currentNodeID)
	b.workflow.ExclusiveGateway(gwID, "Decision Gateway")
	b.connectNode(gwID)

	var condStr string
	switch c := condition.(type) {
	case string:
		condStr = c
	case fmt.Stringer:
		condStr = c.String()
	default:
		condStr = fmt.Sprintf("%v", c)
	}

	builder := &IfElseBranchBuilder{
		branch:    b,
		gatewayID: gwID,
		condition: condStr,
	}

	if len(thenFn) > 0 && thenFn[0] != nil {
		thenBranch := &Branch{
			workflow:      b.workflow,
			gatewayID:     gwID,
			currentNodeID: gwID,
			isConditional: true,
			condition:     condStr,
		}
		thenFn[0](thenBranch)
		if !thenBranch.hasEnded && thenBranch.currentNodeID != gwID {
			b.workflow.pendingMerges = append(b.workflow.pendingMerges, thenBranch.currentNodeID)
		}
	}

	return builder
}

// Then defines branch steps when the nested If condition evaluates to true.
func (ieb *IfElseBranchBuilder) Then(thenFn func(flow *Branch)) *IfElseBranchBuilder {
	thenBranch := &Branch{
		workflow:      ieb.branch.workflow,
		gatewayID:     ieb.gatewayID,
		currentNodeID: ieb.gatewayID,
		isConditional: true,
		condition:     ieb.condition,
	}
	thenFn(thenBranch)
	if !thenBranch.hasEnded && thenBranch.currentNodeID != ieb.gatewayID {
		ieb.branch.workflow.pendingMerges = append(ieb.branch.workflow.pendingMerges, thenBranch.currentNodeID)
	}
	return ieb
}

// ElseIf adds an alternative conditional branch from the same decision gateway.
func (ieb *IfElseBranchBuilder) ElseIf(condition interface{}, thenFn ...func(flow *Branch)) *IfElseBranchBuilder {
	var condStr string
	switch c := condition.(type) {
	case string:
		condStr = c
	case fmt.Stringer:
		condStr = c.String()
	default:
		condStr = fmt.Sprintf("%v", c)
	}
	ieb.condition = condStr

	if len(thenFn) > 0 && thenFn[0] != nil {
		branch := &Branch{
			workflow:      ieb.branch.workflow,
			gatewayID:     ieb.gatewayID,
			currentNodeID: ieb.gatewayID,
			isConditional: true,
			condition:     condStr,
		}
		thenFn[0](branch)
		if !branch.hasEnded && branch.currentNodeID != ieb.gatewayID {
			ieb.branch.workflow.pendingMerges = append(ieb.branch.workflow.pendingMerges, branch.currentNodeID)
		}
	}
	return ieb
}

// Else defines the default fallback path when all nested conditions evaluate to false.
func (ieb *IfElseBranchBuilder) Else(elseFn func(flow *Branch)) *Branch {
	elseBranch := &Branch{
		workflow:      ieb.branch.workflow,
		gatewayID:     ieb.gatewayID,
		currentNodeID: ieb.gatewayID,
		isConditional: false,
	}
	elseFn(elseBranch)
	if !elseBranch.hasEnded && elseBranch.currentNodeID != ieb.gatewayID {
		ieb.branch.workflow.pendingMerges = append(ieb.branch.workflow.pendingMerges, elseBranch.currentNodeID)
	}
	return ieb.branch
}

// Otherwise is a canonical alias for Else.
func (ieb *IfElseBranchBuilder) Otherwise(elseFn func(flow *Branch)) *Branch {
	return ieb.Else(elseFn)
}





type Expression struct {
	expr string
}

func (e Expression) String() string {
	return e.expr
}

type Variable struct {
	name string
}

func V(name string) Variable {
	return Variable{name: name}
}

func Var(name string) Variable {
	return Variable{name: name}
}

func (v Variable) Eq(val interface{}) Expression {
	valStr := fmt.Sprintf("%v", val)
	if b, ok := val.(bool); ok {
		if b {
			valStr = "true"
		} else {
			valStr = "false"
		}
	}
	return Expression{expr: fmt.Sprintf("%s == %s", v.name, valStr)}
}

func (v Variable) Neq(val interface{}) Expression {
	valStr := fmt.Sprintf("%v", val)
	if b, ok := val.(bool); ok {
		if b {
			valStr = "true"
		} else {
			valStr = "false"
		}
	}
	return Expression{expr: fmt.Sprintf("%s != %s", v.name, valStr)}
}

func (v Variable) Gt(val interface{}) Expression {
	return Expression{expr: fmt.Sprintf("%s > %v", v.name, val)}
}

func (v Variable) Gte(val interface{}) Expression {
	return Expression{expr: fmt.Sprintf("%s >= %v", v.name, val)}
}

func (v Variable) Lt(val interface{}) Expression {
	return Expression{expr: fmt.Sprintf("%s < %v", v.name, val)}
}

func (v Variable) Lte(val interface{}) Expression {
	return Expression{expr: fmt.Sprintf("%s <= %v", v.name, val)}
}

func (w *Workflow) ToJSON() ([]byte, error) {
	return json.Marshal(w)
}

