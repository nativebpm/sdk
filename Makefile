.PHONY: generate test
OPENAPI_GEN_IMG = registry.gitlab.com/nativebpm/sdk/openapi-generator-cli:latest
GRADLE_IMG = registry.gitlab.com/nativebpm/sdk/gradle:8-jdk21

test:
	cd kotlin && gradle test

generate:
	docker run --rm -v "$$(pwd):/local" $(OPENAPI_GEN_IMG) generate \
		-i /local/api/openapi.yaml \
		-g kotlin \
		-o /local/kotlin \
		--additional-properties=groupId=com.nativebpm,artifactId=nativebpm-kotlin-client,artifactVersion=1.0.0,packageName=com.nativebpm.client,library=jvm-okhttp4,hideGenerationTimestamp=true
