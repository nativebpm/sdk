<?php

require_once __DIR__ . '/../vendor/autoload.php';

use NativeBPM\Client\Builder\Workflow;
use NativeBPM\Client\ApiClient;
use NativeBPM\Client\Api\DefaultApi;

echo "=== NativeBPM PHP SDK: Workflow as Code Builder Example ===\n";

// 1. Build workflow dynamically using Workflow as Code Fluent API
echo "🔨 Building workflow dynamically using Fluent API...\n";
$workflow = new Workflow("orderProcess", "Order Processing Workflow");
$workflow->startEvent("start")
    ->next("calculateTotal");

$workflow->serviceTask("calculateTotal", "Calculate Totals", "payment-topic")
    ->wasm("./payment.wasm")
    ->next("gateway");

$workflow->exclusiveGateway("gateway", "Urgency Gateway")
    ->condition("reviewOrder", '${isUrgent == true}')
    ->defaultValue("notifyCustomer");

$workflow->userTask("reviewOrder", "Manual Urgency Review")
    ->assignee("sales_representative")
    ->next("end");

$workflow->serviceTask("notifyCustomer", "Send Notification", "notification-topic")
    ->next("end");

$workflow->endEvent("end", "End");

$bpmnXml = $workflow->buildXML();
echo "✓ Successfully compiled workflow AST to BPMN 2.0 XML.\n";

// Print a snippet of generated XML
if (strlen($bpmnXml) > 300) {
    echo "XML snippet:\n" . substr($bpmnXml, 0, 300) . "...\n";
} else {
    echo "XML output:\n" . $bpmnXml . "\n";
}
