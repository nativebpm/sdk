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
