[← Back to Platform Root](../../README.md)

# NativeBPM .NET SDK

Official .NET/C# Client SDK for NativeBPM Cloud-Native engine. Built using `HttpClient` for modern asynchronous HTTP transport.

## Building the SDK

To compile the SDK project locally:

```bash
dotnet build
```

To package the SDK into a NuGet package (`.nupkg`):

```bash
dotnet pack -c Release -o ./out
```

## Installation

To consume the package from the GitLab NuGet Package Registry:

1. Add the NuGet source (replace `<your_gitlab_token>` if repository requires auth; for public reading, the password can be any token or skipped if using public download endpoints):
```bash
dotnet nuget add source "https://gitlab.com/api/v4/projects/nativebpm%2Fsdk/packages/nuget/index.json" \
  --name gitlab \
  --username gitlab-ci-token \
  --password your_personal_access_or_deploy_token \
  --store-password-in-clear-text
```

2. Add the dependency to your project:
```bash
dotnet add package NativeBPM.Client
```

## Namespace & Imports

The generated classes reside in the `NativeBPM.Client` namespaces:

```csharp
using NativeBPM.Client.Api;
using NativeBPM.Client.Client;
using NativeBPM.Client.Model;
```

## Usage Example

Initialize `HostConfiguration` and use the dependency injection-friendly `ApiFactory` to instantiate clients:

```csharp
using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using NativeBPM.Client.Api;
using NativeBPM.Client.Client;
using NativeBPM.Client.Model;

class Program
{
    static async Task Main(string[] args)
    {
        // 1. Configure HttpClient and host properties
        var config = new HostConfiguration<DefaultApi>
        {
            BasePath = "http://localhost:8080",
            Token = "your-api-token" // Automatically sent as Bearer token
        };

        // 2. Instantiate API factory
        var factory = new ApiFactory(config);
        var api = factory.GetApi<DefaultApi>();

        try
        {
            // 3. List deployed process definitions
            List<ProcessDefinition> definitions = await api.ListDefinitionsAsync();
            foreach (var def in definitions)
            {
                Console.WriteLine($"Definition: {def.Name} (ID: {def.Id})");
            }

            // 4. Start a new process instance with business key and variables
            var startRequest = new StartInstanceRequest
            {
                BusinessKey = "REQ-8802",
                Variables = new Dictionary<string, object>
                {
                    { "limit", 15000 },
                    { "escalate", false }
                }
            };

            ProcessInstance instance = await api.StartInstanceAsync("expense-claims-process", startRequest);
            Console.WriteLine($"Started instance: {instance.Id}");
        }
        catch (ApiException e)
        {
            Console.WriteLine($"API Error Code: {e.ErrorCode}");
            Console.WriteLine($"API Error Message: {e.Message}");
        }
    }
}
```

## Developer Guide: Publishing to GitLab Package Registry

### Automated Publishing via CI/CD
This package is automatically built, versioned, and published to the GitLab NuGet Package Registry whenever a git tag matching the pattern `sdk/dotnet/v*` is pushed. For example:
```bash
git tag sdk/dotnet/v1.0.0
git push origin sdk/dotnet/v1.0.0
```

### Manual Publishing
To manually publish a NuGet release package:
1. Package the project locally with the release version:
   ```bash
   dotnet pack -c Release -p:PackageVersion=1.0.0 -o out src/NativeBPM.Client/NativeBPM.Client.csproj
   ```
2. Add the GitLab source and push:
   ```bash
   dotnet nuget add source "https://gitlab.com/api/v4/projects/nativebpm%2Fsdk/packages/nuget/index.json" \
     --name gitlab \
     --username gitlab-ci-token \
     --password your_personal_access_or_deploy_token \
     --store-password-in-clear-text

   dotnet nuget push "out/*.nupkg" --source gitlab --api-key your_personal_access_or_deploy_token
   ```