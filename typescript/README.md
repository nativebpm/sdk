[← Back to Platform Root](../../README.md)

# NativeBPM TypeScript SDK

Official TypeScript Client SDK for the NativeBPM Cloud-Native BPMN 2.0 / DMN 1.3 execution engine.

Built for **Node.js 22+ (native TypeScript execution without compilation via `--experimental-strip-types`)**, Edge CDN runtimes (Cloudflare Workers / Pages), and modern browsers.

---

## Key Features

- **Native Zod 4 Out-of-the-Box**: Re-exports `z`, `ZodError`, and `ZodType`. Define forms and process graphs in a single file (`*.flow.ts`).
- **Automatic JSON Schema Draft 2020-12 Compilation**: Automatically compiles Zod schemas into JSON Schema Draft 2020-12 for Server-Driven UI (`<nativebpm-trigger>`, Vue, React, Capacitor).
- **Three-Tier Architecture**:
  1. *Level 1 (SSOT)*: OpenAPI 3.0 specification (`openapi.yaml`).
  2. *Level 2 (Generated Transport)*: Strictly typed low-level fetch client.
  3. *Level 3 (Fluent API Façade & Worker)*: High-level developer experience for deploying, starting, claiming, and running workers.
- **Clean Code (Zero Parameter Properties)**: 100% compliant with Node 22+ native TypeScript type stripping (`node script.ts`).
- **Standard OMG BPMN 2.0 XML Export**: Compile directly to BPMN 2.0 XML via `workflow.toBPMN()`.
- **DMN 1.3 Decision Tables**: Built-in support for BusinessRuleTask and local rule evaluation.

---

## Installation

```bash
npm install @nativebpm/sdk
```

---

## 1. Single-File Workflow-as-Code & Schema-as-Code with Zod

```typescript
import { z, WorkflowBuilder } from "@nativebpm/sdk";

// 1. Declare form schema using native Zod 4
export const OrderFormSchema = z.object({
  order_id: z.string().min(1, "Order ID is required"),
  amount: z.number().positive("Amount must be positive"),
  promo_code: z.string().optional(),
});

// 2. Build the BPMN process with fluent API
export const orderWorkflow = new WorkflowBuilder("order_flow", "Order Process")
  .start("start")
  .userTask("fill_order", "Enter Order", {
    form: OrderFormSchema,
    formId: "order_form_v1",
    candidateGroups: "sales",
  })
  .exclusiveGateway("check_amount", "Check Order Amount")
    .when("amount > 100000").then("vip_review").userTask("vip_review", "VIP Review", { candidateGroups: "vip_managers" })
    .else("auto_approve").serviceTask("auto_approve", "Auto Approval", "auto_approver")
  .end("end", "Process Finished");

// 3. Extract compiled JSON Schema Draft 2020-12 for Server-Driven UI (BDUI)
const forms = orderWorkflow.extractForms();
console.log(forms.order_form_v1);
// Output: { $schema: "https://json-schema.org/draft/2020-12/schema", type: "object", ... }

// 4. Export standard OMG BPMN 2.0 XML
const bpmnXml = orderWorkflow.toBPMN();
```

---

## 2. Fluent Client API

The NativeBPM client provides a modular, fluent chain interface:

```typescript
import { Client } from "@nativebpm/sdk";
import { orderWorkflow } from "./order.flow.js";

const client = new Client("http://localhost:8080", "your-api-token");

// 1. Deploy workflow definition
const definition = await client.definitions().deploy()
  .withWorkflow(orderWorkflow)
  .send();

// 2. Start a new process instance
const instance = await client.instances()
  .start("order_flow")
  .withBusinessKey("ORD-2026-001")
  .withVariables({ amount: 1500, promo_code: "DISCOUNT10" })
  .send();

// 3. Query and claim human tasks
const tasks = await client.tasks().list()
  .withAssignee("sales_manager")
  .withStatus("CREATED")
  .send();

if (tasks.length > 0) {
  const task = tasks[0];
  await client.tasks().claim(task.id).withAssignee("sales_manager").send();

  // 4. Complete task with validated output variables
  await client.tasks().complete(task.id)
    .withVariables({ approved: true, comment: "Looks good" })
    .send();
}
```

---

## 3. Background Service Task Workers

NativeBPM features a dedicated `Worker` class with mirroring Go syntax for executing external service tasks:

```typescript
import { Worker } from "@nativebpm/sdk";

const worker = new Worker("http://localhost:8080", "your-api-token")
  .withWorkerId("billing-worker-01")
  .withMaxConcurrency(20)
  .withTopic("payment_gateway", async (task) => {
    console.log(`Processing task ${task.id} for instance ${task.instanceId}`);
    console.log("Input variables:", task.variables);

    if (task.variables.amount > 500000) {
      // Throwing an Error automatically transitions the task/instance into an Incident
      throw new Error("Credit limit exceeded");
    }

    // Return output variables to advance the process token
    return {
      payment_status: "PAID",
      tx_id: `tx_${Date.now()}`,
    };
  });

// Start background poll loop
await worker.start();
```

---

## Running Tests

The SDK includes a comprehensive 6-pattern TDD test suite covering:
1. Linear process with Zod validation (UserTask $\rightarrow$ ServiceTask).
2. DMN Decision Tables (BusinessRuleTask).
3. Conditional Branching (ExclusiveGateway When/Then/Else).
4. Service Task Workers & Incident transitions.
5. Boundary Timer Events & SLA escalation (`PT15M`).
6. Server-Driven UI (BDUI) JSON Schema Draft 2020-12 export.

```bash
# Build and run all tests
npm test
```
