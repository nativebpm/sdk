import test from 'node:test';
import assert from 'node:assert';
import * as http from 'node:http';
import { Worker } from '../dist/index.js';

test('Кейс 4: Воркеры сервисных задач (Service Tasks & Workers)', async (t) => {
  let completedTaskId = '';
  let completedVariables: any = null;
  let incidentReported: any = null;

  // 1. Создаем локальный mock-сервер движка NativeBPM
  const server = http.createServer((req, res) => {
    let body = '';
    req.on('data', chunk => { body += chunk; });
    req.on('end', () => {
      const parsedBody = body ? JSON.parse(body) : {};

      // Завершение задачи
      if (req.url?.startsWith('/api/tasks/') && req.url.endsWith('/complete') && req.method === 'POST') {
        const parts = req.url.split('/');
        completedTaskId = parts[3];
        completedVariables = parsedBody.variables;

        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({
          id: 'inst-001',
          process_id: 'order_flow',
          completed: true,
          variables: completedVariables,
        }));
        return;
      }

      // Перевод в Incident при исключении
      if (req.url?.includes('/incidents') && req.method === 'POST') {
        incidentReported = parsedBody;
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ status: 'incident_created', id: 'inc-001' }));
        return;
      }

      res.writeHead(404);
      res.end();
    });
  });

  await new Promise<void>(resolve => {
    server.listen(0, 'localhost', () => resolve());
  });
  const port = (server.address() as any).port;
  const baseUrl = `http://localhost:${port}`;

  try {
    // 2. Инициализация воркера с точным Go-синтаксисом
    const worker = new Worker(baseUrl, 'test-token')
      .withWorkerId('billing-worker')
      .withMaxConcurrency(20)
      .withTopic('payment_gateway', async (task) => {
        if (task.variables.amount > 50000) {
          throw new Error('Credit limit exceeded');
        }
        return { payment_status: 'PAID', tx_id: 'tx_123' };
      });

    assert.strictEqual(worker.getWorkerId(), 'billing-worker');
    assert.deepStrictEqual(worker.getRegisteredTopics(), ['payment_gateway']);

    // 3. Тестирование успешного выполнения задачи воркером
    const taskSuccess = {
      id: 'task-success-1',
      topic: 'payment_gateway',
      instance_id: 'inst-001',
      activity_id: 'payment_step',
      variables: { amount: 1000, order_id: 'ORD-100' },
    };

    const resSuccess = await worker.processTask(taskSuccess);
    assert.strictEqual(resSuccess.status, 'completed', 'Статус успешного выполнения должен быть completed');
    assert.strictEqual(resSuccess.result?.payment_status, 'PAID');
    assert.strictEqual(resSuccess.result?.tx_id, 'tx_123');
    assert.strictEqual(completedTaskId, 'task-success-1');
    assert.strictEqual(completedVariables?.payment_status, 'PAID');
    assert.strictEqual(completedVariables?.tx_id, 'tx_123');

    // 4. Тестирование поведения при выбросе ошибки (throw Error) -> перевод в Incident
    const taskFail = {
      id: 'task-fail-2',
      topic: 'payment_gateway',
      instance_id: 'inst-002',
      activity_id: 'payment_step',
      variables: { amount: 99999, order_id: 'ORD-999' },
    };

    const resFail = await worker.processTask(taskFail);
    assert.strictEqual(resFail.status, 'incident', 'При ошибке статус должен перейти в incident');
    assert.strictEqual(resFail.error, 'Credit limit exceeded');
    assert.ok(incidentReported, 'Сервер должен получить запрос на создание инцидента');
    assert.strictEqual(incidentReported.error_type, 'WorkerExecutionError');
    assert.strictEqual(incidentReported.error_message, 'Credit limit exceeded');
    assert.strictEqual(incidentReported.worker_id, 'billing-worker');

  } finally {
    server.close();
  }
});
