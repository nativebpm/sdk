package com.nativebpm.example;

import com.nativebpm.client.ApiClient;
import com.nativebpm.client.ApiException;
import com.nativebpm.client.api.DefaultApi;
import com.nativebpm.client.model.ProcessDefinition;
import com.nativebpm.client.model.ProcessInstance;
import com.nativebpm.client.model.StartInstanceRequest;
import com.nativebpm.client.builder.Workflow;

import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.UUID;

public class ClientExample {

    public static void main(String[] args) {
        System.out.println("🚀 NativeBPM Java SDK Example Starting...");

        // 1. Initialize ApiClient
        ApiClient client = new ApiClient();
        client.setBasePath("http://localhost:8080");
        client.addDefaultHeader("Authorization", "Bearer test-bearer-token");

        DefaultApi api = new DefaultApi(client);

        // Temp file for deployment multipart upload
        File tempFile = null;
        try {
            // Build workflow dynamically using Workflow as Code Fluent API
            System.out.println("🔨 Building workflow dynamically using Workflow as Code Fluent API...");
            Workflow workflow = new Workflow("orderProcess", "Order Processing Workflow");
            workflow.startEvent("start")
                    .next("calculateTotal");

            workflow.serviceTask("calculateTotal", "Calculate Totals", "payment-topic")
                    .wasm("./payment.wasm")
                    .next("gateway");

            workflow.exclusiveGateway("gateway", "Urgency Gateway")
                    .condition("reviewOrder", "${isUrgent == true}")
                    .defaultValue("notifyCustomer");

            workflow.userTask("reviewOrder", "Manual Urgency Review")
                    .assignee("sales_representative")
                    .next("end");

            workflow.serviceTask("notifyCustomer", "Send Notification", "notification-topic")
                    .next("end");

            workflow.endEvent("end", "End");

            String bpmnXml = workflow.buildXML();

            tempFile = File.createTempFile("simpleProcess", ".bpmn");
            try (FileWriter writer = new FileWriter(tempFile, StandardCharsets.UTF_8)) {
                writer.write(bpmnXml);
            }

            // 2. Deploy Definition
            System.out.println("📦 Deploying process definition...");
            ProcessDefinition definition = api.deployDefinition(tempFile);
            System.out.println("✅ Deployed! ID: " + definition.getId() + ", Hash: " + definition.getHash());

            // 3. Start Process Instance
            System.out.println("⚡ Starting process instance...");
            String instanceId = UUID.randomUUID().toString();
            String businessKey = "bk-java-" + System.currentTimeMillis();

            StartInstanceRequest startRequest = new StartInstanceRequest();
            startRequest.setInstanceId(UUID.fromString(instanceId));
            startRequest.setBusinessKey(businessKey);
            startRequest.putVariablesItem("started_by", "java-example-client");
            startRequest.putVariablesItem("count", 100);

            ProcessInstance instance = api.startInstance(definition.getId(), startRequest);
            System.out.println("✅ Instance started! ID: " + instance.getId() + ", Completed: " + instance.getCompleted());

            // 4. List Active Definitions
            System.out.println("📋 Listing deployed definitions...");
            List<ProcessDefinition> definitions = api.listDefinitions();
            for (ProcessDefinition def : definitions) {
                System.out.println(" - Definition: " + def.getId() + " (name: " + def.getName() + ")");
            }

            // 5. List Instances
            System.out.println("📋 Listing process instances...");
            List<ProcessInstance> instances = api.listInstances();
            for (ProcessInstance inst : instances) {
                System.out.println(" - Instance: " + inst.getId() + ", Process: " + inst.getProcessId() + ", Completed: " + inst.getCompleted());
            }

            System.out.println("🎉 Java SDK usage example finished successfully!");

        } catch (ApiException e) {
            System.err.println("❌ API Request failed! Response code: " + e.getCode());
            System.err.println("Response body: " + e.getResponseBody());
            e.printStackTrace();
        } catch (IOException e) {
            System.err.println("❌ IO error occurred!");
            e.printStackTrace();
        } finally {
            if (tempFile != null && tempFile.exists()) {
                tempFile.delete();
            }
        }
    }
}
