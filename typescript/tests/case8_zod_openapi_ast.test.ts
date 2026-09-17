import test from 'node:test';
import assert from 'node:assert';
import {
  z,
  WorkflowBuilder,
  WorkflowASTSchema,
  exportWorkflowOpenAPISchema,
  exportWorkflowJSONSchema,
  fromJSONSchema,
  ZodError,
} from '../dist/index.js';

test('Кейс 8: Полная интеграция Zod 4 + NativeBPM SDK = OpenAPI 3.0 = JSON Schema = BPMN 2.0', async (t) => {
  // 1. Описание схем данных на Zod 4
  const ProcessVariablesSchema = z.object({
    customer_id: z.string().min(1, 'Customer ID is required'),
    tier: z.enum(['standard', 'premium', 'vip']),
  });

  const OrderFormSchema = z.object({
    order_id: z.string().min(1, 'Order ID is required'),
    items_count: z.number().int().positive('Must be positive integer'),
    notes: z.string().optional(),
  });

  const AIClassificationSchema = z.object({
    sentiment: z.enum(['positive', 'neutral', 'negative']),
    confidence: z.number().min(0).max(1),
    summary: z.string(),
  });

  // 2. Описание рабочего процесса через WorkflowBuilder с новыми возможностями
  const workflow = new WorkflowBuilder('order_orchestration', 'Order Orchestration Process')
    .variables(ProcessVariablesSchema)
    .start('start')
    .userTask('fill_order', 'Fill Order Details', {
      form: OrderFormSchema,
      formId: 'order_form_v1',
      candidateGroups: 'sales',
    })
    .boundaryTimer('sla_timer', 'SLA Timeout 15 Minutes', 'PT15M')
    .sequenceFlow('sla_timer', 'end_escalated')
    .end('end_escalated', 'SLA Escalation End')
    .sequenceFlow('fill_order', 'classify_notes')
    .aiTask('classify_notes', 'Classify Customer Notes', {
      prompt: 'Classify sentiment of order notes: ${notes}',
      responseSchema: AIClassificationSchema,
      resultVar: 'classification',
    })
    .serviceTask('finalize_order', 'Finalize Order in DB', 'db_writer')
    .end('end', 'Process Completed');

  // 3. Получение AST и валидация через Zod 4 WorkflowASTSchema
  const ast = workflow.toAST();
  assert.strictEqual(ast.id, 'order_orchestration');
  assert.strictEqual(ast.name, 'Order Orchestration Process');
  assert.ok(ast.inputSchema, 'AST должен содержать входную схему переменных процесса inputSchema');

  // Валидируем AST с помощью схемы Zod 4
  const parsedAST = WorkflowASTSchema.parse(ast);
  assert.strictEqual(parsedAST.id, 'order_orchestration');
  assert.strictEqual(parsedAST.nodes.length, 7, 'Должно быть 7 узлов (start, userTask, timer, aiTask, serviceTask, end, end_escalated)');

  // 4. Проверка экспорта в OpenAPI 3.0
  const openapiSchema = exportWorkflowOpenAPISchema();
  assert.strictEqual(openapiSchema.type, 'object');
  assert.ok(openapiSchema.properties.id, 'OpenAPI схема должна содержать свойство id');
  assert.ok(openapiSchema.properties.name, 'OpenAPI схема должна содержать свойство name');
  assert.ok(openapiSchema.properties.nodes, 'OpenAPI схема должна содержать свойство nodes');
  assert.ok(openapiSchema.properties.flows, 'OpenAPI схема должна содержать свойство flows');
  assert.ok(openapiSchema.properties.inputSchema, 'OpenAPI схема должна содержать свойство inputSchema');

  // 5. Проверка экспорта в JSON Schema Draft 2020-12
  const jsonSchema = exportWorkflowJSONSchema();
  assert.strictEqual(jsonSchema.$schema, 'https://json-schema.org/draft/2020-12/schema');
  assert.strictEqual(jsonSchema.type, 'object');
  assert.ok(jsonSchema.properties.id);
  assert.ok(jsonSchema.properties.nodes);

  // 6. Проверка генерации BPMN 2.0 XML
  const bpmnXml = workflow.toBPMN();
  assert.ok(bpmnXml.includes('<bpmn:definitions'), 'Должен содержать bpmn:definitions');
  assert.ok(bpmnXml.includes('id="order_orchestration"'), 'Должен содержать ID процесса');
  assert.ok(bpmnXml.includes('<bpmn:userTask id="fill_order"'), 'Должен содержать userTask');
  assert.ok(bpmnXml.includes('camunda:formKey="order_form_v1"'), 'Должен содержать formKey');
  assert.ok(bpmnXml.includes('<bpmn:boundaryEvent id="sla_timer"'), 'Должен содержать boundaryTimer');
  assert.ok(bpmnXml.includes('PT15M'), 'Должен содержать длительность таймера ISO 8601');
  assert.ok(bpmnXml.includes('nativebpm:responseSchema='), 'AI Task должен сериализовать responseSchema в XML');

  // 7. Проверка двусторонней регидратации формы: JSON Schema -> Zod 4 через fromJSONSchema
  const userNode = ast.nodes.find((n: any) => n.id === 'fill_order');
  assert.ok(userNode && userNode.inputSchema, 'Узел userTask обязан иметь inputSchema');
  const rawFormSchema = typeof userNode.inputSchema === 'string'
    ? JSON.parse(userNode.inputSchema)
    : userNode.inputSchema;

  const restoredFormValidator = fromJSONSchema(rawFormSchema);
  assert.ok(restoredFormValidator, 'fromJSONSchema должна успешно восстановить Zod-схему');

  // Валидные данные формы проходят
  const validFormData = { order_id: 'ORD-999', items_count: 5, notes: 'Deliver fast' };
  const parsedForm = restoredFormValidator.parse(validFormData);
  assert.strictEqual(parsedForm.order_id, 'ORD-999');
  assert.strictEqual(parsedForm.items_count, 5);

  // Невалидные данные отклоняются с типизированной ошибкой ZodError
  assert.throws(
    () => {
      restoredFormValidator.parse({ order_id: '', items_count: -3 });
    },
    (err: any) => {
      assert.ok(err instanceof ZodError, 'Ошибка должна быть ZodError');
      return true;
    }
  );

  // 8. Проверка Fail-Fast валидации таймера с невалидной ISO 8601 строкой
  assert.throws(
    () => {
      new WorkflowBuilder('invalid_timer_wf', 'Invalid Timer')
        .start('start')
        .userTask('task1', 'Task')
        .boundaryTimer('bad_timer', 'Bad Timer', 'not-an-iso-duration')
        .end('end');
    },
    (err: any) => {
      assert.ok(err.message.includes('ISO') || err.message.includes('duration'), 'Должна быть ошибка валидации ISO duration');
      return true;
    }
  );
});
