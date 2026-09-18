.PHONY: generate test
OPENAPI_GEN_IMG = registry.gitlab.com/nativebpm/sdk/openapi-generator-cli:latest
PHP_IMG = registry.gitlab.com/nativebpm/sdk/php:8.3-cli

test:
	cd php && vendor/bin/phpunit

generate:
	docker run --rm -v "$$(pwd):/local" $(OPENAPI_GEN_IMG) generate \
		-i /local/api/openapi.yaml \
		-g php \
		-o /local/php \
		--additional-properties=invokerPackage=NativeBPM\\Client,packageName=nativebpm/client,hideGenerationTimestamp=true
