package nativebpm_test

import (
	"encoding/xml"
	"strings"
	"testing"

	nativebpm "gitlab.com/nativebpm/sdk/go"
)

func TestWorkflow_CollaborationMultiPool(t *testing.T) {
	wf := nativebpm.NewCollaboration("sre_infrastructure", "Cluster & Test Stand Topology")

	// Pool 1: Production Cluster
	prodPool := wf.AddPool("pool_prod", "🛡️ Production Cluster (5 Nodes)", "process_prod")
	prodPool.Start("start_prod", "Incoming Traffic").
		Service("task_argo_prod", "Argo Edge Router", "argo").
		Service("task_mesh_prod", "Service Mesh Ingress", "mesh").
		End("end_prod", "Cluster Response 200")

	// Pool 2: Isolated Test Stand (Zero sequence flows to prod)
	testPool := wf.AddPool("pool_test", "🧪 Test Stand & Dev (Spot-Test)", "process_test")
	testPool.Start("start_test", "Test Probe Initiated").
		Service("task_cf_access_test", "CF Access Auth", "cf_access").
		Service("task_tunnel_test", "Argo Tunnel Test", "tunnel").
		End("end_test", "Spot Check Passed")

	xmlBytes, err := wf.ToBPMNXML()
	if err != nil {
		t.Fatalf("Collaboration ToBPMNXML failed: %v", err)
	}

	xmlStr := string(xmlBytes)

	// 1. Verify standard Collaboration & Participant tags
	expectedTags := []string{
		`<collaboration id="sre_infrastructure"`,
		`<participant id="pool_prod" name="🛡️ Production Cluster (5 Nodes)" processRef="process_prod"`,
		`<participant id="pool_test" name="🧪 Test Stand &amp; Dev (Spot-Test)" processRef="process_test"`,
		`<process id="process_prod" name="🛡️ Production Cluster (5 Nodes)" isExecutable="true"`,
		`<process id="process_test" name="🧪 Test Stand &amp; Dev (Spot-Test)" processRef=""`,
		`<bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="sre_infrastructure"`,
		`<bpmndi:BPMNShape id="pool_prod_di" bpmnElement="pool_prod" isHorizontal="true"`,
		`<bpmndi:BPMNShape id="pool_test_di" bpmnElement="pool_test" isHorizontal="true"`,
	}

	// For process_test, check isExecutable="true"
	if !strings.Contains(xmlStr, `<process id="process_test"`) {
		t.Errorf("Expected <process id=\"process_test\"> in generated XML")
	}

	for _, tag := range expectedTags[:4] {
		if !strings.Contains(xmlStr, tag) {
			t.Errorf("Expected tag/substring '%s' in generated BPMN XML", tag)
		}
	}
	for _, tag := range expectedTags[5:] {
		if !strings.Contains(xmlStr, tag) {
			t.Errorf("Expected BPMNDI tag/substring '%s' in generated BPMN XML", tag)
		}
	}

	// 2. Verify XML well-formedness by unmarshaling
	var parsed struct {
		XMLName       xml.Name `xml:"definitions"`
		Collaboration struct {
			ID           string `xml:"id,attr"`
			Participants []struct {
				ID         string `xml:"id,attr"`
				Name       string `xml:"name,attr"`
				ProcessRef string `xml:"processRef,attr"`
			} `xml:"participant"`
		} `xml:"collaboration"`
		Processes []struct {
			ID           string `xml:"id,attr"`
			IsExecutable bool   `xml:"isExecutable,attr"`
		} `xml:"process"`
	}

	if err := xml.Unmarshal(xmlBytes, &parsed); err != nil {
		t.Fatalf("XML unmarshal failed: %v\nXML:\n%s", err, xmlStr)
	}

	if parsed.Collaboration.ID != "sre_infrastructure" {
		t.Errorf("Expected collaboration ID sre_infrastructure, got %s", parsed.Collaboration.ID)
	}
	if len(parsed.Collaboration.Participants) != 2 {
		t.Fatalf("Expected 2 participants, got %d", len(parsed.Collaboration.Participants))
	}
	if len(parsed.Processes) != 2 {
		t.Fatalf("Expected 2 processes, got %d", len(parsed.Processes))
	}

	// 3. Verify that Pool 2 is positioned below Pool 1 in BPMNDI
	type TestShape struct {
		ID           string  `xml:"id,attr"`
		BPMNElement  string  `xml:"bpmnElement,attr"`
		IsHorizontal bool    `xml:"isHorizontal,attr"`
		Bounds       struct {
			X      float64 `xml:"x,attr"`
			Y      float64 `xml:"y,attr"`
			Width  float64 `xml:"width,attr"`
			Height float64 `xml:"height,attr"`
		} `xml:"Bounds"`
	}

	type TestDiagram struct {
		BPMNPlane struct {
			BPMNElement string      `xml:"bpmnElement,attr"`
			Shapes      []TestShape `xml:"BPMNShape"`
		} `xml:"BPMNPlane"`
	}

	var parsedDoc struct {
		Diagram TestDiagram `xml:"BPMNDiagram"`
	}
	if err := xml.Unmarshal(xmlBytes, &parsedDoc); err != nil {
		t.Fatalf("Failed to unmarshal BPMNDI: %v", err)
	}

	poolYCoords := make(map[string]float64)
	for _, s := range parsedDoc.Diagram.BPMNPlane.Shapes {
		if s.BPMNElement == "pool_prod" || s.BPMNElement == "pool_test" {
			poolYCoords[s.BPMNElement] = s.Bounds.Y
		}
	}

	if prodY, ok := poolYCoords["pool_prod"]; ok {
		if testY, ok := poolYCoords["pool_test"]; ok {
			if testY <= prodY {
				t.Errorf("Expected pool_test Y (%.1f) to be strictly below pool_prod Y (%.1f)", testY, prodY)
			}
		} else {
			t.Errorf("pool_test bounds not found in BPMNDI")
		}
	} else {
		t.Errorf("pool_prod bounds not found in BPMNDI")
	}
}

func TestLayout_DisjointComponents_NoHorizontalSpanning(t *testing.T) {
	// Build a single workflow with TWO disconnected components:
	// Component 1 (Prod Cluster with Hub):
	// start_prod -> gw_hub -> branch1 -> end_prod
	// Component 2 (Test Stand, disconnected):
	// start_test -> task_test1 -> task_test2 -> end_test
	wf := nativebpm.NewWorkflow("mono_sre", "Multi-Environment SRE")
	wf.SetLayoutPreset(nativebpm.LayoutCenterHub)
	wf.SetCenterHub("gw_hub")

	// Component 1: Prod
	wf.Start("start_prod", "Prod Start").
		Parallel("gw_hub", "Argo Gateway").
		Service("task_prod1", "Node 1", "node_topic").
		End("end_prod", "Prod End")

	// Component 2: Disjoint Test Stand
	wf.Start("start_test", "Test Start").
		Service("task_test1", "CF Access Auth", "test_auth").
		Service("task_test2", "Argo Tunnel Test", "test_tunnel").
		End("end_test", "Test End")

	xmlBytes, err := wf.ToBPMNXML()
	if err != nil {
		t.Fatalf("ToBPMNXML failed: %v", err)
	}

	// Extract shape coordinates using TestDiagram
	var monoDoc struct {
		Diagram struct {
			BPMNPlane struct {
				Shapes []struct {
					BPMNElement string `xml:"bpmnElement,attr"`
					Bounds      struct {
						X float64 `xml:"x,attr"`
						Y float64 `xml:"y,attr"`
					} `xml:"Bounds"`
				} `xml:"BPMNShape"`
			} `xml:"BPMNPlane"`
		} `xml:"BPMNDiagram"`
	}
	if err := xml.Unmarshal(xmlBytes, &monoDoc); err != nil {
		t.Fatalf("Failed to unmarshal BPMNDI: %v", err)
	}

	coords := make(map[string][2]float64)
	for _, s := range monoDoc.Diagram.BPMNPlane.Shapes {
		coords[s.BPMNElement] = [2]float64{s.Bounds.X, s.Bounds.Y}
	}

	startTest := coords["start_test"]
	taskTest1 := coords["task_test1"]
	taskTest2 := coords["task_test2"]
	endTest := coords["end_test"]

	// 1. Verify start_test is placed at left margin (near startX = 150)
	if startTest[0] > 300.0 {
		t.Errorf("start_test X should be near left margin (~150), got %.1f", startTest[0])
	}

	// 2. Verify task_test1 is adjacent to start_test (not flung thousands of pixels to the right)
	flowDistance := taskTest1[0] - startTest[0]
	if flowDistance > 500.0 || flowDistance < 50.0 {
		t.Errorf("Distance between start_test and task_test1 should be normal (~240px), got %.1f (task_test1 X=%.1f, start_test X=%.1f)",
			flowDistance, taskTest1[0], startTest[0])
	}

	// 3. Verify monotonic left-to-right progression for test component
	if taskTest2[0] <= taskTest1[0] {
		t.Errorf("task_test2 X (%.1f) should be greater than task_test1 X (%.1f)", taskTest2[0], taskTest1[0])
	}
	if endTest[0] <= taskTest2[0] {
		t.Errorf("end_test X (%.1f) should be greater than task_test2 X (%.1f)", endTest[0], taskTest2[0])
	}

	// 4. Verify test component Y is positioned below prod component Y
	startProd := coords["start_prod"]
	if startTest[1] <= startProd[1] {
		t.Errorf("Expected secondary test component Y (%.1f) to be below prod component Y (%.1f)", startTest[1], startProd[1])
	}
}
