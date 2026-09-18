[← Back to Platform Root](../../README.md)

# NativeBPM TypeScript SDK

Official TypeScript Client SDK for the NativeBPM Cloud-Native BPMN 2.0 / DMN 1.3 execution engine.

Built for **Node.js 22+ (native TypeScript execution without compilation via `--experimental-strip-types`)**, Edge CDN runtimes (Cloudflare Workers / Pages), and modern browsers.

---

## Key Features

- **Native Zod 4 Out-of-the-Box**: Re-exports `z`, `ZodError`, `ZodType`, and `fromJSONSchema`. Define forms, variable contracts, and process graphs in a single file (`*.flow.ts`).
- **Formal Workflow AST Schemas**: Includes 7 modular Zod schemas (`WorkflowASTSchema`, `NodeASTSchema`, `FlowASTSchema`, `DMNRuleSchema`, `DMNInputSchema`, `DMNOutputSchema`, `InVariableSchema`) serving as the Single Source of Truth for OpenAPI 3.0 and JSON Schema Draft 2020-12.
- **Process Variables Contract (`.variables(schema)`)**: Enforce typed process start variables directly at the API edge.
- **AI Task Structured Outputs (`aiTask({ responseSchema })`)**: Pass Zod schemas directly to AI tasks; NativeBPM serializes them to standard JSON Schema to guarantee structured LLM responses.
- **Fail-Fast Boundary Timers**: Validates ISO 8601 duration strings at build time using `z.iso.duration()`.
- **Shortened Branching Fluent API**: Ultra-concise `when(cond).then(target).otherwise(target)` syntax with automatic decision gateway generation, plus closure-block DSL.
- **Two-Way Server-Driven UI (BDUI)**: Compile Zod forms to JSON Schema for `<nativebpm-trigger>`, and re-hydrate live Zod validators on the frontend using `fromJSONSchema()`.
- **Standard OMG BPMN 2.0 XML Export**: Compile directly to standard BPMN 2.0 XML via `workflow.toBPMN()`.
- **Clean Code (Zero Parameter Properties)**: 100% compliant with Node 22+ native TypeScript type stripping (`node script.ts`).

---

## Installation

```bash
npm install @nativebpm/sdk
```

---

## 1. Single-File Workflow-as-Code & Schema-as-Code with Zod 4

```typescript
import { z, WorkflowBuilder, fromJSONSchema } from "@nativebpm/sdk";

// 1. Process start variables contract
export const ProcessInputSchema = z.object({
  customerId: z.string().min(1, "Customer ID is required"),
  orderAmount: z.number().positive("Amount must be positive"),
  tier: z.enum(["standard", "premium", "vip"]),
});

// 2. UserTask UI Form Schema (Server-Driven UI)
export const OrderFormSchema = z.object({
  order_id: z.string().min(1, "Order ID is required"),
  amount: z.number().positive("Amount must be positive"),
  promo_code: z.string().optional(),
});

// 3. AI Service Task Structured Output Schema
export const RiskAnalysisSchema = z.object({
  riskScore: z.number().min(0).max(100),
  recommendation: z.enum(["approve", "manual_review", "reject"]),
  reasons: z.array(z.string()),
});

// 4. Build the BPMN process with fluent API
export const orderWorkflow = new WorkflowBuilder("order_flow", "Order Process")
  .variables(ProcessInputSchema) // Enforce process input contract
  .start("start")
  .aiTask("analyze_risk", "Analyze Order Risk", {
    prompt: "Assess fraud risk for customer: ${customerId}, amount: ${orderAmount}",
    responseSchema: RiskAnalysisSchema,
    resultVar: "risk",
  })
  // Shortened Fluent API: automatically creates and connects exclusiveGateway
  .when("risk.recommendation === 'manual_review'")
    .then("fill_order")
    .userTask("fill_order", "Manual Order Review", {
      form: OrderFormSchema,
      formId: "order_form_v1",
      candidateGroups: "risk_team",
    })
    .boundaryTimer("sla_timer", "SLA 15 Minutes", "PT15M") // ISO 8601 validated
  .otherwise("auto_approve")
    .serviceTask("auto_approve", "Auto Approval", "auto_approver")
  .end("end", "Process Finished");

// 5. Extract compiled JSON Schema Draft 2020-12 for Server-Driven UI (BDUI)
const forms = orderWorkflow.extractForms();
console.log(forms.order_form_v1);

// 6. Two-Way Re-hydration on Frontend:
const clientValidator = fromJSONSchema(forms.order_form_v1);
const validation = clientValidator.safeParse({ order_id: "ORD-1", amount: 250 });

// 7. Export standard OMG BPMN 2.0 XML
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
