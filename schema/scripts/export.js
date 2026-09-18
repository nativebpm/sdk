import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import {
  exportWorkflowOpenAPIDocument,
  exportWorkflowOpenAPISchema,
  exportWorkflowJSONSchema,
} from '../dist/index.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const outDir = path.resolve(__dirname, '../../api/schemas');
fs.mkdirSync(outDir, { recursive: true });

// 1. Standalone OpenAPI 3.0.3 specification document (GitLab / Swagger UI preview compatible)
const openapiDoc = exportWorkflowOpenAPIDocument();
fs.writeFileSync(path.join(outDir, 'workflow-ast.openapi.json'), JSON.stringify(openapiDoc, null, 2) + '\n');

// 2. Standard JSON Schema Draft 2020-12
const jsonSchema = exportWorkflowJSONSchema();
fs.writeFileSync(path.join(outDir, 'workflow-ast.schema.json'), JSON.stringify(jsonSchema, null, 2) + '\n');

console.log('✓ Successfully exported Zod schemas to api/schemas/:');
console.log('  - api/schemas/workflow-ast.openapi.json (OpenAPI 3.0.3 document)');
console.log('  - api/schemas/workflow-ast.schema.json (JSON Schema Draft 2020-12)');

