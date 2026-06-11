.PHONY: build test

build:
	npm install
	npm run build
	cp ../python/nativebpm/core.wasm ./dist/
	cp ../python/nativebpm/jsonschema.wasm ./dist/

test: build
	node test.js
