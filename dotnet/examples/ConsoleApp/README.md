[← Back to .NET SDK](../../README.md) | [← Back to Platform Root](../../../../README.md)

# NativeBPM .NET SDK Usage Example

This example demonstrates how to integrate the generated .NET Client SDK into a .NET Console Application using Microsoft Dependency Injection.

## Prerequisites

- .NET 8.0 SDK installed on your system.
- A running NativeBPM engine (typically at `http://localhost:8080`).

## Project Configuration

The project references the local SDK client library by using a `<ProjectReference>` in the `.csproj` file.

## Running the Example

Make sure the NativeBPM server is running, then execute the following command:

```bash
dotnet run
```

This will run the program, which executes the following steps:
1. Registers the `IDefaultApi` using Microsoft's `IServiceCollection` extension `AddApi`.
2. Configures base address to `http://localhost:8080` and adds a Bearer Authorization token.
3. Automatically creates a temporary BPMN XML file.
4. Deploys the workflow definition using `DeployDefinitionAsync`.
5. Starts a process instance of the deployed definition using `StartInstanceAsync`.
6. Queries and lists all deployed process definitions using `ListDefinitionsAsync`.
7. Queries and lists all process instances using `ListInstancesAsync`.
