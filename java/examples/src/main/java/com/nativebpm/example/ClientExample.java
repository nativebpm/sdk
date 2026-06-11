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

        // Temp files for deployment multipart upload
        File tempFile1 = null;
        File tempFile2 = null;
        try {
            // SCENARIO 1: Workflow as Code (Without custom Guest WASM tasks)
            System.out.println("--------------------------------------------------");
            String bpmnXml1 = WorkflowWasmExample.buildWorkflow();
            tempFile1 = File.createTempFile("nativeProcess", ".bpmn");
            try (FileWriter writer = new FileWriter(tempFile1, StandardCharsets.UTF_8)) {
                writer.write(bpmnXml1);
            }

            System.out.println("📦 Deploying native process definition...");
            ProcessDefinition def1 = api.deployDefinition(tempFile1);
            System.out.println("✅ Deployed! ID: " + def1.getId() + ", Hash: " + def1.getHash());

            System.out.println("⚡ Starting native process instance...");
            StartInstanceRequest startRequest1 = new StartInstanceRequest();
            startRequest1.setInstanceId(UUID.randomUUID());
            startRequest1.setBusinessKey("bk-java-native-" + System.currentTimeMillis());
            startRequest1.putVariablesItem("isUrgent", true);

            ProcessInstance inst1 = api.startInstance(def1.getId(), startRequest1);
            System.out.println("✅ Instance started! ID: " + inst1.getId() + ", Completed: " + inst1.getCompleted());

            // SCENARIO 2: Workflow with Guest WASM Plugins
            System.out.println("--------------------------------------------------");
            String bpmnXml2 = WorkflowWithWasmPluginsExample.buildWorkflow();
            tempFile2 = File.createTempFile("wasmProcess", ".bpmn");
            try (FileWriter writer = new FileWriter(tempFile2, StandardCharsets.UTF_8)) {
                writer.write(bpmnXml2);
            }

            System.out.println("📦 Deploying WASM process definition...");
            ProcessDefinition def2 = api.deployDefinition(tempFile2);
            System.out.println("✅ Deployed! ID: " + def2.getId() + ", Hash: " + def2.getHash());

            System.out.println("⚡ Starting WASM process instance...");
            StartInstanceRequest startRequest2 = new StartInstanceRequest();
            startRequest2.setInstanceId(UUID.randomUUID());
            startRequest2.setBusinessKey("bk-java-wasm-" + System.currentTimeMillis());
            startRequest2.putVariablesItem("orderAmount", 2500);

            ProcessInstance inst2 = api.startInstance(def2.getId(), startRequest2);
            System.out.println("✅ Instance started! ID: " + inst2.getId() + ", Completed: " + inst2.getCompleted());

            System.out.println("--------------------------------------------------");

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
            if (tempFile1 != null && tempFile1.exists()) {
                tempFile1.delete();
            }
            if (tempFile2 != null && tempFile2.exists()) {
                tempFile2.delete();
            }
        }
    }
}
