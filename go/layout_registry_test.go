package nativebpm_test

import (
	"testing"

	nativebpm "gitlab.com/nativebpm/sdk/go"
)

func TestLayout_Catalog(t *testing.T) {
	presets := nativebpm.ListLayoutPresets()
	if len(presets) < 4 {
		t.Fatalf("Expected at least 4 presets, got %d", len(presets))
	}

	presetMap := make(map[nativebpm.LayoutPreset]nativebpm.LayoutPresetInfo)
	for _, p := range presets {
		presetMap[p.Preset] = p
	}

	expected := []nativebpm.LayoutPreset{
		nativebpm.LayoutAuto,
		nativebpm.LayoutTiered,
		nativebpm.LayoutCenterHub,
		nativebpm.LayoutLinear,
	}

	for _, exp := range expected {
		info, ok := presetMap[exp]
		if !ok {
			t.Errorf("Expected preset '%s' in catalog", exp)
			continue
		}
		if info.Name == "" || info.Description == "" {
			t.Errorf("Preset '%s' must have Name and Description", exp)
		}
	}

	// Test GetLayoutPresetInfo
	autoInfo, ok := nativebpm.GetLayoutPresetInfo(nativebpm.LayoutAuto)
	if !ok || autoInfo.Preset != nativebpm.LayoutAuto {
		t.Errorf("GetLayoutPresetInfo failed for LayoutAuto")
	}
}

func TestLayout_AutoSelection_Linear(t *testing.T) {
	wf := nativebpm.NewWorkflow("linear_process", "Simple Linear Process")
	wf.Start("start").
		Service("step1", "Step 1", "topic1").
		Service("step2", "Step 2", "topic2").
		End("end", "End")

	// Set LayoutAuto explicitly
	wf.SetLayoutPreset(nativebpm.LayoutAuto)

	xmlBytes, err := wf.ToBPMNXML()
	if err != nil {
		t.Fatalf("ToBPMNXML failed: %v", err)
	}
	if len(xmlBytes) == 0 {
		t.Fatal("Expected non-empty BPMN XML")
	}

	// Metrics should report Linear
	metrics := nativebpm.AnalyzeGraphComplexity(wf.Nodes, wf.Flows)
	if metrics.NodeCount != 4 {
		t.Errorf("Expected 4 nodes, got %d", metrics.NodeCount)
	}
	selected := nativebpm.AutoSelectLayoutPreset(wf.Nodes, wf.Flows)
	if selected != nativebpm.LayoutLinear {
		t.Errorf("Expected auto-selection '%s', got '%s'", nativebpm.LayoutLinear, selected)
	}
}

func TestLayout_AutoSelection_CenterHub(t *testing.T) {
	wf := nativebpm.NewWorkflow("hub_process", "Central Hub Topology")
	wf.Start("start").
		ParallelGateway("gw_fork", "Smart Router Fork")

	// 4 fan-out branches from gw_fork
	for i := 1; i <= 4; i++ {
		taskID := "task_branch_" + string(rune('0'+i))
		wf.ServiceTask(taskID, "Worker Branch", "topic")
		wf.SequenceFlow("gw_fork", taskID)
		wf.SequenceFlow(taskID, "gw_join")
	}

	wf.ParallelGateway("gw_join", "Join Gateway").
		End("end", "End")

	metrics := nativebpm.AnalyzeGraphComplexity(wf.Nodes, wf.Flows)
	if metrics.DetectedHubID != "gw_fork" && metrics.DetectedHubID != "gw_join" {
		t.Errorf("Expected detected hub 'gw_fork' or 'gw_join', got '%s'", metrics.DetectedHubID)
	}
	if metrics.MaxOutDegree != 4 {
		t.Errorf("Expected MaxOutDegree 4, got %d", metrics.MaxOutDegree)
	}

	selected := nativebpm.AutoSelectLayoutPreset(wf.Nodes, wf.Flows)
	if selected != nativebpm.LayoutCenterHub {
		t.Errorf("Expected auto-selection '%s', got '%s'", nativebpm.LayoutCenterHub, selected)
	}
}

func TestLayout_RegisterCustomPreset(t *testing.T) {
	customPreset := nativebpm.LayoutPreset("custom_grid_matrix")
	customCalled := false

	customStrategy := func(nodes []map[string]interface{}, flows []map[string]interface{}, opts nativebpm.LayoutOptions) map[string]nativebpm.Coords {
		customCalled = true
		res := make(map[string]nativebpm.Coords)
		for i, n := range nodes {
			id, _ := n["id"].(string)
			res[id] = nativebpm.Coords{
				X:      float64(i * 100),
				Y:      float64(i * 50),
				Width:  100,
				Height: 80,
			}
		}
		return res
	}

	nativebpm.RegisterLayoutPreset(customPreset, nativebpm.LayoutPresetInfo{
		Preset:      customPreset,
		Name:        "Custom Grid Matrix",
		Description: "Experimental diagonal layout",
		BestFor:     "Custom test topologies",
	}, customStrategy)

	// Verify preset is listed in catalog
	info, ok := nativebpm.GetLayoutPresetInfo(customPreset)
	if !ok || info.Name != "Custom Grid Matrix" {
		t.Fatalf("Expected custom preset in catalog")
	}

	wf := nativebpm.NewWorkflow("custom_proc", "Custom Layout Test")
	wf.Start("start").Service("task1", "Task 1", "topic").End("end", "End")
	wf.SetLayoutPreset(customPreset)

	xmlBytes, err := wf.ToBPMNXML()
	if err != nil {
		t.Fatalf("ToBPMNXML failed with custom preset: %v", err)
	}
	if len(xmlBytes) == 0 {
		t.Fatal("Expected non-empty XML")
	}
	if !customCalled {
		t.Errorf("Custom layout strategy was not called during ToBPMNXML")
	}
}

func TestLayout_PresetsProduceDistinctCoordinates(t *testing.T) {
	buildSampleWorkflow := func(preset nativebpm.LayoutPreset) *nativebpm.Workflow {
		wf := nativebpm.NewWorkflow("distinct_test", "Distinct Coordinates Test")
		wf.Start("start").
			ParallelGateway("gw_hub", "Argo Hub Router")

		for i := 1; i <= 4; i++ {
			tID := "task_branch_" + string(rune('0'+i))
			wf.ServiceTask(tID, "Parallel Service", "topic")
			wf.SequenceFlow("gw_hub", tID)
			wf.SequenceFlow(tID, "gw_join")
		}

		wf.ParallelGateway("gw_join", "Join Convergence").
			ServiceTask("task_commit", "Commit Result", "topic").
			SequenceFlow("gw_join", "task_commit").
			End("end", "End")
		wf.SequenceFlow("task_commit", "end")

		opts := nativebpm.DefaultLayoutOptions()
		opts.Preset = preset
		opts.CenterHubID = "gw_hub"
		wf.SetLayoutOptions(opts)
		return wf
	}

	tieredXML, err := buildSampleWorkflow(nativebpm.LayoutTiered).ToBPMNXML()
	if err != nil {
		t.Fatalf("Tiered XML failed: %v", err)
	}

	centerHubXML, err := buildSampleWorkflow(nativebpm.LayoutCenterHub).ToBPMNXML()
	if err != nil {
		t.Fatalf("CenterHub XML failed: %v", err)
	}

	clusterMeshXML, err := buildSampleWorkflow(nativebpm.LayoutClusterMesh).ToBPMNXML()
	if err != nil {
		t.Fatalf("ClusterMesh XML failed: %v", err)
	}

	linearXML, err := buildSampleWorkflow(nativebpm.LayoutLinear).ToBPMNXML()
	if err != nil {
		t.Fatalf("Linear XML failed: %v", err)
	}

	// The generated BPMN DI shapes must NOT be identical between different presets!
	if string(tieredXML) == string(centerHubXML) {
		t.Errorf("LayoutTiered and LayoutCenterHub must produce distinctly different XML layouts!")
	}
	if string(centerHubXML) == string(clusterMeshXML) {
		t.Errorf("LayoutCenterHub and LayoutClusterMesh must produce distinctly different XML layouts!")
	}
	if string(tieredXML) == string(linearXML) {
		t.Errorf("LayoutTiered and LayoutLinear must produce distinctly different XML layouts!")
	}
}

func TestLayout_Orientation_TransposeBPMNXML(t *testing.T) {
	wf := nativebpm.NewWorkflow("transposition_test", "Vertical Transposition Test")
	wf.Start("start").
		Service("task1", "Step 1", "topic1").
		Service("task2", "Step 2", "topic2").
		End("end", "End")

	// Standard horizontal XML
	horizXML, err := wf.ToBPMNXML()
	if err != nil {
		t.Fatalf("ToBPMNXML failed: %v", err)
	}

	// Transpose to vertical
	vertXML, err := nativebpm.TransposeBPMNXML(horizXML)
	if err != nil {
		t.Fatalf("TransposeBPMNXML failed: %v", err)
	}

	if string(vertXML) == string(horizXML) {
		t.Errorf("Transposed XML must be different from horizontal XML!")
	}

	// Also test workflow with Orientation: OrientationVertical
	opts := nativebpm.DefaultLayoutOptions()
	opts.Orientation = nativebpm.OrientationVertical
	wf.SetLayoutOptions(opts)

	autoVertXML, err := wf.ToBPMNXML()
	if err != nil {
		t.Fatalf("ToBPMNXML with OrientationVertical failed: %v", err)
	}

	if string(autoVertXML) == string(horizXML) {
		t.Errorf("ToBPMNXML with OrientationVertical must produce vertical layout!")
	}
}


