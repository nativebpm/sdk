import { test } from 'node:test';
import assert from 'node:assert';
import { z } from 'zod';
import {
  defineWorkflow,
  setDefaultSigningKey,
  getDefaultSigningKey,
} from '../dist/index.js';

test('Кейс 12: Code-First Native Workflows (defineWorkflow) & Worker/Topic Model (NB-282)', async (t) => {

  await t.test('1. Linear Async Workflow -> BPMN AST & XML Generation', async () => {
    // Developers write standard async functions with zero BPMN boilerplate
    const AddressSchema = z.object({
      city: z.string(),
      street: z.string(),
    });

    const linearFlow = defineWorkflow('linear_checkout', async ({ form, step }) => {
      const address = await form('checkout_address', {
        form: AddressSchema,
      });

      const scoring = await step('installment_scoring', {
        iin: '990101450123',
      });

      return { status: 'DONE' };
    });

    assert.strictEqual(linearFlow.id, 'linear_checkout');

    // 1. Check generated AST
    const ast = linearFlow.toAST();
    assert.strictEqual(ast.id, 'linear_checkout');
    assert.ok(ast.nodes.length >= 4, `Expected at least 4 nodes, got ${ast.nodes.length}`);

    const startNode = ast.nodes.find(n => n.type === 'startEvent');
    assert.ok(startNode, 'Should have startEvent node');

    const formNode = ast.nodes.find(n => n.type === 'userTask');
    assert.ok(formNode, 'Should have userTask node');
    assert.strictEqual(formNode?.formId, 'checkout_address');
    assert.ok(formNode?.inputSchema?.includes('city'), 'UserTask inputSchema should contain city from Zod schema');

    const stepNode = ast.nodes.find(n => n.type === 'serviceTask');
    assert.ok(stepNode, 'Should have serviceTask node');
    assert.strictEqual(stepNode?.topic, 'installment_scoring');

    const endNode = ast.nodes.find(n => n.type === 'endEvent');
    assert.ok(endNode, 'Should have endEvent node');

    // 2. Check Sequence Flows
    assert.ok(ast.flows.length >= 3, `Expected at least 3 flows, got ${ast.flows.length}`);
    const flow1 = ast.flows.find(f => f.source === startNode?.id && f.target === formNode?.id);
    assert.ok(flow1, 'Flow from start to userTask must exist');

    const flow2 = ast.flows.find(f => f.source === formNode?.id && f.target === stepNode?.id);
    assert.ok(flow2, 'Flow from userTask to serviceTask must exist');

    const flow3 = ast.flows.find(f => f.source === stepNode?.id && f.target === endNode?.id);
    assert.ok(flow3, 'Flow from serviceTask to endEvent must exist');

    // 3. Check BPMN 2.0 XML
    const xml = linearFlow.toBPMN();
    assert.ok(xml.includes('<bpmn:definitions'), 'XML should be standard BPMN 2.0 definitions');
    assert.ok(xml.includes('<bpmn:userTask'), 'XML should contain userTask');
    assert.ok(xml.includes('<bpmn:serviceTask'), 'XML should contain serviceTask');
    assert.ok(xml.includes('camunda:topic="installment_scoring"'), 'XML should configure camunda:topic');
  });

  await t.test('2. Branching Async Workflow (if/else) -> Exclusive Gateway & Conditional Flows', async () => {
    const branchFlow = defineWorkflow('branch_checkout', async ({ form, step }) => {
      const client = await form('client_profile');
      const scoring = await step('credit_scoring', { iin: client.iin });

      if (scoring.approved) {
        await step('reserve_stock');
        return { status: 'SUCCESS' };
      } else {
        await step('notify_rejection');
        return { status: 'REJECTED' };
      }
    });

    const ast = branchFlow.toAST();
    assert.strictEqual(ast.id, 'branch_checkout');

    // Check Gateway exists
    const gatewayNode = ast.nodes.find(n => n.type === 'exclusiveGateway');
    assert.ok(gatewayNode, 'Should automatically generate an exclusiveGateway for if/else branch');

    const reserveNode = ast.nodes.find(n => n.type === 'serviceTask' && n.topic === 'reserve_stock');
    assert.ok(reserveNode, 'Should have reserve_stock serviceTask');

    const rejectNode = ast.nodes.find(n => n.type === 'serviceTask' && n.topic === 'notify_rejection');
    assert.ok(rejectNode, 'Should have notify_rejection serviceTask');

    // Check conditional flows from gateway
    const trueFlow = ast.flows.find(f => f.source === gatewayNode?.id && f.target === reserveNode?.id);
    assert.ok(trueFlow, 'Conditional flow to reserve_stock must exist');
    assert.ok(trueFlow?.condition && trueFlow.condition.includes('approved'), 'Condition should reflect scoring.approved');

    const falseFlow = ast.flows.find(f => f.source === gatewayNode?.id && f.target === rejectNode?.id);
    assert.ok(falseFlow, 'Conditional flow to notify_rejection must exist');

    // Check BPMN XML contains exclusiveGateway
    const xml = branchFlow.toBPMN();
    assert.ok(xml.includes('<bpmn:exclusiveGateway'), 'XML must contain exclusiveGateway');
  });

  await t.test('3. Declarative State Machine Definition (AI Agent / XState friendly)', async () => {
    const stateMachineFlow = defineWorkflow({
      id: 'sm_loan_flow',
      initial: 'fill_form',
      states: {
        fill_form: {
          type: 'user_task',
          form: 'loan_application_form',
          on: { SUBMIT: 'scoring' }
        },
        scoring: {
          type: 'service_task',
          topic: 'client.scoring',
          on: {
            APPROVED: 'disburse_money',
            REJECTED: 'finish_rejected'
          }
        },
        disburse_money: {
          type: 'service_task',
          topic: 'client.disbursement',
          on: { DONE: 'finish_approved' }
        },
        finish_approved: { type: 'end' },
        finish_rejected: { type: 'end' }
      }
    });

    assert.strictEqual(stateMachineFlow.id, 'sm_loan_flow');
    const ast = stateMachineFlow.toAST();

    assert.strictEqual(ast.id, 'sm_loan_flow');
    assert.ok(ast.nodes.some(n => n.type === 'startEvent'));
    assert.ok(ast.nodes.some(n => n.id === 'fill_form' && n.type === 'userTask'));
    assert.ok(ast.nodes.some(n => n.id === 'scoring' && n.type === 'serviceTask' && n.topic === 'client.scoring'));
    assert.ok(ast.nodes.some(n => n.id === 'disburse_money' && n.type === 'serviceTask'));
    assert.ok(ast.nodes.some(n => n.id === 'finish_approved' && n.type === 'endEvent'));
    assert.ok(ast.nodes.some(n => n.id === 'finish_rejected' && n.type === 'endEvent'));

    // Check flows
    assert.ok(ast.flows.some(f => f.source === 'fill_form' && f.target === 'scoring'));
    assert.ok(ast.flows.some(f => f.source === 'scoring' && f.target === 'disburse_money'));
    assert.ok(ast.flows.some(f => f.source === 'scoring' && f.target === 'finish_rejected'));

    const xml = stateMachineFlow.toBPMN();
    assert.ok(xml.includes('id="fill_form"'));
    assert.ok(xml.includes('camunda:topic="client.scoring"'));
  });

  await t.test('4. One-Line Execution (.run) with Auto-Signature & Variable Handling', async () => {
    setDefaultSigningKey('test-secret-key-12345');
    assert.strictEqual(getDefaultSigningKey(), 'test-secret-key-12345');

    const simpleFlow = defineWorkflow('quick_run_flow', async ({ step }) => {
      await step('client.fast_step');
      return { ok: true };
    });

    let capturedPayload: any = null;
    let capturedHeaders: Record<string, string> = {};

    // Mock client to intercept JIT execution call
    const mockClient = {
      executeProcess: async (req: any, opts: any) => {
        capturedPayload = req;
        capturedHeaders = opts?.headers || {};
        return {
          instanceId: 'inst_mock_9999',
          status: 'ACTIVE',
          definitionId: req.definitionId,
          version: 1,
          variables: req.variables,
        };
      }
    };

    const res = await simpleFlow.run({ amount: 50000 }, { client: mockClient as any });

    assert.strictEqual(res.instanceId, 'inst_mock_9999');
    assert.ok(capturedPayload, 'executeProcess should have been called');
    assert.strictEqual(capturedPayload.definitionId, 'quick_run_flow');
    assert.strictEqual(capturedPayload.variables?.amount, 50000);
    assert.ok(capturedPayload.ast, 'AST should be provided in JIT execution');
    assert.ok(capturedHeaders['X-NativeBPM-Signature'], 'HMAC signature should be transmitted in header');
    assert.strictEqual(capturedPayload.signature, capturedHeaders['X-NativeBPM-Signature']);
  });
});
