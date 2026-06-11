import os
from nativebpm import Client, Workflow

def main():
    print("=== NativeBPM Python SDK: Workflow with Native Service Tasks ===")

    # 1. Build workflow as code (without WASM tasks) using Fluent API method chaining
    workflow = Workflow('native-demo', 'Workflow with Native Service Tasks')
    
    # Chain starting from the start event
    workflow.start_event('start').next('gateway')
    
    # Configure exclusive gateway with transition conditions and default flow
    workflow.exclusive_gateway('gateway', 'Urgency Gateway')\
        .condition('reviewOrder', '${isUrgent == true}')\
        .default('notifyCustomer')
    
    # Add standard Service Task. Since no `.wasm(...)` call is made,
    # this task is executed natively as a topic-based worker queue execution.
    workflow.service_task('notifyCustomer', 'Send Confirmation Email', 'email_topic')\
        .next('end')
    
    # Add User Task and chain to the end event
    workflow.user_task('reviewOrder', 'Review Order Details')\
        .assignee('sales_representative')\
        .next('end')
        
    workflow.end_event('end', 'Process Finished')
    
    # Compile the workflow AST to standard BPMN 2.0 XML using the embedded Go engine
    bpmn_xml = workflow.build_xml()
    print("✓ Successfully compiled native workflow AST to BPMN 2.0 XML.")
    
    # 2. Deploy and start process definition using the Fluent Client API
    client = Client("http://localhost:8080", "test-bearer-token")
    
    print("\nDeploying to NativeBPM engine...")
    try:
        # Deploy process definition
        definition = client.definitions().deploy()\
            .with_id("native-demo")\
            .with_name("Workflow with Native Service Tasks")\
            .with_bpmn(bpmn_xml.encode('utf-8'))\
            .send()
        print(f"✓ Deployed process definition (hash: {definition.hash})")
        
        # Start a process instance with input variables
        instance = client.instances().start("native-demo")\
            .with_business_key("order-5541")\
            .with_variable("customerEmail", "customer@example.com")\
            .with_variable("isUrgent", True)\
            .send()
        print(f"✓ Started process instance ID: {instance.id} (completed: {instance.completed})")
        
    except Exception as e:
        print(f"Note: Local API Engine deployment skipped (ensure local server is running on :8080). Details: {e}")

if __name__ == '__main__':
    main()
