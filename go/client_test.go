package nativebpm

import (
	"context"
	"encoding/json"
	"testing"
)

func TestWorkflowBuilderJSON(t *testing.T) {
	wf := NewWorkflow("test-process", "Test Process")
	wf.Service("service1", "Service 1", "topic1", M{"wasm": "./my_wasm_module.wasm"}).
		End("end", "End")

	jsonData, err := wf.ToJSON()
	if err != nil {
		t.Fatalf("failed to serialize to JSON: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if parsed["id"] != "test-process" {
		t.Errorf("expected id test-process, got %v", parsed["id"])
	}

	nodes, ok := parsed["nodes"].([]interface{})
	if !ok || len(nodes) < 3 {
		t.Fatalf("expected at least 3 nodes, got %d", len(nodes))
	}

	// Verify service1 node has wasmPath
	var foundService bool
	for _, n := range nodes {
		m, ok := n.(map[string]interface{})
		if !ok {
			continue
		}
		if m["id"] == "service1" {
			foundService = true
			if m["wasmPath"] != "./my_wasm_module.wasm" {
				t.Errorf("expected wasmPath ./my_wasm_module.wasm, got %v", m["wasmPath"])
			}
		}
	}
	if !foundService {
		t.Errorf("service1 node not found in nodes list")
	}
}

func TestBlockClosureDSL(t *testing.T) {
	wf := NewWorkflow("closure-process", "Closure Process")
	wf.
		User("task1", "User Approval").
		When(V("approved").Eq(true)).
		Then(func(flow *Branch) {
			flow.Service("publish", "Publish Page", "publish-topic", M{"wasm": "./publish.wasm"})
		}).
		Else(func(flow *Branch) {
			flow.Service("reject", "Notify Reject", "reject-topic")
		})

	jsonData, err := wf.ToJSON()
	if err != nil {
		t.Fatalf("failed to serialize to JSON: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	flows, ok := parsed["flows"].([]interface{})
	if !ok || len(flows) == 0 {
		t.Fatalf("expected flows, got none")
	}

	var foundCondition bool
	for _, f := range flows {
		m, ok := f.(map[string]interface{})
		if !ok {
			continue
		}
		if m["condition"] == "approved == true" {
			foundCondition = true
		}
	}
	if !foundCondition {
		t.Errorf("expected to find flow condition 'approved == true'")
	}
}

func TestClient_NewWorker(t *testing.T) {
	client, err := NewClient("unix:///tmp/nbpm.sock", "token-abc")
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	worker := client.NewWorker("inventory_check", func(ctx context.Context, task *TaskContext) (map[string]interface{}, error) {
		return map[string]interface{}{"available": true}, nil
	})

	if worker == nil {
		t.Fatal("expected non-nil worker")
	}
	if worker.serverURL != "unix:///tmp/nbpm.sock" {
		t.Errorf("expected serverURL unix:///tmp/nbpm.sock, got %s", worker.serverURL)
	}
	if worker.apiToken != "token-abc" {
		t.Errorf("expected apiToken token-abc, got %s", worker.apiToken)
	}
	topics := worker.getTopicList()
	if len(topics) != 1 || topics[0] != "inventory_check" {
		t.Errorf("expected topics [inventory_check], got %v", topics)
	}
}
