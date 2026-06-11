import { validateJSON, parseSchema } from './dist/index.js';
import * as assert from 'node:assert';
import * as path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const wasmPath = path.resolve(__dirname, 'dist/jsonschema.wasm');

async function testValidation() {
  console.log("Running JSON Schema validation tests via WASM...");

  const schema = {
    type: "object",
    required: ["username", "email", "age"],
    properties: {
      username: { type: "string", minLength: 3, maxLength: 10 },
      email: { type: "string", pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$" },
      age: { type: "integer", minimum: 18, maximum: 99 },
      isAdmin: { type: "boolean" }
    }
  };

  // Case 1: Valid payload
  {
    const payload = { username: "alice", email: "alice@example.com", age: 30, isAdmin: true };
    const res = await validateJSON(schema, payload, wasmPath);
    console.log("Case 1 (Valid Payload):", res);
    assert.strictEqual(res.valid, true, "Payload should be valid");
    assert.strictEqual(res.errors.length, 0, "Should have 0 errors");
  }

  // Case 2: Missing required field
  {
    const payload = { username: "alice", email: "alice@example.com" };
    const res = await validateJSON(schema, payload, wasmPath);
    console.log("Case 2 (Missing Age):", res);
    assert.strictEqual(res.valid, false, "Payload should be invalid due to missing field");
    assert.ok(res.errors.some(e => e.includes("Field 'age' is required")), "Error message must report missing age");
  }

  // Case 3: Value constraints (minLength & maximum & pattern)
  {
    const payload = { username: "al", email: "invalid-email", age: 100, isAdmin: "yes" };
    const res = await validateJSON(schema, payload, wasmPath);
    console.log("Case 3 (Constraints violation):", res);
    assert.strictEqual(res.valid, false, "Payload should be invalid");
    assert.ok(res.errors.some(e => e.includes("username") && e.includes("at least 3")), "Should fail minLength");
    assert.ok(res.errors.some(e => e.includes("email") && e.includes("pattern")), "Should fail email regex pattern");
    assert.ok(res.errors.some(e => e.includes("age") && e.includes("less than or equal to maximum")), "Should fail age maximum");
    assert.ok(res.errors.some(e => e.includes("isAdmin") && e.includes("must be a boolean")), "Should fail type match");
  }

  console.log("Validation tests completed successfully!");
}

async function testParsing() {
  console.log("Running JSON Schema parsing/compilation tests via WASM...");

  const schema = {
    type: "object",
    required: ["name"],
    "ui:order": ["age", "name", "bio"],
    properties: {
      name: { type: "string", title: "Full Name", default: "John Doe" },
      age: { type: "integer", minimum: 0, "x-step": 1 },
      bio: { type: "string", maxLength: 200 },
      status: { type: "string", enum: ["active", "inactive"] }
    }
  };

  const variables = { age: 25 };

  const widgets = await parseSchema(schema, variables, wasmPath);
  console.log("Parsed Widgets:", JSON.stringify(widgets, null, 2));

  assert.strictEqual(widgets.length, 4, "Should have compiled 4 widgets");

  // Verify ui:order and sorting (ordered keys first, then remaining sorted alphabetically)
  assert.strictEqual(widgets[0].name, "age", "First should be age");
  assert.strictEqual(widgets[1].name, "name", "Second should be name");
  assert.strictEqual(widgets[2].name, "bio", "Third should be bio");
  assert.strictEqual(widgets[3].name, "status", "Fourth should be status (remaining sorted)");

  // Verify constraints and mapping
  const ageWidget = widgets[0];
  assert.strictEqual(ageWidget.type, "integer", "age type is integer");
  assert.strictEqual(ageWidget.widget, "number", "age widget is number");
  assert.strictEqual(ageWidget.value, 25, "age injected value is 25");
  assert.strictEqual(ageWidget.xStep, 1, "age step is 1");

  const nameWidget = widgets[1];
  assert.strictEqual(nameWidget.label, "Full Name", "name title mapped to label");
  assert.strictEqual(nameWidget.required, true, "name is required");
  assert.strictEqual(nameWidget.value, "John Doe", "name value fallback to default");

  const bioWidget = widgets[2];
  assert.strictEqual(bioWidget.widget, "textarea", "long strings fallback to textarea");

  const statusWidget = widgets[3];
  assert.strictEqual(statusWidget.widget, "select", "enum string maps to select");
  assert.deepStrictEqual(statusWidget.options, ["active", "inactive"], "options mapped");

  console.log("Schema parsing tests completed successfully!");
}

async function main() {
  try {
    await testValidation();
    await testParsing();
    console.log("All TS/JS WebAssembly JSON Schema tests passed successfully!");
  } catch (error) {
    console.error("Test execution failed:", error);
    process.exit(1);
  }
}

main();
