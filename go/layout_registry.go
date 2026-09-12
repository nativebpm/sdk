package nativebpm

import (
	"log/slog"
	"sync"
)

const (
	LayoutAuto        LayoutPreset = "auto"
	LayoutClusterMesh LayoutPreset = "cluster_mesh"
)

// LayoutPresetInfo encapsulates metadata and human-readable guidance for layout presets.
type LayoutPresetInfo struct {
	Preset      LayoutPreset `json:"preset"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	BestFor     string       `json:"bestFor"`
	IsDefault   bool         `json:"isDefault,omitempty"`
}

// GraphMetrics captures topological metrics of a BPMN workflow.
type GraphMetrics struct {
	NodeCount       int    `json:"nodeCount"`
	FlowCount       int    `json:"flowCount"`
	GatewayCount    int    `json:"gatewayCount"`
	MaxOutDegree    int    `json:"maxOutDegree"`
	MaxInDegree     int    `json:"maxInDegree"`
	DetectedHubID   string `json:"detectedHubId,omitempty"`
	ComplexityLevel string `json:"complexityLevel"`
}

// LayoutStrategy is a pluggable function that computes coordinates for workflow elements.
type LayoutStrategy func(nodes []map[string]interface{}, flows []map[string]interface{}, opts LayoutOptions) map[string]Coords

type layoutEntry struct {
	info     LayoutPresetInfo
	strategy LayoutStrategy
}

var (
	registryMu sync.RWMutex
	registry   = make(map[LayoutPreset]layoutEntry)
)

func init() {
	// Register default built-in layout presets
	RegisterLayoutPreset(LayoutAuto, LayoutPresetInfo{
		Preset:      LayoutAuto,
		Name:        "Intelligent Auto-Selection",
		Description: "Dynamically analyzes graph complexity, element count, and hub fan-out to pick the best layout with structured logging.",
		BestFor:     "Zero-configuration workflows (recommended default).",
		IsDefault:   true,
	}, nil)

	RegisterLayoutPreset(LayoutTiered, LayoutPresetInfo{
		Preset:      LayoutTiered,
		Name:        "Tiered Multilevel Pipeline",
		Description: "Arranges nodes in topological columns left-to-right with Zero-Overlap collision resolution.",
		BestFor:     "Complex business pipelines, multi-step approvals, and standard DAGs.",
	}, nil)

	RegisterLayoutPreset(LayoutCenterHub, LayoutPresetInfo{
		Preset:      LayoutCenterHub,
		Name:        "Radial Center-Hub & Parallel Tunnels",
		Description: "Centers a primary gateway/router with balanced parallel fan-out branches radiating outwards.",
		BestFor:     "Network topologies, Cloudflare Argo smart routing, load balancers, and multi-region hubs.",
	}, nil)

	RegisterLayoutPreset(LayoutLinear, LayoutPresetInfo{
		Preset:      LayoutLinear,
		Name:        "Compact Linear Sequence",
		Description: "Single-tier horizontal flow with minimal vertical displacement.",
		BestFor:     "Short sequential processes (<= 8 nodes) without branching.",
	}, nil)

	RegisterLayoutPreset(LayoutClusterMesh, LayoutPresetInfo{
		Preset:      LayoutClusterMesh,
		Name:        "Cluster Mesh Interconnect",
		Description: "Matrix layout designed for distributed quorum consensus, multi-master meshes, and interconnected nodes.",
		BestFor:     "Distributed databases, Raft quorums, KeyDB/Redis multi-master clusters.",
	}, nil)
}

// RegisterLayoutPreset registers a layout preset in the global registry.
func RegisterLayoutPreset(preset LayoutPreset, info LayoutPresetInfo, strategy LayoutStrategy) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[preset] = layoutEntry{
		info:     info,
		strategy: strategy,
	}
}

// ListLayoutPresets returns all available layout presets with metadata.
func ListLayoutPresets() []LayoutPresetInfo {
	registryMu.RLock()
	defer registryMu.RUnlock()

	presets := make([]LayoutPresetInfo, 0, len(registry))
	// Ensure stable order
	order := []LayoutPreset{LayoutAuto, LayoutTiered, LayoutCenterHub, LayoutLinear, LayoutClusterMesh}
	seen := make(map[LayoutPreset]bool)

	for _, p := range order {
		if entry, ok := registry[p]; ok {
			presets = append(presets, entry.info)
			seen[p] = true
		}
	}

	for p, entry := range registry {
		if !seen[p] {
			presets = append(presets, entry.info)
		}
	}

	return presets
}

// GetLayoutPresetInfo retrieves metadata for a specific layout preset.
func GetLayoutPresetInfo(preset LayoutPreset) (LayoutPresetInfo, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	entry, ok := registry[preset]
	return entry.info, ok
}

// GetLayoutStrategy retrieves the custom strategy for a preset if registered.
func GetLayoutStrategy(preset LayoutPreset) (LayoutStrategy, bool) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	entry, ok := registry[preset]
	if !ok || entry.strategy == nil {
		return nil, false
	}
	return entry.strategy, true
}

// AnalyzeGraphComplexity analyzes graph structure and returns topological metrics.
func AnalyzeGraphComplexity(nodes []map[string]interface{}, flows []map[string]interface{}) GraphMetrics {
	m := GraphMetrics{
		NodeCount: len(nodes),
		FlowCount: len(flows),
	}

	outDegree := make(map[string]int)
	inDegree := make(map[string]int)
	isGateway := make(map[string]bool)

	for _, n := range nodes {
		id, _ := n["id"].(string)
		nType, _ := n["type"].(string)
		if nType == "exclusiveGateway" || nType == "parallelGateway" || nType == "eventBasedGateway" {
			m.GatewayCount++
			isGateway[id] = true
		}
	}

	for _, f := range flows {
		src, _ := f["source"].(string)
		tgt, _ := f["target"].(string)
		if src != "" {
			outDegree[src]++
		}
		if tgt != "" {
			inDegree[tgt]++
		}
	}

	maxOut := 0
	maxIn := 0
	hubCandidate := ""
	hubScore := 0

	for id, deg := range outDegree {
		if deg > maxOut {
			maxOut = deg
		}
		score := deg
		if isGateway[id] {
			score += 3
		}
		if score > hubScore && deg >= 3 {
			hubScore = score
			hubCandidate = id
		}
	}

	for id, deg := range inDegree {
		if deg > maxIn {
			maxIn = deg
		}
		score := deg
		if isGateway[id] {
			score += 3
		}
		if score > hubScore && deg >= 3 {
			hubScore = score
			hubCandidate = id
		}
	}

	m.MaxOutDegree = maxOut
	m.MaxInDegree = maxIn
	m.DetectedHubID = hubCandidate

	// Determine complexity level
	if m.NodeCount <= 8 && m.GatewayCount == 0 && m.MaxOutDegree <= 1 {
		m.ComplexityLevel = "simple"
	} else if m.FlowCount >= m.NodeCount*2 && m.NodeCount >= 10 {
		m.ComplexityLevel = "cluster_mesh"
	} else if m.MaxOutDegree >= 3 || m.NodeCount > 16 {
		m.ComplexityLevel = "high"
	} else {
		m.ComplexityLevel = "moderate"
	}

	return m
}

// AutoSelectLayoutPreset deterministically picks the best layout preset based on graph topology.
func AutoSelectLayoutPreset(nodes []map[string]interface{}, flows []map[string]interface{}) LayoutPreset {
	metrics := AnalyzeGraphComplexity(nodes, flows)

	if metrics.DetectedHubID != "" && (metrics.MaxOutDegree >= 3 || metrics.MaxInDegree >= 3) {
		return LayoutCenterHub
	}
	if metrics.ComplexityLevel == "simple" && metrics.MaxOutDegree <= 1 {
		return LayoutLinear
	}
	if metrics.ComplexityLevel == "cluster_mesh" {
		return LayoutClusterMesh
	}

	return LayoutTiered
}

// LogAutoSelection logs the reasoning behind the auto-selected layout preset using slog.
func LogAutoSelection(workflowID string, selected LayoutPreset, metrics GraphMetrics) {
	slog.Info("NativeBPM Layout Engine: auto-selected layout preset",
		"workflow_id", workflowID,
		"selected_preset", string(selected),
		"nodes", metrics.NodeCount,
		"flows", metrics.FlowCount,
		"gateways", metrics.GatewayCount,
		"max_out_degree", metrics.MaxOutDegree,
		"max_in_degree", metrics.MaxInDegree,
		"detected_hub", metrics.DetectedHubID,
		"complexity", metrics.ComplexityLevel,
	)
}
