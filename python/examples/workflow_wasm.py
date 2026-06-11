import os
from nativebpm import Client, Workflow

def main():
    print("=== NativeBPM Python SDK: Workflow with WASM Service Task ===")

    # 1. Build workflow as code using Fluent API method chaining
    workflow = Workflow('wasm-demo', 'Workflow with WASM Service Task', './nativebpm/core.wasm')
    
    # Chain starting from the start event
    workflow.start_event('start').next('calculate')
    
    # Configure WASM Service Task and chain to the AI task
    workflow.service_task('calculate', 'Calculate Totals', 'payment_topic')\
        .wasm('./calculate_total.wasm')\
        .next('aiCheck')
    
    # Configure AI task and chain to the exclusive gateway
    workflow.ai_task('aiCheck', 'AI Fraud Guard')\
        .provider('google')\
        .model('gemini-2.5-flash')\
        .prompt('Analyze transaction for fraud: ${orderAmount}')\
        .result_var('isFraudulent')\
        .next('gateway')
        
    # Set up exclusive gateway with transition conditions and default flow
    workflow.exclusive_gateway('gateway', 'Fraud Gateway')\
        .condition('userTask', '${isFraudulent == true}')\
        .default('end')
    
    # Configure User Task and chain to the end event
    workflow.user_task('userTask', 'Manual Fraud Approval')\
        .assignee('security_officer')\
        .next('end')
        
    workflow.end_event('end', 'Process Finished')
    
    # Compile the workflow AST to standard BPMN 2.0 XML using the embedded Go engine
    bpmn_xml = workflow.build_xml()
    print("✓ Successfully compiled WASM workflow AST to BPMN 2.0 XML.")
    
    # 2. Deploy and start process definition using the Fluent Client API
    client = Client("http://localhost:8080", "test-bearer-token")
    
    print("\nDeploying to NativeBPM engine...")
    try:
        # Deploy process definition
        definition = client.definitions().deploy()\
            .with_id("wasm-demo")\
            .with_name("Workflow with WASM Service Task")\
            .with_bpmn(bpmn_xml.encode('utf-8'))\
            .send()
        print(f"✓ Deployed process definition (hash: {definition.hash})")
        
        # Start a process instance with input variables
        instance = client.instances().start("wasm-demo")\
            .with_business_key("tx-8837")\
            .with_variable("orderAmount", 2500)\
            .send()
        print(f"✓ Started process instance ID: {instance.id} (completed: {instance.completed})")
        
    except Exception as e:
        print(f"Note: Local API Engine deployment skipped (ensure local server is running on :8080). Details: {e}")

if __name__ == '__main__':
    main()
