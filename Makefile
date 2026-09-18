.PHONY: generate test
OPENAPI_GEN_IMG = registry.gitlab.com/nativebpm/sdk/openapi-generator-cli:latest
RUST_IMG = registry.gitlab.com/nativebpm/sdk/rust:1.96

test:
	cd rust && cargo test

generate:
	docker run --rm -v "$$(pwd):/local" $(OPENAPI_GEN_IMG) generate \
		-i /local/api/openapi.yaml \
		-g rust \
		-o /local/rust \
		--additional-properties=packageName=nativebpm-client,hideGenerationTimestamp=true
