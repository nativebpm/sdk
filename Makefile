.PHONY: generate test
OPENAPI_GEN_IMG = registry.gitlab.com/nativebpm/sdk/openapi-generator-cli:latest
DOTNET_IMG = registry.gitlab.com/nativebpm/sdk/dotnet-sdk:9.0

test:
	cd dotnet && dotnet test NativeBPM.Client.sln

generate:
	docker run --rm -v "$$(pwd):/local" $(OPENAPI_GEN_IMG) generate \
		-i /local/api/openapi.yaml \
		-g csharp \
		-o /local/dotnet \
		--additional-properties=packageName=NativeBPM.Client,targetFramework=net9.0,hideGenerationTimestamp=true,packageGuid={5C27DC0F-7267-4E8A-8D82-A50B20450F28}
