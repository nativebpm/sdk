import test from 'node:test';
import assert from 'node:assert';
import { WorkflowBuilder } from '../dist/index.js';

test('Кейс 5: Таймеры и SLA-эскалация (Boundary Timer Events)', async (t) => {
  // 1. Построение процесса с ожиданием формы и таймером SLA-эскалации (PT15M)
  const slaWorkflow = new WorkflowBuilder('application_sla_process', 'Application Processing with SLA')
    .start('start')
    .userTask('fill_application', 'Заполнение заявки клиентом', {
      formId: 'application_form_v1',
      candidateGroups: 'clients',
    })
    // Навешивание граничного таймера ожидания на задачу fill_application: 15 минут (PT15M)
    .boundaryTimerEvent('timer_sla_15m', '15 Minute SLA Timeout', 'fill_application', 'PT15M', true)
    // Ветка эскалации при срабатывании таймера
    .sequenceFlow('timer_sla_15m', 'escalate_to_boss')
    .serviceTask('escalate_to_boss', 'Эскалация начальнику отдела', 'manager_escalation')
    .end('end_escalated', 'Завершено по эскалации')
    // Основная ветка нормального исполнения
    .sequenceFlow('fill_application', 'process_application')
    .serviceTask('process_application', 'Стандартная обработка заявки', 'app_processing')
    .end('end_normal', 'Успешно завершено');

  // 2. Проверка AST
  const ast = slaWorkflow.toAST();
  assert.strictEqual(ast.id, 'application_sla_process');

  const timerNode = ast.nodes.find(n => n.id === 'timer_sla_15m');
  assert.ok(timerNode, 'Узел таймера timer_sla_15m обязан присутствовать в AST');
  assert.strictEqual(timerNode.type, 'boundaryTimerEvent');
  assert.strictEqual(timerNode.attachedToRef, 'fill_application', 'Таймер должен быть привязан к fill_application');
  assert.strictEqual(timerNode.timeDuration, 'PT15M', 'Длительность таймера должна быть PT15M');
  assert.strictEqual(timerNode.cancelActivity, true, 'Таймер должен прерывать активность (cancelActivity=true)');

  const escalationFlow = ast.flows.find(f => f.source === 'timer_sla_15m' && f.target === 'escalate_to_boss');
  assert.ok(escalationFlow, 'Должен существовать sequence flow от таймера к задаче эскалации');

  const normalFlow = ast.flows.find(f => f.source === 'fill_application' && f.target === 'process_application');
  assert.ok(normalFlow, 'Должен существовать нормальный sequence flow от формы к обработке');

  // 3. Проверка результирующего BPMN 2.0 XML
  const xml = slaWorkflow.toBPMN();

  assert.ok(
    xml.includes('<bpmn:boundaryEvent id="timer_sla_15m" name="15 Minute SLA Timeout" attachedToRef="fill_application" cancelActivity="true">'),
    'BPMN XML обязан содержать тег bpmn:boundaryEvent с привязкой attachedToRef и cancelActivity'
  );
  assert.ok(
    xml.includes('<bpmn:timerEventDefinition id="timerDef_timer_sla_15m">'),
    'BPMN XML обязан содержать timerEventDefinition'
  );
  assert.ok(
    xml.includes('<bpmn:timeDuration xsi:type="bpmn:tFormalExpression">PT15M</bpmn:timeDuration>'),
    'BPMN XML обязан содержать спецификацию интервала PT15M в timeDuration'
  );
  assert.ok(
    xml.includes('<bpmn:sequenceFlow id="flow-timer_sla_15m-escalate_to_boss" sourceRef="timer_sla_15m" targetRef="escalate_to_boss" />'),
    'BPMN XML обязан содержать sequenceFlow от таймера к эскалации'
  );
});
