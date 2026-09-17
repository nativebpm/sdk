import test from 'node:test';
import assert from 'node:assert';
import { z, WorkflowBuilder, evaluateDMNRule } from '../dist/index.js';

test('Кейс 2: Таблицы решений DMN (BusinessRuleTask)', async (t) => {
  // 1. Zod-схема для ввода данных о скидке
  const DiscountRequestSchema = z.object({
    customer_id: z.string().min(1),
    promo_code: z.string().optional(),
  });

  // 2. Описание процесса с BusinessRuleTask для DMN-таблицы
  const dmnWorkflow = new WorkflowBuilder('discount_evaluation_flow', 'Discount Evaluation')
    .start('start')
    .userTask('enter_promo', 'Ввод промокода клиентом', {
      form: DiscountRequestSchema,
      formId: 'promo_input_form',
    })
    .businessRuleTask('apply_discount_table', 'Таблица скидок DMN', 'discount_table', {
      hitPolicy: 'UNIQUE',
      inputs: [
        { expression: 'promo_code', type: 'string' },
      ],
      outputs: [
        { name: 'discount_percent', type: 'number' },
      ],
      rules: [
        { inputs: ['"DISCOUNT10"'], outputs: ['10'] },
        { inputs: ['"DISCOUNT20"'], outputs: ['20'] },
        { inputs: ['"-"'], outputs: ['0'] },
      ],
      resultVar: 'discount_percent',
      mapDecisionResult: 'singleEntry',
    })
    .serviceTask('apply_discount_service', 'Применение скидки в биллинге', 'billing_discount')
    .end('end', 'Завершено');

  // 3. Проверка AST и метаданных узла DMN
  const ast = dmnWorkflow.toAST();
  const ruleNode = ast.nodes.find(n => n.id === 'apply_discount_table');
  assert.ok(ruleNode, 'AST содержит узел apply_discount_table');
  assert.strictEqual(ruleNode.type, 'businessRuleTask');
  assert.strictEqual(ruleNode.decisionRef, 'discount_table');
  assert.strictEqual(ruleNode.hitPolicy, 'UNIQUE');
  assert.strictEqual(ruleNode.resultVar, 'discount_percent');
  assert.strictEqual(ruleNode.inputs?.length, 1);
  assert.strictEqual(ruleNode.outputs?.length, 1);
  assert.strictEqual(ruleNode.rules?.length, 3);

  // 4. Проверка результирующего BPMN XML
  const bpmnXml = dmnWorkflow.toBPMN();
  assert.ok(bpmnXml.includes('<bpmn:businessRuleTask id="apply_discount_table"'), 'XML содержит businessRuleTask');
  assert.ok(bpmnXml.includes('camunda:decisionRef="discount_table"'), 'XML содержит decisionRef');
  assert.ok(bpmnXml.includes('camunda:resultVariable="discount_percent"'), 'XML содержит resultVariable');

  // 5. Тестирование передачи переменных формы в DMN и расчет скидок:
  // Тест 5.1: Промокод DISCOUNT10 -> расчет скидки 10%
  const input1 = DiscountRequestSchema.parse({
    customer_id: 'CUST-001',
    promo_code: 'DISCOUNT10',
  });
  const result1 = evaluateDMNRule(ruleNode, input1);
  assert.strictEqual(result1.discount_percent, 10, 'Для промокода DISCOUNT10 скидка обязана составлять 10%');

  // Тест 5.2: Промокод DISCOUNT20 -> расчет скидки 20%
  const input2 = DiscountRequestSchema.parse({
    customer_id: 'CUST-002',
    promo_code: 'DISCOUNT20',
  });
  const result2 = evaluateDMNRule(ruleNode, input2);
  assert.strictEqual(result2.discount_percent, 20, 'Для промокода DISCOUNT20 скидка обязана составлять 20%');

  // Тест 5.3: Несуществующий промокод -> fallback rule "-" (0%)
  const input3 = DiscountRequestSchema.parse({
    customer_id: 'CUST-003',
    promo_code: 'UNKNOWN_PROMO',
  });
  const result3 = evaluateDMNRule(ruleNode, input3);
  assert.strictEqual(result3.discount_percent, 0, 'Для неизвестного промокода скидка должна быть 0%');
});
