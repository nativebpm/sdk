<?php

require_once __DIR__ . '/../vendor/autoload.php';

use NativeBPM\Client\Builder\Workflow;

echo "=== NativeBPM PHP SDK: Workflow with Guest WASM Plugins ===\n";

// 1. Build workflow dynamically using Workflow as Code Fluent API
echo "🔨 Building workflow dynamically using Fluent API...\n";
$workflow = new Workflow("wasm-demo", "Workflow with Guest WASM Plugins");

// Chain starting from start event
$workflow->startEvent("start")
    ->next("calculate");

// Configure WASM Service Task and chain to the AI task
$workflow->serviceTask("calculate", "Calculate Totals", "payment_topic")
    ->wasm("./calculate_total.wasm")
    ->next("aiCheck");

// Configure AI task and chain to the exclusive gateway
$workflow->aiTask("aiCheck", "AI Fraud Guard")
    ->provider("google")
    ->model("gemini-2.5-flash")
    ->prompt('Analyze transaction for fraud: ${orderAmount}')
    ->resultVar("isFraudulent")
    ->next("gateway");

// Set up exclusive gateway with transition conditions and default flow
$workflow->exclusiveGateway("gateway", "Fraud Gateway")
    ->condition("userTask", '${isFraudulent == true}')
    ->defaultValue("end");

// Configure User Task and chain to the end event
$workflow->userTask("userTask", "Manual Fraud Approval")
    ->assignee("security_officer")
    ->next("end");

$workflow->endEvent("end", "Process Finished");

$bpmnXml = $workflow->buildXML();
echo "✓ Successfully compiled WASM workflow AST to BPMN 2.0 XML.\n";

// Print a snippet of generated XML
if (strlen($bpmnXml) > 300) {
    echo "XML snippet:\n" . substr($bpmnXml, 0, 300) . "...\n";
} else {
    echo "XML output:\n" . $bpmnXml . "\n";
}
