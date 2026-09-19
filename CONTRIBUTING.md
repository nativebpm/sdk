# Contributing to NativeBPM Contracts & Native Code Recipes

Thank you for contributing to NativeBPM! This repository maintains the authoritative OpenAPI 3.0 specification (`api/openapi.yaml`), JSON Schema definitions (`schema/`), and Zero-SDK recipes for Native Code integration.

---

## ⚡ Architecture: Pure Native Code (Zero-SDK)

NativeBPM has unified all client integration around **a single standard: Native Code (Workflow-as-Code)**.

Developers interact directly with NativeBPM using:
1. **Native AST Parsers**: Go (`go/parser`), Python (async AST), TypeScript/JavaScript (ECMA-262).
2. **Declarative State Machines**: JSON / YAML / DMN 1.3 definitions.
3. **Wasmee WASM Runtime**: Sandboxed WebAssembly workers via Unix Domain Sockets (< 50µs SLA).
4. **Universal Standard HTTP/UDS**: Direct REST calls using language standard libraries (`net/http`, `urllib`, `fetch`, `HttpClient`, `curl`).

There are **no generated client SDK packages** to build, publish, or install.

---

## 🛠️ API & Schema Contribution Workflow

### 1. Modifying the OpenAPI 3.0 Contract
* Edit [`api/openapi.yaml`](api/openapi.yaml).
* Follow OpenAPI 3.0.3 standards with explicit request/response schemas.
* Verify that all endpoints conform to Pure Engine Core boundaries (BPMN 2.0 / DMN 1.3 / OpenAPI standard).

### 2. Modifying JSON Schemas
* Edit schemas in [`schema/`](schema/).
* Ensure schemas comply with JSON Schema Draft 2020-12.

### 3. Adding Zero-SDK Recipes
* If adding new language examples, submit them to the [Zero-SDK Cookbook](https://gitlab.com/nativebpm/docs/-/blob/main/nativebpm/NB-284/zero_sdk_cookbook_ru.md).
* All examples must use **only built-in standard libraries** without external third-party packages.
