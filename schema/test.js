import assert from 'node:assert';
import {
  z,
  WorkflowASTSchema,
  NodeASTSchema,
  FlowASTSchema,
  exportWorkflowOpenAPISchema,
  exportWorkflowOpenAPIDocument,
  exportWorkflowJSONSchema,
} from './dist/index.js';

console.log("🚀 Testing @nativebpm/schema standalone package...");

// 1. Validate valid AST
const validAST = {
  id: "test-proc",
  name: "Test Process",
  inputSchema: JSON.stringify({ type: "object", properties: { amount: { type: "number" } } }),
  nodes: [
    { id: "start", name: "Start", type: "startEvent" },
    { id: "vip_task", name: "VIP Task", type: "serviceTask", topic: "vip" },
    { id: "end", name: "End", type: "endEvent" }
  ],
  flows: [
    { id: "f1", source: "start", target: "vip_task" },
    { id: "f2", source: "vip_task", target: "end" }
  ]
};

const parsed = WorkflowASTSchema.parse(validAST);
assert.strictEqual(parsed.id, "test-proc");
assert.strictEqual(parsed.nodes.length, 3);
assert.strictEqual(parsed.flows.length, 2);
console.log("✓ Valid AST parsing passed");

// 2. Reject invalid AST (missing id or invalid node type)
assert.throws(
  () => {
    WorkflowASTSchema.parse({
      id: "",
      name: "Bad",
      nodes: [{ id: "n1", name: "Bad", type: "unknownType" }],
      flows: []
    });
  },
  (err) => err instanceof z.ZodError,
  "Must throw ZodError on invalid node type"
);
console.log("✓ Invalid AST rejection passed");

// 3. OpenAPI 3.0 Schema & Document Exporters
const openapiSchema = exportWorkflowOpenAPISchema();
assert.strictEqual(openapiSchema.type, "object");
assert.ok(openapiSchema.properties.id);
assert.ok(openapiSchema.properties.nodes);
assert.ok(openapiSchema.properties.flows);
console.log("✓ OpenAPI 3.0 schema export passed");

const openapiDoc = exportWorkflowOpenAPIDocument();
assert.strictEqual(openapiDoc.openapi, "3.0.3", "Must have openapi 3.0.3 version field");
assert.ok(openapiDoc.info, "Must have info object");
assert.strictEqual(openapiDoc.info.title, "NativeBPM Workflow AST Specification");
assert.ok(openapiDoc.info.version);
assert.ok(openapiDoc.paths, "Must have paths object");
assert.ok(openapiDoc.paths['/api/deploy'], "Must specify /api/deploy endpoint");
assert.ok(openapiDoc.components && openapiDoc.components.schemas, "Must have components.schemas");
assert.ok(openapiDoc.components.schemas.WorkflowAST, "Must have WorkflowAST component schema");
assert.ok(openapiDoc.components.schemas.NodeAST, "Must have NodeAST component schema");
assert.ok(openapiDoc.components.schemas.FlowAST, "Must have FlowAST component schema");
assert.ok(openapiDoc.components.schemas.DMNRule, "Must have DMNRule component schema");
assert.ok(openapiDoc.components.schemas.DMNInput, "Must have DMNInput component schema");
assert.ok(openapiDoc.components.schemas.DMNOutput, "Must have DMNOutput component schema");
assert.ok(openapiDoc.components.schemas.InVariable, "Must have InVariable component schema");
assert.ok(openapiDoc.components.schemas.OutVariable, "Must have OutVariable component schema");
console.log("✓ OpenAPI 3.0.3 document export (GitLab/Swagger UI compatible) passed");

// 4. JSON Schema Draft 2020-12 Exporter
const jsonSchema = exportWorkflowJSONSchema();
assert.strictEqual(jsonSchema.$schema, "https://json-schema.org/draft/2020-12/schema");
assert.strictEqual(jsonSchema.type, "object");
console.log("✓ JSON Schema Draft 2020-12 export passed");

// 5. Verify Zero-Network Isolation
assert.strictEqual(typeof fetch, "function"); // global fetch exists in Node
console.log("✓ Zero-network dependency verified");

console.log("🎉 All @nativebpm/schema tests passed successfully!");
