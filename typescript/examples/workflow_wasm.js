import { Workflow, Client } from '../dist/index.js';

export async function run() {
  console.log("=== NativeBPM TS SDK: Workflow as Code ===");

  // 1. Build workflow as code (without WASM tasks)
  const workflow = new Workflow('native-demo', 'Workflow as Code');
  workflow.startEvent('start');
  
  // Add a standard Service Task. Since no `.wasm(...)` call is made,
  // this task is executed natively as a topic-based worker queue execution.
  workflow.serviceTask('notifyCustomer', 'Send Confirmation Email', 'email_topic');
  
  // Add an exclusive gateway to show transition conditions
  workflow.exclusiveGateway('gateway', 'Urgency Gateway');
  
  // Add a User (human) task for standard workflow steps
  workflow.userTask('reviewOrder', 'Review Order Details')
    .assignee('sales_representative');
      
  workflow.endEvent('end', 'Process Finished');
  
  // Connect elements with sequence flows
  workflow.sequenceFlow('start', 'gateway');
  
  // transition conditions configured as process expressions (e.g. ${isUrgent == true})
  workflow.sequenceFlowWithCondition('gateway', 'reviewOrder', '${isUrgent == true}');
  workflow.sequenceFlowWithCondition('gateway', 'notifyCustomer', '${isUrgent == false}');
  
  workflow.sequenceFlow('notifyCustomer', 'end');
  workflow.sequenceFlow('reviewOrder', 'end');
  
  // Compile the workflow AST to standard BPMN 2.0 XML using the embedded Go engine
  const bpmnXML = await workflow.buildXML();
  console.log("✓ Successfully compiled native workflow AST to BPMN 2.0 XML.");
  
  // 2. Deploy and start process definition using the Fluent Client API
  const client = new Client("http://localhost:8080", "test-bearer-token");
  
  console.log("\nDeploying to NativeBPM engine...");
  try {
    // Deploy process definition
    const definition = await client.definitions().deploy()
      .withID("native-demo")
      .withName("Workflow as Code")
      .withBPMN(bpmnXML)
      .send();
    console.log(`✓ Deployed process definition (hash: ${definition.hash})`);
    
    // Start a process instance with input variables
    const instance = await client.instances().start("native-demo")
      .withBusinessKey("order-5541")
      .withVariable("customerEmail", "customer@example.com")
      .withVariable("isUrgent", true)
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
