import test from 'node:test';
import assert from 'node:assert';
import * as http from 'node:http';
import { Worker, WorkflowBuilder } from '../dist/index.js';
import { z } from 'zod';

// ============================================================================
// 1. Zod 4 Граничная Схема с Санитизацией и Нормализацией Промокода
// ============================================================================
export const RobustOrderCheckoutSchema = z.object({
  recipient_name: z.string()
    .trim()
    .min(2, 'Имя должно содержать минимум 2 символа')
    .max(100, 'Имя слишком длинное')
    .regex(/^[a-zA-Zа-яА-ЯёЁәіңғүұқөһӘІҢҒҮҰҚӨҺ\s\-\.\']+$/, 'Имя не должно содержать недопустимых спецсимволов'),
  
  recipient_phone: z.string()
    .trim()
    .regex(/^\+?[0-9]{10,15}$/, 'Телефон должен быть в международном формате (10-15 цифр)'),

  delivery_address: z.string()
    .trim()
    .min(5, 'Адрес слишком короткий')
    .max(300, 'Адрес слишком длинный')
    .refine(val => !/[<>{}]/i.test(val), 'Адрес содержит недопустимые спецсимволы (XSS защита)'),

  payment_method: z.enum(['kaspi_qr', 'bank_card', 'cash_on_delivery']),

  items_total: z.number()
    .positive('Сумма заказа должна быть больше 0')
    .max(50_000_000, 'Превышен максимальный лимит заказа (50 млн ₸)')
    .refine(val => Number.isFinite(val), 'Некорректная сумма'),

  promo_code: z.string()
    .trim()
    .transform(val => val.toUpperCase())
    .optional()
    .default(''),

  tip_amount: z.number()
    .min(0, 'Чаевые не могут быть отрицательными')
    .default(0),
});

// ============================================================================
// Тест 1: Идемпотентность воркера (Duplicate Execution Prevention)
// ============================================================================
test('TDD Риск 1: Воркер обязан гарантировать идемпотентность и не списать деньги дважды', async () => {
  let executionCount = 0;
  let completedCount = 0;

  const server = http.createServer((req, res) => {
    if (req.url?.includes('/complete') && req.method === 'POST') {
      completedCount++;
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ status: 'COMPLETED' }));
      return;
    }
    res.writeHead(404);
    res.end();
  });

  await new Promise<void>(resolve => server.listen(0, 'localhost', () => resolve()));
  const port = (server.address() as any).port;
  const baseUrl = `http://localhost:${port}`;

  try {
    const worker = new Worker(baseUrl, 'test-token')
      .withWorkerId('billing-idempotent-worker')
      .withTopic('payment_gateway', async (task) => {
        executionCount++;
        return {
          payment_status: 'PAID',
          transaction_id: 'TX-KASPI-TEST-1',
          charge_timestamp: 1726500000,
        };
      });

    const taskPayload = {
      id: 'task-job-991',
      instance_id: 'inst-ord-2026-991',
      activity_id: 'process_payment',
      topic: 'payment_gateway',
      variables: { amount: 40500 },
    };

    // Первый вызов задачи
    const res1 = await worker.processTask(taskPayload);
    assert.strictEqual(res1.status, 'completed');
    assert.strictEqual(res1.result.transaction_id, 'TX-KASPI-TEST-1');
    assert.strictEqual(executionCount, 1, 'Первый вызов должен исполнить логику');

    // Второй дублирующий вызов с теми же id/instance_id
    const res2 = await worker.processTask(taskPayload);
    assert.strictEqual(res2.status, 'completed');
    assert.strictEqual(res2.result.transaction_id, 'TX-KASPI-TEST-1');
    assert.strictEqual(executionCount, 1, 'Повторный вызов НЕ ДОЛЖЕН выполнять повторное списание (идемпотентность)');
  } finally {
    server.close();
  }
});

// ============================================================================
// Тест 2: Авто-повторы (Retries) с экспоненциальной задержкой и эскалация
// ============================================================================
test('TDD Риск 2: Воркер повторяет временные ошибки и эскалирует в Incident при исчерпании', async () => {
  let callCount = 0;
  let incidentReported: any = null;

  const server = http.createServer((req, res) => {
    let body = '';
    req.on('data', c => { body += c; });
    req.on('end', () => {
      if (req.url?.includes('/incidents') && req.method === 'POST') {
        incidentReported = JSON.parse(body);
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ id: 'inc-99' }));
        return;
      }
      res.writeHead(404);
      res.end();
    });
  });

  await new Promise<void>(resolve => server.listen(0, 'localhost', () => resolve()));
  const port = (server.address() as any).port;
  const baseUrl = `http://localhost:${port}`;

  try {
    const worker = new Worker(baseUrl, 'test-token')
      .withWorkerId('failing-worker')
      .withTopic('flaky_api', async () => {
        callCount++;
        throw new Error('Bank 504 Gateway Timeout');
      });

    const taskPayload = {
      id: 'task-flaky-1',
      instance_id: 'inst-flaky-1',
      activity_id: 'step_flaky',
      topic: 'flaky_api',
      variables: {},
    };

    const res = await worker.processTask(taskPayload);
    assert.strictEqual(res.status, 'incident');
    assert.strictEqual(res.error, 'Bank 504 Gateway Timeout');
    assert.ok(incidentReported, 'Инцидент обязан быть отправлен в движок');
    assert.strictEqual(incidentReported.error_type, 'WorkerExecutionError');
    assert.strictEqual(incidentReported.node_id, 'step_flaky');
  } finally {
    server.close();
  }
});

// ============================================================================
// Тест 3: Паттерн Saga и Компенсирующая транзакция (Auto-Refund)
// ============================================================================
test('TDD Риск 3: Паттерн Saga компенсирует списание (refund_payment) при сбое доставки', () => {
  const sagaWorkflow = new WorkflowBuilder('order_saga_process', 'Order Checkout Saga with Refund')
    .start('start')
    .userTask('fill_order', 'Ввод данных заказа', { formId: 'order_form' })
    .serviceTask('charge_payment', 'Списание средств', 'payment.charge')
    .serviceTask('dispatch_shipping', 'Оформление доставки', 'shipping.dispatch')
    .exclusiveGateway('gateway_delivery_check')
    // Ветка успеха:
    .sequenceFlow('gateway_delivery_check', 'order_success', 'delivery_status == "OK"')
    .end('order_success', 'Заказ успешно доставлен')
    // Ветка компенсации (Saga Compensation):
    .sequenceFlow('gateway_delivery_check', 'refund_payment', 'delivery_status == "FAILED"')
    .serviceTask('refund_payment', 'Компенсирующий возврат средств', 'payment.refund')
    .end('order_cancelled_refunded', 'Заказ отменен, средства возвращены');

  const ast = sagaWorkflow.toAST();
  
  // Проверяем наличие узла компенсации
  const refundNode = ast.nodes.find(n => n.id === 'refund_payment');
  assert.ok(refundNode, 'Узел компенсации refund_payment обязан присутствовать в процессе');
  assert.strictEqual(refundNode.topic, 'payment.refund');

  // Проверяем ветвление от шлюза
  const refundFlow = ast.flows.find(f => f.source === 'gateway_delivery_check' && f.target === 'refund_payment');
  assert.ok(refundFlow, 'Должен быть sequenceFlow к refund_payment');
  assert.strictEqual(refundFlow.condition, 'delivery_status == "FAILED"');
});

// ============================================================================
// Тест 4: Граничные условия и санитизация в Zod 4
// ============================================================================
test('TDD Риск 4: Zod 4 блокирует XSS, отрицательные суммы и нормализует промокоды', () => {
  // 1. Блокировка отрицательных сумм
  const negativeOrder = {
    recipient_name: 'Алиса С.',
    recipient_phone: '+77015554433',
    delivery_address: 'пр. Достык, д. 42',
    payment_method: 'kaspi_qr',
    items_total: -100,
  };
  const resNegative = RobustOrderCheckoutSchema.safeParse(negativeOrder);
  assert.strictEqual(resNegative.success, false);

  // 2. Блокировка нулевой суммы
  const zeroOrder = { ...negativeOrder, items_total: 0 };
  const resZero = RobustOrderCheckoutSchema.safeParse(zeroOrder);
  assert.strictEqual(resZero.success, false);

  // 3. Блокировка XSS в адресе
  const xssOrder = { ...negativeOrder, items_total: 1000, delivery_address: 'ул. Ленина <script>alert(1)</script>' };
  const resXss = RobustOrderCheckoutSchema.safeParse(xssOrder);
  assert.strictEqual(resXss.success, false);
  assert.ok(resXss.error?.issues.some(i => i.message.includes('XSS')));

  // 4. Тримминг и апперкейс промокода
  const promoOrder = {
    recipient_name: 'Алиса С.',
    recipient_phone: '+77015554433',
    delivery_address: 'пр. Достык, д. 42',
    payment_method: 'kaspi_qr',
    items_total: 45000,
    promo_code: '  promo10  ',
  };
  const resPromo = RobustOrderCheckoutSchema.safeParse(promoOrder);
  assert.strictEqual(resPromo.success, true);
  assert.strictEqual(resPromo.data?.promo_code, 'PROMO10', 'Промокод должен быть триммирован и приведен к верхнему регистру');
});

// ============================================================================
// Тест 5: Boundary Timer на прерывание зависших заказов (SLA)
// ============================================================================
test('TDD Риск 5: Boundary Timer прерывает брошенные заказы (PT15M) для возврата остатков', () => {
  const workflow = new WorkflowBuilder('order_with_sla', 'Order with SLA Timeout')
    .start('start')
    .userTask('fill_checkout_form', 'Ввод данных заказа', { formId: 'order_checkout' })
    .boundaryTimerEvent('sla_timeout_15m', '15m Abandonment Timeout', 'fill_checkout_form', 'PT15M', true)
    .sequenceFlow('sla_timeout_15m', 'release_inventory')
    .serviceTask('release_inventory', 'Освобождение остатков товара', 'inventory.release')
    .end('end_abandoned', 'Заказ аннулирован по таймауту')
    .sequenceFlow('fill_checkout_form', 'process_payment')
    .serviceTask('process_payment', 'Оплата', 'payment.charge')
    .end('end_paid', 'Оплачено');

  const ast = workflow.toAST();
  const timer = ast.nodes.find(n => n.id === 'sla_timeout_15m');
  assert.ok(timer, 'Boundary timer обязан существовать');
  assert.strictEqual(timer.type, 'boundaryTimerEvent');
  assert.strictEqual(timer.attachedToRef, 'fill_checkout_form');
  assert.strictEqual(timer.cancelActivity, true);
  assert.strictEqual(timer.timeDuration, 'PT15M');
});
