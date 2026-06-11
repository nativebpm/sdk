package com.nativebpm.example;

import com.nativebpm.client.builder.Workflow;

public class WorkflowWasmExample {
    public static String buildWorkflow() throws Exception {
        System.out.println("=== NativeBPM Java SDK: Workflow as Code ===");

        // 1. Build workflow dynamically using Workflow as Code Fluent API
        System.out.println("🔨 Building workflow dynamically using Fluent API...");
        Workflow workflow = new Workflow("native-demo", "Workflow as Code");

        // Chain starting from the start event
        workflow.startEvent("start")
                .next("gateway");

        // Configure exclusive gateway with transition conditions and default flow
        workflow.exclusiveGateway("gateway", "Urgency Gateway")
                .condition("reviewOrder", "${isUrgent == true}")
                .defaultValue("notifyCustomer");

        // Add standard Service Task. Since no .wasm(...) call is made,
        // this task is executed natively as a topic-based worker queue execution.
        workflow.serviceTask("notifyCustomer", "Send Confirmation Email", "email_topic")
                .next("end");

        // Add User Task and chain to the end event
        workflow.userTask("reviewOrder", "Review Order Details")
                .assignee("sales_representative")
                .next("end");

        workflow.endEvent("end", "Process Finished");

        String bpmnXml = workflow.buildXML();
        System.out.println("✓ Successfully compiled native workflow AST to BPMN 2.0 XML.");
        return bpmnXml;
    }
}
