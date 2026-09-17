import test from 'node:test';
import assert from 'node:assert';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import {
  WorkflowASTSchema,
  exportWorkflowOpenAPISchema,
  exportWorkflowJSONSchema,
} from '../dist/index.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const sdkRoot = path.resolve(__dirname, '../..');

test('Кейс 9: 360-градусное тестирование генерации SDK на все 10 языков из единого Zod 4 / OpenAPI контракта', async (t) => {
  // -------------------------------------------------------------------------
  // 1. Zod 4: Первоисточник контракта и схема Workflow AST
  // -------------------------------------------------------------------------
  await t.test('1. Zod 4 Schema Authority & Exporters', () => {
    const openapiAST = exportWorkflowOpenAPISchema();
    assert.strictEqual(openapiAST.type, 'object', 'OpenAPI Schema must be object');
    assert.ok(openapiAST.properties.id, 'WorkflowAST must have id property');
    assert.ok(openapiAST.properties.name, 'WorkflowAST must have name property');
    assert.ok(openapiAST.properties.nodes, 'WorkflowAST must have nodes property');
    assert.ok(openapiAST.properties.flows, 'WorkflowAST must have flows property');
    assert.ok(openapiAST.properties.inputSchema, 'WorkflowAST must have inputSchema property');

    const jsonSchemaAST = exportWorkflowJSONSchema();
    assert.strictEqual(jsonSchemaAST.$schema, 'https://json-schema.org/draft/2020-12/schema');
    assert.strictEqual(jsonSchemaAST.type, 'object');

    // Валидация эталонного AST
    const testAST = {
      id: 'polyglot_ast_process',
      name: 'Polyglot 360 Process',
      inputSchema: JSON.stringify({
        type: 'object',
        properties: { orderId: { type: 'string' } },
        required: ['orderId'],
      }),
      nodes: [
        {
          id: 'start',
          name: 'Start',
          type: 'startEvent',
        },
        {
          id: 'process_order',
          name: 'Process Order Task',
          type: 'serviceTask',
          topic: 'order_processing',
        },
        {
          id: 'end',
          name: 'End',
          type: 'endEvent',
        },
      ],
      flows: [
        { id: 'f1', source: 'start', target: 'process_order' },
        { id: 'f2', source: 'process_order', target: 'end' },
      ],
    };

    const parsed = WorkflowASTSchema.parse(testAST);
    assert.strictEqual(parsed.id, 'polyglot_ast_process');
    assert.strictEqual(parsed.nodes.length, 3);
    assert.strictEqual(parsed.flows.length, 2);
  });

  // -------------------------------------------------------------------------
  // 2. OpenAPI 3.0 Спецификация (sdk/api/openapi.yaml)
  // -------------------------------------------------------------------------
  await t.test('2. OpenAPI Specification AST Components Sync', () => {
    const openapiYamlPath = path.join(sdkRoot, 'api/openapi.yaml');
    assert.ok(fs.existsSync(openapiYamlPath), 'api/openapi.yaml must exist');
    const content = fs.readFileSync(openapiYamlPath, 'utf8');

    // Проверяем наличие всех 6 AST компонентов в components/schemas
    assert.ok(content.includes('WorkflowAST:'), 'openapi.yaml must declare WorkflowAST schema');
    assert.ok(content.includes('NodeAST:'), 'openapi.yaml must declare NodeAST schema');
    assert.ok(content.includes('FlowAST:'), 'openapi.yaml must declare FlowAST schema');
    assert.ok(content.includes('InVariable:'), 'openapi.yaml must declare InVariable schema');
    assert.ok(content.includes('OutVariable:'), 'openapi.yaml must declare OutVariable schema');
    assert.ok(content.includes('DMNRuleAST:'), 'openapi.yaml must declare DMNRuleAST schema');

    // Проверяем использование в POST /api/deploy
    assert.ok(
      content.includes("$ref: '#/components/schemas/WorkflowAST'"),
      '/api/deploy must accept WorkflowAST'
    );
  });

  // -------------------------------------------------------------------------
  // 3. 360-градусная верификация генерации моделей для всех 10 SDK
  // -------------------------------------------------------------------------
  const polyglotTargets = [
    {
      language: 'Go',
      modelFile: 'go/api/api.gen.go',
      checks: [
        'type WorkflowAST struct {',
        'type NodeAST struct {',
        'type FlowAST struct {',
        'InputSchema *string',
      ],
    },
    {
      language: 'Python',
      modelFile: 'python/nativebpm_client/models/workflow_ast.py',
      checks: [
        'class WorkflowAST(BaseModel):',
        'id: Annotated[str, Field(',
        'nodes: List[NodeAST]',
        'flows: List[FlowAST]',
        'input_schema: Optional[StrictStr]',
      ],
    },
    {
      language: 'TypeScript',
      modelFile: 'typescript/src/api/src/models/WorkflowAST.ts',
      checks: [
        'export interface WorkflowAST {',
        'id: string;',
        'name: string;',
        'nodes: Array<NodeAST>;',
        'flows: Array<FlowAST>;',
        'inputSchema?: string;',
      ],
    },
    {
      language: 'Java',
      modelFile: 'java/src/main/java/com/nativebpm/client/model/WorkflowAST.java',
      checks: [
        'public class WorkflowAST {',
        'SERIALIZED_NAME_ID = "id"',
        'SERIALIZED_NAME_NODES = "nodes"',
        'SERIALIZED_NAME_FLOWS = "flows"',
        'SERIALIZED_NAME_INPUT_SCHEMA = "inputSchema"',
      ],
    },
    {
      language: 'C# (.NET)',
      modelFile: 'dotnet/src/NativeBPM.Client/Model/WorkflowAST.cs',
      checks: [
        'public partial class WorkflowAST',
        'public string Id { get; set; }',
        'public List<NodeAST> Nodes { get; set; }',
        'public List<FlowAST> Flows { get; set; }',
      ],
    },
    {
      language: 'Rust',
      modelFile: 'rust/src/models/workflow_ast.rs',
      checks: [
        'pub struct WorkflowAst {',
        'pub id: String,',
        'pub nodes: Vec<models::NodeAst>,',
        'pub flows: Vec<models::FlowAst>,',
      ],
    },
    {
      language: 'PHP',
      modelFile: 'php/lib/Model/WorkflowAST.php',
      checks: [
        'class WorkflowAST implements ModelInterface',
        "'id' => 'string'",
        "'nodes' => '\\NativeBPMClient\\Model\\NodeAST[]'",
        "'flows' => '\\NativeBPMClient\\Model\\FlowAST[]'",
      ],
    },
    {
      language: 'Kotlin',
      modelFile: 'kotlin/src/main/kotlin/com/nativebpm/client/models/WorkflowAST.kt',
      checks: [
        'data class WorkflowAST (',
        '@Json(name = "id")',
        'val nodes: kotlin.collections.List<NodeAST>',
        'val flows: kotlin.collections.List<FlowAST>',
      ],
    },
    {
      language: 'Swift 5',
      modelFile: 'swift/NativeBPMClient/Classes/OpenAPIs/Models/WorkflowAST.swift',
      checks: [
        'public struct WorkflowAST: Codable',
        'public var id: String',
        'public var nodes: [NodeAST]',
        'public var flows: [FlowAST]',
      ],
    },
    {
      language: 'Dart',
      modelFile: 'dart/lib/model/workflow_ast.dart',
      checks: [
        'class WorkflowAST {',
        'required this.id,',
        'required this.name,',
        'List<NodeAST> nodes;',
        'List<FlowAST> flows;',
      ],
    },
  ];

  for (const target of polyglotTargets) {
    await t.test(`3.${polyglotTargets.indexOf(target) + 1} Target SDK: ${target.language}`, () => {
      const fullPath = path.join(sdkRoot, target.modelFile);
      assert.ok(
        fs.existsSync(fullPath),
        `Generated model file must exist: ${target.modelFile} for ${target.language}`
      );
      const fileContent = fs.readFileSync(fullPath, 'utf8');
      for (const check of target.checks) {
        assert.ok(
          fileContent.includes(check),
          `Target ${target.language} (${target.modelFile}) must contain expected token: "${check}"`
        );
      }
    });
  }
});
