.PHONY: generate test
OPENAPI_GEN_IMG = registry.gitlab.com/nativebpm/sdk/openapi-generator-cli:latest
DART_IMG = registry.gitlab.com/nativebpm/sdk/dart:stable

test:
	cd dart && dart test

generate:
	docker run --rm -v "$$(pwd):/local" $(OPENAPI_GEN_IMG) generate \
		-i /local/api/openapi.yaml \
		-g dart \
		-o /local/dart \
		--additional-properties=pubName=nativebpm_client,pubVersion=1.0.0,pubDescription="NativeBPM Client SDK for Dart and Flutter",hideGenerationTimestamp=true
