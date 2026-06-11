import { Workflow, Client } from '../dist/index.js';

export async function run() {
  console.log("=== NativeBPM TS SDK: Workflow with WASM Service Task ===");

  // 1. Build workflow as code
  const workflow = new Workflow('wasm-demo', 'Workflow with WASM Service Task', './src/core.wasm');
  workflow.startEvent('start');
  
  // Add a Service Task configured with a WASM file path/alias.
  // The process engine will execute this compiled WASM binary inside a secure WASM sandbox.
  workflow.serviceTask('calculate', 'Calculate Totals', 'payment_topic')
    .wasm('./calculate_total.wasm');
  
  // Add an AI task that delegates decision checking to an LLM
  workflow.aiTask('aiCheck', 'AI Fraud Guard')
    .provider('google')
    .model('gemini-2.5-flash')
    .prompt('Analyze transaction for fraud: ${orderAmount}')
    .resultVar('isFraudulent');
      
  workflow.exclusiveGateway('gateway', 'Fraud Gateway');
  
  // Add a User (human) task for high-risk approval
  workflow.userTask('userTask', 'Manual Fraud Approval')
    .assignee('security_officer');
      
  workflow.endEvent('end', 'Process Finished');
  
  // Connect elements with sequence flows
  workflow.sequenceFlow('start', 'calculate');
  workflow.sequenceFlow('calculate', 'aiCheck');
  workflow.sequenceFlow('aiCheck', 'gateway');
  
  // Sequence flows with conditions
  workflow.sequenceFlowWithCondition('gateway', 'userTask', '${isFraudulent == true}');
  workflow.sequenceFlowWithCondition('gateway', 'end', '${isFraudulent == false}');
  workflow.sequenceFlow('userTask', 'end');
  
  // Compile the workflow AST to standard BPMN 2.0 XML using the embedded Go engine
  const bpmnXML = await workflow.buildXML();
  console.log("✓ Successfully compiled WASM workflow AST to BPMN 2.0 XML.");
  
  // 2. Deploy and start process definition using the Fluent Client API
  const client = new Client("http://localhost:8080", "test-bearer-token");
  
  console.log("\nDeploying to NativeBPM engine...");
  try {
    // Deploy process definition
    const definition = await client.definitions().deploy()
      .withID("wasm-demo")
      .withName("Workflow with WASM Service Task")
      .withBPMN(bpmnXML)
      .send();
    console.log(`✓ Deployed process definition (hash: ${definition.hash})`);
    
    // Start a process instance with input variables
    const instance = await client.instances().start("wasm-demo")
      .withBusinessKey("tx-8837")
      .withVariable("orderAmount", 2500)
      .send();
    console.log(`✓ Started process instance ID: ${instance.id} (completed: ${instance.completed})`);
    
  } catch (error) {
    console.log(`Note: Local API Engine deployment skipped (ensure local server is running on :8080). Details: ${error.message || error}`);
  }
}

// If run directly
if (import.meta.url.endsWith(process.argv[1])) {
  run();
}
