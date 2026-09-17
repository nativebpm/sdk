import test from 'node:test';
import assert from 'node:assert';
import { WorkflowBuilder } from '../dist/index.js';

test('Кейс 3: Условное ветвление (ExclusiveGateway: When / Then / Else)', async (t) => {
  const builder = new WorkflowBuilder('order_approval_gateway', 'Order Approval Process');

  // Использование точного эргономичного синтаксиса из технического задания:
  builder
    .start('start')
    .exclusiveGateway('check_amount', 'Проверка суммы заказа')
    .when('amount > 100000').then('vip_review').userTask('vip_review', 'VIP Проверка', { candidateGroups: 'vip_managers' })
    .else('auto_approve').serviceTask('auto_approve', 'Автоодобрение', 'auto_approver')
    .end('end', 'Процесс завершен');

  // 1. Проверка структуры AST
  const ast = builder.toAST();
  assert.strictEqual(ast.id, 'order_approval_gateway');

  // Проверка наличия шлюза
  const gwNode = ast.nodes.find(n => n.id === 'check_amount');
  assert.ok(gwNode, 'Шлюз check_amount должен присутствовать в AST');
  assert.strictEqual(gwNode.type, 'exclusiveGateway');

  // Проверка целевых задач
  const vipNode = ast.nodes.find(n => n.id === 'vip_review');
  assert.ok(vipNode, 'Узел vip_review должен присутствовать');
  assert.strictEqual(vipNode.type, 'userTask');
  assert.strictEqual(vipNode.candidateGroups, 'vip_managers');

  const autoNode = ast.nodes.find(n => n.id === 'auto_approve');
  assert.ok(autoNode, 'Узел auto_approve должен присутствовать');
  assert.strictEqual(autoNode.type, 'serviceTask');
  assert.strictEqual(autoNode.topic, 'auto_approver');

  // 2. Проверка sequence flows и условий
  const condFlow = ast.flows.find(f => f.source === 'check_amount' && f.target === 'vip_review');
  assert.ok(condFlow, 'Должен существовать sequence flow от check_amount к vip_review');
  assert.strictEqual(condFlow.condition, 'amount > 100000', 'Условие должно быть строго amount > 100000');

  const elseFlow = ast.flows.find(f => f.source === 'check_amount' && f.target === 'auto_approve');
  assert.ok(elseFlow, 'Должен существовать sequence flow от check_amount к auto_approve');
  assert.strictEqual(elseFlow.condition, '', 'Ветка else не должна содержать условий');

  // 3. Проверка генерации правильного BPMN 2.0 XML с conditionExpression
  const xml = builder.toBPMN();

  assert.ok(xml.includes('<bpmn:exclusiveGateway id="check_amount"'), 'XML содержит exclusiveGateway');
  assert.ok(xml.includes('<bpmn:sequenceFlow id="flow-check_amount-vip_review" sourceRef="check_amount" targetRef="vip_review">'), 'XML содержит conditional sequence flow');
  assert.ok(xml.includes('<bpmn:conditionExpression xsi:type="bpmn:tFormalExpression">amount &gt; 100000</bpmn:conditionExpression>'), 'XML экранирует и содержит корректный conditionExpression');
  assert.ok(xml.includes('<bpmn:sequenceFlow id="flow-check_amount-auto_approve" sourceRef="check_amount" targetRef="auto_approve" />'), 'XML содержит else sequence flow без условий');
});
