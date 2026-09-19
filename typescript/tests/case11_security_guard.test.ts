import test from 'node:test';
import assert from 'node:assert/strict';
import http from 'node:http';
import {
  Workflow,
  computeWorkflowHash,
  computeWorkflowSignature,
  verifyWorkflowSignature,
  setDefaultSigningKey,
  clearDeployedHashCache,
  nativebpmVitePlugin,
  generateWorkflowSignaturesManifest,
  Client,
  setDefaultClient,
} from '../dist/index.js';

test('Кейс 11: Out-of-the-Box Security Guard for JIT Execution (NB-281)', async (t) => {
  await t.test('1. Topic Sandboxing & Verification', () => {
    const wfSafe = new Workflow('order_flow', 'Order Flow')
      .start()
      .service('charge_card', 'Charge Card', 'client.checkout')
      .end();

    assert.equal(wfSafe.toAST().nodes[1].topic, 'client.checkout');

    const wfAdmin = new Workflow('admin_flow', 'Admin Flow')
      .start()
      .service('drop_users', 'Drop Users', 'admin.users.delete')
      .end();

    assert.equal(wfAdmin.toAST().nodes[1].topic, 'admin.users.delete');
  });

  await t.test('2. Bundler HMAC-SHA256 Signatures & Manifest Generation', () => {
    const secret = 'prod-signing-secret-key-123';
    const wf1 = new Workflow('order_flow', 'Order Flow').start().end();
    const hash1 = wf1.getContentHash();

    const sig1 = computeWorkflowSignature(hash1, secret);
    assert.match(sig1, /^[a-f0-9]{64}$/, 'Signature must be a 64-char hex string');

    // Valid signature verifies true
    assert.equal(verifyWorkflowSignature(hash1, sig1, secret), true);

    // Tampered hash or wrong secret verifies false
    assert.equal(verifyWorkflowSignature('sha256:fakehash', sig1, secret), false);
    assert.equal(verifyWorkflowSignature(hash1, sig1, 'wrong-secret'), false);

    // Manifest generation
    const manifest = generateWorkflowSignaturesManifest([wf1], secret);
    assert.equal(manifest.version, '1.0.0');
    assert.equal(manifest.signatures['order_flow'].contentHash, hash1);
    assert.equal(manifest.signatures['order_flow'].signature, sig1);
  });

  await t.test('3. Client Transmits X-NativeBPM-Signature and Handles Security Responses', async () => {
    clearDeployedHashCache();
    setDefaultSigningKey('test-secret');

    let lastRequest: any = null;
    let lastHeaders: any = null;

    const mockServer = http.createServer((req, res) => {
      let body = '';
      req.on('data', chunk => { body += chunk; });
      req.on('end', () => {
        lastHeaders = req.headers;
        lastRequest = JSON.parse(body || '{}');

        // Security check simulation:
        // Check for forbidden topic
        const nodes = lastRequest.ast?.nodes || [];
        for (const n of nodes) {
          if (n.topic && (n.topic.startsWith('admin.') || n.topic.includes('refund'))) {
            res.writeHead(403, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'FORBIDDEN_TOPIC_IN_CLIENT_WORKFLOW' }));
            return;
          }
        }

        // Check for complexity
        if (nodes.length > 50) {
          res.writeHead(400, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: 'AST_NODE_LIMIT_EXCEEDED' }));
          return;
        }

        // Check signature header
        if (req.headers['x-nativebpm-signature'] === 'invalid') {
          res.writeHead(401, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: 'INVALID_WORKFLOW_SIGNATURE' }));
          return;
        }

        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({
          instanceId: 'inst-sec-123',
          definitionId: lastRequest.definitionId || lastRequest.ast?.id,
          version: 1,
          isNewVersionDeployed: true,
          status: 'ACTIVE',
        }));
      });
    });

    await new Promise<void>(resolve => mockServer.listen(0, '127.0.0.1', () => resolve()));
    const port = (mockServer.address() as any).port;
    const client = new Client(`http://127.0.0.1:${port}`, 'test-token');
    setDefaultClient(client);

    try {
      // 3.1 Normal run attaches signature header automatically
      const wfNormal = new Workflow('safe_flow', 'Safe Flow')
        .start()
        .service('task1', 'Client Task', 'client.do_something')
        .end();

      const inst = await wfNormal.run({ amount: 500 });
      assert.equal(inst.instanceId, 'inst-sec-123');
      assert.ok(lastHeaders['x-nativebpm-signature'], 'X-NativeBPM-Signature header must be present');
      assert.equal(lastRequest.signature, lastHeaders['x-nativebpm-signature']);

      // 3.2 Forbidden topic throws 403 Forbidden
      const wfForbidden = new Workflow('bad_flow', 'Bad Flow')
        .start()
        .service('hack', 'Drop DB', 'admin.wipe_database')
        .end();

      await assert.rejects(
        async () => {
          await wfForbidden.run();
        },
        (err: any) => {
          assert.equal(err.status, 403);
          assert.match(err.message, /403/);
          return true;
        }
      );

      // 3.3 Invalid signature throws 401 Unauthorized
      await assert.rejects(
        async () => {
          await wfNormal.run({}, { signature: 'invalid' });
        },
        (err: any) => {
          assert.equal(err.status, 401);
          assert.match(err.message, /401/);
          return true;
        }
      );
    } finally {
      setDefaultSigningKey(undefined);
      await new Promise<void>(resolve => mockServer.close(() => resolve()));
    }
  });
});
