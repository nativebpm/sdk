package main

import (
	"context"
	"fmt"
	"os"

	"gitlab.com/nativebpm/platform/sdk/go"
)

func main() {
	fmt.Println("=== NativeBPM Go SDK: Workflow with Native Service Tasks ===")
	ctx := context.Background()

	// 1. Build workflow as code (without WASM tasks) using Fluent API method chaining
	workflow := nativebpm.NewWorkflow("native-demo", "Workflow with Native Service Tasks")

	// Chain starting from the start event
	workflow.StartEvent("start").
		Next("gateway").
		Builder().
		ExclusiveGateway("gateway", "Urgency Gateway").
		Condition("reviewOrder", "${isUrgent == true}").
		Default("notifyCustomer")

	// Add standard Service Task. Since no `.Wasm(...)` call is made,
	// this task is executed natively as a topic-based worker queue execution.
	workflow.ServiceTask("notifyCustomer", "Send Confirmation Email", "email_topic").
		Next("end").
		Builder().
		UserTask("reviewOrder", "Review Order Details").
		Assignee("sales_representative").
		Next("end").
		Builder().
		EndEvent("end", "Process Finished")

	// Compile the workflow AST to standard BPMN 2.0 XML using the default embedded Go engine
	bpmnXML, err := workflow.BuildXML(ctx)
	if err != nil {
		fmt.Printf("Error compiling workflow: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Successfully compiled native workflow AST to BPMN 2.0 XML.")
	
	// Print a small snippet of the XML
	if len(bpmnXML) > 300 {
		fmt.Printf("XML snippet:\n%s...\n", bpmnXML[:300])
	} else {
		fmt.Printf("XML output:\n%s\n", bpmnXML)
	}
}
