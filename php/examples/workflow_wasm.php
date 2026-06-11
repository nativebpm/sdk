<?php

require_once __DIR__ . '/../vendor/autoload.php';

use NativeBPM\Client\Builder\Workflow;

echo "=== NativeBPM PHP SDK: Workflow as Code ===\n";

// 1. Build workflow dynamically using Workflow as Code Fluent API
echo "🔨 Building workflow dynamically using Fluent API...\n";
$workflow = new Workflow("native-demo", "Workflow as Code");

// Chain starting from the start event
$workflow->startEvent("start")
    ->next("gateway");

// Configure exclusive gateway with transition conditions and default flow
$workflow->exclusiveGateway("gateway", "Urgency Gateway")
    ->condition("reviewOrder", '${isUrgent == true}')
    ->defaultValue("notifyCustomer");

// Add standard Service Task. Since no `->wasm(...)` call is made,
// this task is executed natively as a topic-based worker queue execution.
$workflow->serviceTask("notifyCustomer", "Send Confirmation Email", "email_topic")
    ->next("end");

// Add User Task and chain to the end event
$workflow->userTask("reviewOrder", "Review Order Details")
    ->assignee("sales_representative")
    ->next("end");

$workflow->endEvent("end", "Process Finished");

$bpmnXml = $workflow->buildXML();
echo "✓ Successfully compiled native workflow AST to BPMN 2.0 XML.\n";

// Print a snippet of generated XML
if (strlen($bpmnXml) > 300) {
    echo "XML snippet:\n" . substr($bpmnXml, 0, 300) . "...\n";
} else {
    echo "XML output:\n" . $bpmnXml . "\n";
}
