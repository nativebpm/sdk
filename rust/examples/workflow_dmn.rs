use nativebpm_client::builder::Workflow;
use nativebpm_client::apis::{configuration, default_api};

#[tokio::main]
async fn main() {
    println!("=== NativeBPM Rust SDK: Workflow with DMN Decision Table ===");

    // 1. Build workflow as code using Fluent API method chaining & inline DMN table definition
    let mut workflow = Workflow::new("dmn-demo", "Workflow with DMN Table");

    // Chain starting from the start event
    workflow.start_event("start")
        .connect_to("calculate-discount")
        .business_rule_task("calculate-discount", "Determine Discount", "discount-decision")
        .hit_policy("UNIQUE")
        .result_variable("discountResult")
        .input("customerType", "string")
        .input("orderAmount", "double")
        .output("discountPercent", "double")
        .rule()
            .when("customerType", serde_json::json!("VIP"))
            .when("orderAmount", serde_json::json!(">= 1000"))
            .then("discountPercent", serde_json::json!(0.2))
        .rule()
            .when("customerType", serde_json::json!("VIP"))
            .when("orderAmount", serde_json::json!("< 1000"))
            .then("discountPercent", serde_json::json!(0.15))
        .rule()
            .when("customerType", serde_json::json!("Regular"))
            .when("orderAmount", serde_json::json!(">= 500"))
            .then("discountPercent", serde_json::json!(0.1))
        .rule()
            .then("discountPercent", serde_json::json!(0.05))
        .connect_to("end")
        .end_event("end", "Process Finished");

    // Compile the workflow AST to standard BPMN 2.0 XML using the embedded Go engine
    let bpmn_xml = match workflow.build_xml() {
        Ok(xml) => {
            println!("✓ Successfully compiled DMN workflow AST to BPMN 2.0 XML.");
            xml
        }
        Err(e) => {
            eprintln!("✗ Failed to compile workflow: {}", e);
            return;
        }
    };

    // 2. Deploy and start process definition using the REST API client
    let mut config = configuration::Configuration::new();
    config.base_path = "http://localhost:8080".to_string();
    config.bearer_access_token = Some("test-bearer-token".to_string());

    println!("\nDeploying to NativeBPM engine...");
    
    let bpmn_path = std::env::temp_dir().join("dmn-demo.bpmn");
    if std::fs::write(&bpmn_path, bpmn_xml).is_ok() {
        match default_api::deploy_definition(&config, Some(bpmn_path)).await {
            Ok(definition) => {
                println!("✓ Deployed process definition (hash: {:?})", definition.hash);

                // Start a process instance
                let mut variables = std::collections::HashMap::new();
                variables.insert("customerType".to_string(), serde_json::json!("VIP"));
                variables.insert("orderAmount".to_string(), serde_json::json!(1250.0));

                let start_request = nativebpm_client::models::StartInstanceRequest {
                    instance_id: None,
                    business_key: Some("order-9981".to_string()),
                    variables: Some(variables),
                };

                match default_api::start_instance(&config, "dmn-demo", Some(start_request)).await {
                    Ok(instance) => {
                        println!("✓ Started process instance ID: {} (completed: {})", instance.id, instance.completed);
                    }
                    Err(e) => {
                        println!("Note: Failed to start process instance: {:?}", e);
                    }
                }
            }
            Err(e) => {
                println!("Note: Local API Engine deployment skipped (ensure local server is running on :8080). Details: {:?}", e);
            }
        }
    }
}
