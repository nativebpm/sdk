import test from 'node:test';
import assert from 'node:assert/strict';
import http from 'node:http';
import { Workflow, computeWorkflowHash, clearDeployedHashCache, Client, setDefaultClient } from '../dist/index.js';

test('Кейс 10: Just-In-Time (JIT) Auto-Deploy & Single-Call Process Execution', async (t) => {
  await t.test('1. Deterministic Content Hash (SHA-256)', () => {
    const wf1 = new Workflow('order_checkout', 'Order Checkout Process')
      .start()
      .user('fill_checkout_form', 'Fill Checkout Form')
      .service('charge_card', 'Charge Card', 'payment_topic')
      .end();

    const hash1 = wf1.getContentHash();
    assert.match(hash1, /^sha256:[a-f0-9]{64}$/, 'Content hash must be valid sha256: format');

    // Same definition produces identical hash
    const wf2 = new Workflow('order_checkout', 'Order Checkout Process')
      .start()
      .user('fill_checkout_form', 'Fill Checkout Form')
      .service('charge_card', 'Charge Card', 'payment_topic')
      .end();
    assert.equal(wf1.getContentHash(), wf2.getContentHash(), 'Identical ASTs must produce identical hashes');

    // Modified definition produces different hash
    const wf3 = new Workflow('order_checkout', 'Order Checkout Process')
      .start()
      .user('fill_checkout_form', 'Fill Checkout Form')
      .service('charge_card_v2', 'Charge Card with 3DS', 'payment_topic_3ds')
      .end();
    assert.notEqual(wf1.getContentHash(), wf3.getContentHash(), 'Modified AST must produce different hash');
  });

  await t.test('2. Single-Call Execution: await workflow.run(payload) with Cold and Hot Paths', async () => {
    clearDeployedHashCache();

    let requestCount = 0;
    const receivedRequests: any[] = [];

    // Spin up mock server for testing .run()
    const server = http.createServer((req, res) => {
      if (req.method === 'POST' && req.url === '/api/process/execute') {
        let body = '';
        req.on('data', (chunk) => (body += chunk));
        req.on('end', () => {
          requestCount++;
          const parsed = JSON.parse(body);
          receivedRequests.push(parsed);

          const isCold = !!parsed.ast;
          res.writeHead(200, { 'Content-Type': 'application/json' });
          res.end(
            JSON.stringify({
              instanceId: `inst_${requestCount}`,
              definitionId: parsed.definitionId || parsed.ast?.id,
              version: 1,
              isNewVersionDeployed: isCold,
              status: 'ACTIVE',
              currentTasks: [
                {
                  id: 'task_1',
                  activity_id: 'fill_checkout_form',
                  name: 'Fill Checkout Form',
                  assignee: '',
                  candidate_groups: '',
                  status: 'CREATED',
                  created_at: new Date().toISOString(),
                },
              ],
            })
          );
        });
        return;
      }
      res.writeHead(404);
      res.end();
    });

    await new Promise<void>((resolve) => server.listen(0, resolve));
    const port = (server.address() as any).port;
    const client = new Client(`http://127.0.0.1:${port}`, 'test-token');

    try {
      const orderCheckoutWorkflow = new Workflow('order_checkout', 'Order Checkout')
        .start()
        .user('fill_checkout_form', 'Fill Checkout Form')
        .end();

      // 2.1 Cold path: First call sends full AST
      const instance1 = await orderCheckoutWorkflow.run(
        {
          recipient_phone: '+77015551234',
          items_total: 45000,
        },
        { client, businessKey: 'BK-1001' }
      );

      assert.equal(instance1.instanceId, 'inst_1');
      assert.equal(instance1.isNewVersionDeployed, true);
      assert.equal(instance1.version, 1);
      assert.equal(instance1.status, 'ACTIVE');
      assert.equal(instance1.currentTasks.length, 1);
      assert.equal(instance1.currentTasks[0].activity_id, 'fill_checkout_form');

      assert.equal(receivedRequests.length, 1);
      assert.ok(receivedRequests[0].ast, 'Cold request must contain full AST');
      assert.equal(receivedRequests[0].businessKey, 'BK-1001');

      // 2.2 Hot path: Second call with same workflow uses cached contentHash (0ms AST payload)
      const instance2 = await orderCheckoutWorkflow.run(
        {
          recipient_phone: '+77015555678',
          items_total: 62000,
        },
        { client, businessKey: 'BK-1002' }
      );

      assert.equal(instance2.instanceId, 'inst_2');
      assert.equal(instance2.isNewVersionDeployed, false);
      assert.equal(receivedRequests.length, 2);
      assert.equal(receivedRequests[1].ast, undefined, 'Hot request must NOT re-transmit AST body');
      assert.equal(receivedRequests[1].definitionId, 'order_checkout');
      assert.equal(receivedRequests[1].contentHash, orderCheckoutWorkflow.getContentHash());
    } finally {
      server.close();
    }
  });

  await t.test('3. Transparent Recovery on Cache Eviction / Unknown Hash', async () => {
    clearDeployedHashCache();

    let attempts = 0;
    const requests: any[] = [];

    const server = http.createServer((req, res) => {
      if (req.method === 'POST' && req.url === '/api/process/execute') {
        let body = '';
        req.on('data', (chunk) => (body += chunk));
        req.on('end', () => {
          attempts++;
          const parsed = JSON.parse(body);
          requests.push(parsed);

          // If hot request without AST, simulate server cache restart with 404
          if (!parsed.ast) {
            res.writeHead(404, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'DEFINITION_HASH_UNKNOWN' }));
            return;
          }

          res.writeHead(200, { 'Content-Type': 'application/json' });
          res.end(
            JSON.stringify({
              instanceId: 'inst_recovered',
              definitionId: parsed.ast.id,
              version: 2,
              isNewVersionDeployed: true,
              status: 'ACTIVE',
              currentTasks: [],
            })
          );
        });
        return;
      }
      res.writeHead(404);
      res.end();
    });

    await new Promise<void>((resolve) => server.listen(0, resolve));
    const port = (server.address() as any).port;
    const client = new Client(`http://127.0.0.1:${port}`, 'test-token');

    try {
      const flow = new Workflow('auto_recover_flow', 'Auto Recover Flow').start().end();

      // Seed the cache manually
      flow.withClient(client);
      // Run once with full AST
      await flow.run({ x: 1 });
      assert.equal(attempts, 1);

      // Now server restarts and drops def. Next run should attempt hot, get 404, and transparently retry cold
      const inst = await flow.run({ x: 2 });
      assert.equal(inst.instanceId, 'inst_recovered');
      assert.equal(attempts, 3, 'Should have attempted hot (404) and then cold (200)');
      assert.equal(requests[1].ast, undefined, 'Attempt 2 was hot');
      assert.ok(requests[2].ast, 'Attempt 3 was recovered cold with full AST');
    } finally {
      server.close();
    }
  });
});
