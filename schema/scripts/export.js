import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { exportWorkflowOpenAPISchema, exportWorkflowJSONSchema } from '../dist/index.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const outDir = path.resolve(__dirname, '../../api/schemas');
fs.mkdirSync(outDir, { recursive: true });

const openapi = exportWorkflowOpenAPISchema();
fs.writeFileSync(path.join(outDir, 'workflow-ast.openapi.json'), JSON.stringify(openapi, null, 2) + '\n');

const jsonSchema = exportWorkflowJSONSchema();
fs.writeFileSync(path.join(outDir, 'workflow-ast.schema.json'), JSON.stringify(jsonSchema, null, 2) + '\n');

console.log('✓ Successfully exported Zod schemas to api/schemas/:');
console.log('  - api/schemas/workflow-ast.openapi.json');
console.log('  - api/schemas/workflow-ast.schema.json');
