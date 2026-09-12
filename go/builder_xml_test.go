package nativebpm_test

import (
	"encoding/xml"
	"strings"
	"testing"

	nativebpm "gitlab.com/nativebpm/sdk/go"
)

func TestWorkflow_ToBPMNXML(t *testing.T) {
	wf := nativebpm.NewWorkflow("order-processing", "Order Processing Process")
	wf.Start("start_event").
		Service("validate_order", "Validate Order", "order_topic").
		When(nativebpm.V("is_valid").Eq(true)).
		Then(func(flow *nativebpm.Branch) {
			flow.Service("process_payment", "Process Payment", "payment_topic").
				Service("ship_order", "Ship Order", "shipping_topic")
		}).
		Else(func(flow *nativebpm.Branch) {
			flow.Service("notify_failure", "Notify Customer", "notification_topic")
		}).
		End("end_event", "Order Finished")

	xmlBytes, err := wf.ToBPMNXML()
	if err != nil {
		t.Fatalf("ToBPMNXML failed: %v", err)
	}
	if len(xmlBytes) == 0 {
		t.Fatal("ToBPMNXML returned empty output")
	}

	xmlStr := string(xmlBytes)

	// Verify well-formed XML
	var parsed struct {
		XMLName xml.Name `xml:"definitions"`
		ID      string   `xml:"id,attr"`
	}
	if err := xml.Unmarshal(xmlBytes, &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v\nXML:\n%s", err, xmlStr)
	}

	// Verify essential BPMN tags
	for _, tag := range []string{
		"<definitions",
		"<process",
		"<startEvent id=\"start_event\"",
		"<serviceTask id=\"validate_order\"",
		"<exclusiveGateway",
		"<sequenceFlow",
		"<bpmndi:BPMNDiagram",
		"<bpmndi:BPMNPlane",
		"<bpmndi:BPMNShape",
		"<bpmndi:BPMNEdge",
		"<omgdi:waypoint",
	} {
		if !strings.Contains(xmlStr, tag) {
			t.Errorf("Expected tag/substring '%s' in generated BPMN XML", tag)
		}
	}
}

func TestWorkflow_LayoutCenterHub_And_Color(t *testing.T) {
	wf := nativebpm.NewWorkflow("net_hub", "Network Hub Topology")
	wf.SetLayoutPreset(nativebpm.LayoutCenterHub).
		SetCenterHub("gw_hub")

	wf.StartEvent("client_start").
		ServiceTask("cf_edge", "Cloudflare Edge", "cf_edge").
		ParallelGateway("gw_hub", "Argo Smart Hub").
		ServiceTask("worker_1", "Worker 1", "worker").
		ServiceTask("worker_2", "Worker 2", "worker").
		EndEvent("end_sink", "Response Delivered")

	wf.SequenceFlow("client_start", "cf_edge").
		SequenceFlow("cf_edge", "gw_hub").
		SequenceFlow("gw_hub", "worker_1").
		SequenceFlow("gw_hub", "worker_2").
		SequenceFlow("worker_1", "end_sink").
		SequenceFlow("worker_2", "end_sink")

	// Set colors
	wf.SetNodeColor("cf_edge", "#0ea5e9", "rgba(14,165,233,0.15)").
		SetNodeColor("gw_hub", "#f59e0b", "rgba(245,158,11,0.15)").
		SetNodeColor("worker_1", "#10b981", "rgba(16,185,129,0.15)")

	xmlBytes, err := wf.ToBPMNXML()
	if err != nil {
		t.Fatalf("ToBPMNXML failed: %v", err)
	}

	xmlStr := string(xmlBytes)
	if !strings.Contains(xmlStr, "xmlns:bioc=\"http://bpmn.io/schema/bpmn/biocolor/1.0\"") {
		t.Error("Expected xmlns:bioc in definitions")
	}
	if !strings.Contains(xmlStr, "bioc:stroke=\"#0ea5e9\"") {
		t.Error("Expected bioc:stroke on cf_edge shape")
	}
	if !strings.Contains(xmlStr, "bioc:stroke=\"#f59e0b\"") {
		t.Error("Expected bioc:stroke on gw_hub shape")
	}
}
