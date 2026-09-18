.PHONY: generate test

OPENAPI_GEN_IMG = registry.gitlab.com/nativebpm/sdk/openapi-generator-cli:latest
PYTHON_IMG = registry.gitlab.com/nativebpm/sdk/python:3.12

test:
	cd python && python3 -m unittest discover -s test && python3 test.py && python3 test_client.py

generate:
	docker run --rm -v "$$(pwd):/local" $(OPENAPI_GEN_IMG) generate \
		-i /local/api/openapi.yaml \
		-g python \
		-o /local/python \
		--additional-properties=packageName=nativebpm_client,hideGenerationTimestamp=true
