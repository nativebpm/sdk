.PHONY: generate generate-typescript build test test-schema test-typescript

OPENAPI_GEN_IMG = registry.gitlab.com/nativebpm/sdk/openapi-generator-cli:latest
NODE_IMG = registry.gitlab.com/nativebpm/sdk/node:22-alpine

build:
	npm run build

test: test-schema test-typescript

test-schema:
	cd schema && npm test

test-typescript:
	cd typescript && npm test

generate: generate-typescript

generate-typescript:
	docker run --rm -v "$$(pwd):/local" $(OPENAPI_GEN_IMG) generate \
		-i /local/api/openapi.yaml \
		-g typescript-fetch \
		-o /local/typescript/src/api \
		--additional-properties=npmName=@nativebpm/client,npmVersion=1.0.0,hideGenerationTimestamp=true
