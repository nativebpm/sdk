import { writeFileSync, mkdirSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';
import {
  exportWorkflowOpenAPISchema,
  exportWorkflowJSONSchema,
} from '../dist/index.js';

const __dirname = dirname(fileURLToPath(import.meta.url));
const targetSchemasDir = join(__dirname, '../../api/schemas');

mkdirSync(targetSchemasDir, { recursive: true });

// 1. Generate OpenAPI 3.0 schema
const openapiSchema = exportWorkflowOpenAPISchema();
const openapiPath = join(targetSchemasDir, 'workflow-ast.openapi.json');
writeFileSync(openapiPath, JSON.stringify(openapiSchema, null, 2), 'utf-8');
console.log(`✅ Generated OpenAPI 3.0 schema: ${openapiPath}`);

// 2. Generate JSON Schema Draft 2020-12
const jsonSchema = exportWorkflowJSONSchema();
const jsonSchemaPath = join(targetSchemasDir, 'workflow-ast.schema.json');
writeFileSync(jsonSchemaPath, JSON.stringify(jsonSchema, null, 2), 'utf-8');
console.log(`✅ Generated JSON Schema Draft 2020-12: ${jsonSchemaPath}`);
