# @nativebpm/schema

> Pure Zod 4 Schemas & AST Type Definitions for NativeBPM Workflows, BPMN 2.0, and DMN 1.3.

Zero runtime dependencies (only `zod ^4.6.5`). Zero network calls. Hermetic, ultra-lightweight (~15 KB), and portable across frontend, backend, edge workers, and CLI tools.

---

## 📦 Installation

```bash
npm install @nativebpm/schema
```

*(Or use the unified `@nativebpm/sdk` bundle which re-exports this schema package via `@nativebpm/sdk/schemas`).*

---

## 🚀 Quick Start

### 1. Validating a Workflow AST

```typescript
import { WorkflowASTSchema } from '@nativebpm/schema';

const workflow = {
  id: "order-fulfillment",
  name: "Order Fulfillment Process",
  nodes: [
    { id: "start", name: "Order Placed", type: "startEvent" },
    { id: "charge", name: "Charge Card", type: "serviceTask", topic: "payment" },
    { id: "end", name: "Order Complete", type: "endEvent" }
  ],
  flows: [
    { id: "f1", source: "start", target: "charge" },
    { id: "f2", source: "charge", target: "end" }
  ]
};

// Parse and validate strictly
const validatedAST = WorkflowASTSchema.parse(workflow);
```

### 2. Exporting to OpenAPI 3.0 or JSON Schema (Draft 2020-12)

```typescript
import { exportWorkflowOpenAPISchema, exportWorkflowJSONSchema } from '@nativebpm/schema';

// Export as OpenAPI 3.0 Schema Object
const openApiSchema = exportWorkflowOpenAPISchema();

// Export as JSON Schema Draft 2020-12 document
const jsonSchema = exportWorkflowJSONSchema();
```

---

## 🏛️ Architecture & Philosophy

1. **Single Source of Truth (SSoT)**: Zod schemas define the runtime validation rules and TypeScript types simultaneously.
2. **Deterministic Generation**: OpenAPI 3.0 and JSON Schema specifications are derived directly from the Zod schemas via `toJSONSchema()`.
3. **Zero Network Dependency**: This package does not perform HTTP requests or rely on external network access. It can safely run in edge workers, sandboxes, and offline environments.
