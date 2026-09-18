.PHONY: generate generate-go generate-python generate-typescript generate-java generate-php generate-dotnet generate-rust generate-kotlin generate-swift generate-dart
.PHONY: test test-schema test-go test-python test-typescript test-java test-kotlin test-php test-dotnet test-rust test-dart test-swift
.PHONY: push-images login-registry push-gradle push-composer push-php push-dotnet push-rust push-dart push-maven push-node push-docker-git push-docker-dind push-openapi-gen push-golang push-python push-node-alpine

# Registry-backed images to bypass Docker Hub rate limits in CI/CD
OPENAPI_GEN_IMG = registry.gitlab.com/nativebpm/sdk/openapi-generator-cli:latest
GOLANG_IMG = registry.gitlab.com/nativebpm/sdk/golang:1.26-alpine
PYTHON_IMG = registry.gitlab.com/nativebpm/sdk/python:3.12
NODE_IMG = registry.gitlab.com/nativebpm/sdk/node:22-alpine
GRADLE_IMG = registry.gitlab.com/nativebpm/sdk/gradle:8-jdk21
COMPOSER_IMG = registry.gitlab.com/nativebpm/sdk/composer:latest
PHP_IMG = registry.gitlab.com/nativebpm/sdk/php:8.3-cli
DOTNET_IMG = registry.gitlab.com/nativebpm/sdk/dotnet-sdk:9.0
RUST_IMG = registry.gitlab.com/nativebpm/sdk/rust:1.96
DART_IMG = registry.gitlab.com/nativebpm/sdk/dart:stable
MAVEN_IMG = registry.gitlab.com/nativebpm/sdk/maven:3.9-eclipse-temurin-21
NODE22_IMG = registry.gitlab.com/nativebpm/sdk/node:22
DOCKER_GIT_IMG = registry.gitlab.com/nativebpm/sdk/docker:27
DOCKER_DIND_IMG = registry.gitlab.com/nativebpm/sdk/docker:27-dind

# BuildKit-powered code generator for polyglot SDKs
BUILDX = docker buildx build -f Dockerfile.gen

generate:
	$(BUILDX) --target export-all --output type=local,dest=$${DEST:-./out} .

generate-go:
	$(BUILDX) --target export-go --output type=local,dest=$${DEST:-./out/go} .

generate-python:
	$(BUILDX) --target export-python --output type=local,dest=$${DEST:-./out/python} .

generate-typescript:
	$(BUILDX) --target export-typescript --output type=local,dest=$${DEST:-./out/typescript} .

generate-java:
	$(BUILDX) --target export-java --output type=local,dest=$${DEST:-./out/java} .

generate-php:
	$(BUILDX) --target export-php --output type=local,dest=$${DEST:-./out/php} .

generate-dotnet:
	$(BUILDX) --target export-dotnet --output type=local,dest=$${DEST:-./out/dotnet} .

generate-rust:
	$(BUILDX) --target export-rust --output type=local,dest=$${DEST:-./out/rust} .

generate-kotlin:
	$(BUILDX) --target export-kotlin --output type=local,dest=$${DEST:-./out/kotlin} .

generate-swift:
	$(BUILDX) --target export-swift --output type=local,dest=$${DEST:-./out/swift} .

generate-dart:
	$(BUILDX) --target export-dart --output type=local,dest=$${DEST:-./out/dart} .

test: test-schema test-go test-python test-typescript test-java test-kotlin test-php test-dotnet test-rust test-dart

test-schema:
	docker run --rm -v "$$(pwd):/local" -w /local/schema $(NODE_IMG) sh -c "npm install --no-package-lock && npm test"


test-go:
	docker run --rm -v "$$(pwd):/local" -w /local/go $(GOLANG_IMG) sh -c "go get github.com/stretchr/testify/assert && go test -v ./..."

test-python:
	docker run --rm -v "$$(pwd):/local" -w /local/python $(PYTHON_IMG) sh -c "pip install -r requirements.txt -r test-requirements.txt && python -m unittest discover -s test"

test-typescript:
	docker run --rm -v "$$(pwd):/local" -w /local/typescript $(NODE_IMG) sh -c "apk add --no-cache make && make test"

test-java:
	docker run --rm -v "$$(pwd):/local" -w /local/java $(GRADLE_IMG) gradle test

test-kotlin:
	docker run --rm -v "$$(pwd):/local" -w /local/kotlin $(GRADLE_IMG) gradle test

test-php:
	docker run --rm -v "$$(pwd):/local" -w /local/php $(COMPOSER_IMG) composer install --no-interaction
	docker run --rm -v "$$(pwd):/local" -w /local/php $(PHP_IMG) vendor/bin/phpunit

test-dotnet:
	docker run --rm -v "$$(pwd):/local" -w /local/dotnet $(DOTNET_IMG) dotnet test NativeBPM.Client.sln

test-rust:
	docker run --rm -v "$$(pwd):/local" -w /local/rust $(RUST_IMG) cargo test

test-dart:
	docker run --rm -v "$$(pwd):/local" -w /local/dart $(DART_IMG) dart test

test-swift:
	@echo "No tests configured for Swift"

push-images: login-registry push-gradle push-composer push-php push-dotnet push-rust push-dart push-maven push-node push-docker-git push-docker-dind push-openapi-gen push-golang push-python push-node-alpine

login-registry:
	@if [ -z "$$GL_PAT" ]; then echo "Error: GL_PAT environment variable is not set." && exit 1; fi
	docker login -u sauran -p $$GL_PAT registry.gitlab.com

define mirror_image
	@docker manifest inspect $(2) >/dev/null 2>&1 || ( \
		docker pull $(1) && \
		docker tag $(1) $(2) && \
		docker push $(2) \
	)
endef

push-gradle:
	$(call mirror_image,mirror.gcr.io/library/gradle:8-jdk21,$(GRADLE_IMG))

push-composer:
	$(call mirror_image,mirror.gcr.io/library/composer:latest,$(COMPOSER_IMG))

push-php:
	$(call mirror_image,mirror.gcr.io/library/php:8.3-cli,$(PHP_IMG))

push-dotnet:
	$(call mirror_image,mcr.microsoft.com/dotnet/sdk:9.0,$(DOTNET_IMG))

push-rust:
	$(call mirror_image,mirror.gcr.io/library/rust:1.96,$(RUST_IMG))

push-dart:
	$(call mirror_image,mirror.gcr.io/library/dart:stable,$(DART_IMG))

push-maven:
	$(call mirror_image,mirror.gcr.io/library/maven:3.9-eclipse-temurin-21,$(MAVEN_IMG))

push-node:
	$(call mirror_image,mirror.gcr.io/library/node:22,$(NODE22_IMG))

push-docker-git:
	$(call mirror_image,mirror.gcr.io/library/docker:27,$(DOCKER_GIT_IMG))

push-docker-dind:
	$(call mirror_image,mirror.gcr.io/library/docker:27-dind,$(DOCKER_DIND_IMG))

push-openapi-gen:
	$(call mirror_image,openapitools/openapi-generator-cli:latest,$(OPENAPI_GEN_IMG))

push-golang:
	$(call mirror_image,mirror.gcr.io/library/golang:1.26-alpine,$(GOLANG_IMG))

push-python:
	$(call mirror_image,mirror.gcr.io/library/python:3.12,$(PYTHON_IMG))

push-node-alpine:
	$(call mirror_image,mirror.gcr.io/library/node:22-alpine,$(NODE_IMG))


