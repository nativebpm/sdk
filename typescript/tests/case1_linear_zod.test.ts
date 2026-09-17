import test from 'node:test';
import assert from 'node:assert';
import { z, WorkflowBuilder, ZodError } from '../dist/index.js';

test('Кейс 1: Линейный процесс с валидацией Zod (UserTask -> ServiceTask)', async (t) => {
  // 1. Описание схемы формы через Zod
  const OrderFormSchema = z.object({
    order_id: z.string().min(1, 'ID заказа обязателен'),
    amount: z.number().positive('Сумма должна быть положительной'),
    promo_code: z.string().optional(),
  });

  // 2. Описание графа процесса через WorkflowBuilder
  const orderWorkflow = new WorkflowBuilder('order_flow', 'Order Process')
    .start('start')
    .userTask('fill_order', 'Ввод заказа', {
      form: OrderFormSchema,
      formId: 'order_form_v1',
      candidateGroups: 'sales',
    })
    .serviceTask('process_payment', 'Обработка платежа', 'payment_gateway')
    .end('end', 'Завершение процесса');

  // 3. Проверка компиляции Zod-схемы в Draft 2020-12 JSON Schema
  const forms = orderWorkflow.extractForms();
  assert.ok(forms.order_form_v1, 'Словарь форм обязан содержать order_form_v1');
  const schema = forms.order_form_v1;

  assert.strictEqual(schema.$schema, 'https://json-schema.org/draft/2020-12/schema', 'Схема должна соответствовать Draft 2020-12');
  assert.strictEqual(schema.type, 'object', 'Корневой тип схемы должен быть object');
  assert.ok(schema.properties.order_id, 'Схема должна содержать свойство order_id');
  assert.strictEqual(schema.properties.order_id.type, 'string');
  assert.strictEqual(schema.properties.order_id.minLength, 1);
  assert.ok(schema.properties.amount, 'Схема должна содержать свойство amount');
  assert.strictEqual(schema.properties.amount.type, 'number');
  assert.strictEqual(schema.properties.amount.exclusiveMinimum, 0);
  assert.ok(schema.required.includes('order_id'), 'order_id должен быть в списке required');
  assert.ok(schema.required.includes('amount'), 'amount должен быть в списке required');

  // 4. Проверка корректности AST
  const ast = orderWorkflow.toAST();
  assert.strictEqual(ast.id, 'order_flow');
  assert.strictEqual(ast.name, 'Order Process');
  assert.strictEqual(ast.nodes.length, 4, 'Должно быть 4 узла (start, userTask, serviceTask, end)');
  assert.strictEqual(ast.flows.length, 3, 'Должно быть 3 перехода (sequenceFlows)');

  const userNode = ast.nodes.find(n => n.id === 'fill_order');
  assert.ok(userNode, 'AST содержит узел fill_order');
  assert.strictEqual(userNode.type, 'userTask');
  assert.strictEqual(userNode.candidateGroups, 'sales');
  assert.strictEqual(userNode.formId, 'order_form_v1');

  const serviceNode = ast.nodes.find(n => n.id === 'process_payment');
  assert.ok(serviceNode, 'AST содержит узел process_payment');
  assert.strictEqual(serviceNode.type, 'serviceTask');
  assert.strictEqual(serviceNode.topic, 'payment_gateway');

  // 5. Проверка результирующего BPMN 2.0 XML
  const bpmnXml = orderWorkflow.toBPMN();
  assert.ok(bpmnXml.includes('<bpmn:definitions'), 'XML должен содержать bpmn:definitions');
  assert.ok(bpmnXml.includes('xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL"'), 'Стандартный неймспейс OMG BPMN 2.0');
  assert.ok(bpmnXml.includes('<bpmn:userTask id="fill_order" name="Ввод заказа"'), 'XML содержит bpmn:userTask');
  assert.ok(bpmnXml.includes('camunda:candidateGroups="sales"'), 'XML содержит candidateGroups');
  assert.ok(bpmnXml.includes('camunda:formKey="order_form_v1"'), 'XML содержит formKey');
  assert.ok(bpmnXml.includes('<bpmn:serviceTask id="process_payment" name="Обработка платежа"'), 'XML содержит bpmn:serviceTask');
  assert.ok(bpmnXml.includes('nativebpm:topic="payment_gateway"'), 'XML содержит nativebpm:topic');

  // 6. Валидация входных данных: корректные проходят
  const validData = {
    order_id: 'ORD-2026-001',
    amount: 1500.5,
    promo_code: 'WELCOME10',
  };
  const parsedData = OrderFormSchema.parse(validData);
  assert.deepStrictEqual(parsedData, validData, 'Корректные данные должны успешно пройти валидацию');

  // 7. Невалидные отклоняются с типизированными ошибками ZodError
  assert.throws(
    () => {
      OrderFormSchema.parse({
        order_id: '',
        amount: -100,
      });
    },
    (err: any) => {
      assert.ok(err instanceof ZodError, 'Ошибка должна быть экземпляром ZodError');
      const issues = err.issues;
      assert.ok(issues.some((i: any) => i.path.includes('order_id')), 'Должна быть ошибка для order_id');
      assert.ok(issues.some((i: any) => i.path.includes('amount')), 'Должна быть ошибка для amount');
      return true;
    },
    'Невалидные данные обязаны выбрасывать типизированный ZodError'
  );
});
