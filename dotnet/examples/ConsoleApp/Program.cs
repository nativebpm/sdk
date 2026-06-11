using System;
using System.Collections.Generic;
using System.IO;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Text;
using System.Threading.Tasks;
using Microsoft.Extensions.DependencyInjection;
using NativeBPM.Client.Api;
using NativeBPM.Client.Client;
using NativeBPM.Client.Extensions;
using NativeBPM.Client.Model;

namespace ConsoleApp
{
    class Program
    {
        static async Task Main(string[] args)
        {
            Console.WriteLine("🚀 NativeBPM .NET SDK Example Starting...");

            var services = new ServiceCollection();

            string hostUrl = "http://localhost:8080";
            string apiToken = "test-bearer-token";

            // Configure NativeBPM client API with Base Address and Bearer Token
            services.AddApi(options =>
            {
                options.AddApiHttpClients(client =>
                {
                    client.BaseAddress = new Uri(hostUrl);
                    client.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", apiToken);
                });
            });

            var serviceProvider = services.BuildServiceProvider();
            var api = serviceProvider.GetRequiredService<IDefaultApi>();

            string tempFilePath = Path.GetTempFileName() + ".bpmn";
            try
            {
                // 1. Build workflow as code using the Fluent API builder
                Console.WriteLine("🔨 Building workflow dynamically using Workflow as Code Fluent API...");
                using var workflow = new Workflow("orderProcess", "Order Processing Workflow");
                
                workflow.StartEvent("start")
                    .Next("calculateTotal");
                
                // WASM-based execution logic
                workflow.ServiceTask("calculateTotal", "Calculate Totals", "payment-topic")
                    .Wasm("./payment.wasm")
                    .Next("gateway");

                // Condition-based routing in native BPMN
                workflow.ExclusiveGateway("gateway", "Urgency Gateway")
                    .Condition("reviewOrder", "${isUrgent == true}")
                    .DefaultValue("notifyCustomer");

                // Human-in-the-loop user task
                workflow.UserTask("reviewOrder", "Manual Urgency Review")
                    .Assignee("sales_representative")
                    .Next("end");

                // Native topic-based service task
                workflow.ServiceTask("notifyCustomer", "Send Notification", "notification-topic")
                    .Next("end");

                workflow.EndEvent("end", "End");

                string bpmnXml = workflow.BuildXml();

                await File.WriteAllTextAsync(tempFilePath, bpmnXml, Encoding.UTF8);

                // 2. Deploy Definition
                Console.WriteLine("📦 Deploying process definition...");
                using var fileStream = File.OpenRead(tempFilePath);
                var fileParam = new FileParameter(fileStream, "simpleProcess.bpmn", "application/octet-stream");
                var definitionResponse = await api.DeployDefinitionAsync(fileParam);
                
                if (definitionResponse.IsSuccessStatusCode && definitionResponse.TryOk(out var definition) && definition != null)
                {
                    Console.WriteLine($"✅ Deployed! ID: {definition.Id}, Hash: {definition.Hash}");

                    // 3. Start Process Instance
                    Console.WriteLine("⚡ Starting process instance...");
                    var instanceId = Guid.NewGuid();
                    var businessKey = "bk-dotnet-" + DateTimeOffset.UtcNow.ToUnixTimeMilliseconds();

                    var startRequest = new StartInstanceRequest(
                        instanceId: instanceId,
                        businessKey: businessKey,
                        variables: new Dictionary<string, object>
                        {
                            { "started_by", "dotnet-example-client" },
                            { "count", 100 }
                        }
                    );

                    var startResponse = await api.StartInstanceAsync(definition.Id, startRequest);
                    if (startResponse.IsSuccessStatusCode && startResponse.TryOk(out var instance) && instance != null)
                    {
                        Console.WriteLine($"✅ Instance started! ID: {instance.Id}, Completed: {instance.Completed}");
                    }
                    else
                    {
                        Console.WriteLine("❌ Failed to parse start instance response.");
                    }
                }
                else
                {
                    Console.WriteLine($"❌ Failed to deploy definition. Status: {definitionResponse.StatusCode}");
                }

                // 4. List Active Definitions
                Console.WriteLine("📋 Listing deployed definitions...");
                var definitionsResponse = await api.ListDefinitionsAsync();
                if (definitionsResponse.IsSuccessStatusCode && definitionsResponse.TryOk(out var definitions) && definitions != null)
                {
                    foreach (var def in definitions)
                    {
                        Console.WriteLine($" - Definition: {def.Id} (name: {def.Name})");
                    }
                }

                // 5. List Instances
                Console.WriteLine("📋 Listing process instances...");
                var instancesResponse = await api.ListInstancesAsync();
                if (instancesResponse.IsSuccessStatusCode && instancesResponse.TryOk(out var instances) && instances != null)
                {
                    foreach (var inst in instances)
                    {
                        Console.WriteLine($" - Instance: {inst.Id}, Process: {inst.ProcessId}, Completed: {inst.Completed}");
                    }
                }

                Console.WriteLine("🎉 .NET SDK usage example finished successfully!");
            }
            catch (ApiException e)
            {
                Console.WriteLine($"❌ API Request failed! Response code: {e.ErrorCode}");
                Console.WriteLine($"Response headers: {e.Headers}");
                Console.WriteLine($"Response content: {e.ErrorContent}");
            }
            catch (Exception e)
            {
                Console.WriteLine($"❌ Error occurred: {e.Message}");
                Console.WriteLine(e.StackTrace);
            }
            finally
            {
                if (File.Exists(tempFilePath))
                {
                    File.Delete(tempFilePath);
                }
            }
        }
    }
}
