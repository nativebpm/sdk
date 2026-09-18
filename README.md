# NativeBPM TypeScript SDK & Zod Schema

> Pure TypeScript Monorepo for NativeBPM BPMN 2.0 / DMN 1.3 Execution Engine.

This branch (`typescript`) is dedicated to TypeScript and JavaScript developers. It contains zero non-TypeScript language directories, providing an isolated, pure TypeScript developer experience with root npm workspaces.

---

## 📦 Packages

| Package | Directory | Description |
|---|---|---|
| **`@nativebpm/schema`** | [`schema/`](./schema) | Pure Zod 4 Schemas & AST Type Definitions. Zero network dependencies (~15 KB). |
| **`@nativebpm/sdk`** | [`typescript/`](./typescript) | Full Client SDK, Fluent Workflow & DMN Builder, BDUI Form Extractor, and Service Task Workers. |

---

## 🚀 Quick Start

### 1. Installation

```bash
# Full SDK (includes Fluent Builder, Client, and Schemas)
npm install @nativebpm/sdk

# Or pure schemas only (lightweight, zero network overhead)
npm install @nativebpm/schema
```

### 2. Building a BPMN Process with Zod 4 Schema Validation

```typescript
import { Client, Workflow, z } from "@nativebpm/sdk";

// 1. Declare process input variables contract using Zod 4
const OrderSchema = z.object({
  orderId: z.string().min(1),
  amount: z.number().positive(),
  tier: z.enum(["standard", "vip"]),
});

// 2. Build the workflow graph with embedded JSON Schema and branching
const workflow = new Workflow("order-fulfillment", "Order Fulfillment Process")
  .variables(OrderSchema) // automatically exported to JSON Schema Draft 2020-12
  .start("start", "Start")
  .exclusiveGateway("check_tier", "Check Tier")
  .when("tier == 'vip'").then("vip_task").serviceTask("vip_task", "VIP Processing", "vip_topic")
  .otherwise("std_task").serviceTask("std_task", "Standard Processing", "std_topic")
  .end("end", "End");

// 3. Connect to NativeBPM Engine
const client = new Client("http://localhost:8080", "your-auth-token");

// 4. Deploy process definition
await client.deploy(workflow);

// 5. Fail-Fast validation on the client before network transmission
const rawInput = { orderId: "ORD-991", amount: 150.0, tier: "vip" };
const validated = OrderSchema.parse(rawInput);

// 6. Start process instance
const instance = await client.instances().start("order-fulfillment")
  .withVariables(validated)
  .send();

console.log("Started process instance:", instance.id);
```

---

## 🧪 Testing

Run all test suites across all workspaces directly from the repository root:

```bash
# Run tests for both @nativebpm/schema and @nativebpm/sdk
npm test

# Build all packages
npm run build
```

---

## 🏛️ Architecture & Philosophy

- **Zero Foreign Language Footprint**: Only TypeScript, Zod, and OpenAPI specs live in this branch.
- **Root NPM Workspaces**: Seamless multi-package development without `npm link` friction.
- **Fail-Fast Principles**: Strict schema validation at the perimeter before hitting the network.
