.PHONY: generate test

OPENAPI_GEN_IMG = registry.gitlab.com/nativebpm/sdk/openapi-generator-cli:latest
GOLANG_IMG = registry.gitlab.com/nativebpm/sdk/golang:1.26-alpine

test:
	cd go && GOWORK=off CGO_ENABLED=0 go test -count=1 -v ./...

generate:
	docker run --rm -v "$$(pwd):/local" $(OPENAPI_GEN_IMG) generate \
		-i /local/api/openapi.yaml \
		-g go \
		-o /local/go \
		--additional-properties=packageName=nativebpm,hideGenerationTimestamp=true \
		--git-host gitlab.com \
		--git-user-id nativebpm \
		--git-repo-id sdk/go
	mkdir -p go/api
	docker run --rm -v "$$(pwd):/local" -w /local/go $(GOLANG_IMG) sh -c "go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.7.1 -package api -generate types,std-http,spec,strict-server /local/api/openapi.yaml > /local/go/api/api.gen.go"
	docker run --rm -v "$$(pwd):/local" -w /local/go $(GOLANG_IMG) go mod tidy
