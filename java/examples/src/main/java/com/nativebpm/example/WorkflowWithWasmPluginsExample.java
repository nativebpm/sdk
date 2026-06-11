package com.nativebpm.example;

import com.nativebpm.client.builder.Workflow;

public class WorkflowWithWasmPluginsExample {
    public static String buildWorkflow() throws Exception {
        System.out.println("=== NativeBPM Java SDK: Workflow with Guest WASM Plugins ===");

        // 1. Build workflow dynamically using Workflow as Code Fluent API
        System.out.println("🔨 Building workflow dynamically using Fluent API...");
        Workflow workflow = new Workflow("wasm-demo", "Workflow with Guest WASM Plugins");

        // Chain starting from start event
        workflow.startEvent("start")
                .next("calculate");

        // Configure WASM Service Task and chain to the AI task
        workflow.serviceTask("calculate", "Calculate Totals", "payment_topic")
                .wasm("./calculate_total.wasm")
                .next("aiCheck");

        // Configure AI task and chain to the exclusive gateway
        workflow.aiTask("aiCheck", "AI Fraud Guard")
                .provider("google")
                .model("gemini-2.5-flash")
                .prompt("Analyze transaction for fraud: ${orderAmount}")
                .resultVar("isFraudulent")
                .next("gateway");

        // Set up exclusive gateway with transition conditions and default flow
        workflow.exclusiveGateway("gateway", "Fraud Gateway")
                .condition("userTask", "${isFraudulent == true}")
                .defaultValue("end");

        // Configure User Task and chain to the end event
        workflow.userTask("userTask", "Manual Fraud Approval")
                .assignee("security_officer")
                .next("end");

        workflow.endEvent("end", "Process Finished");

        String bpmnXml = workflow.buildXML();
        System.out.println("✓ Successfully compiled WASM workflow AST to BPMN 2.0 XML.");
        return bpmnXml;
    }
}
