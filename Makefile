.PHONY: generate test
OPENAPI_GEN_IMG = registry.gitlab.com/nativebpm/sdk/openapi-generator-cli:latest

test:
	cd swift && swift test

generate:
	docker run --rm -v "$$(pwd):/local" $(OPENAPI_GEN_IMG) generate \
		-i /local/api/openapi.yaml \
		-g swift5 \
		-o /local/swift \
		--additional-properties=projectName=NativeBPMClient,responseAs=AsyncAwait,hideGenerationTimestamp=true
