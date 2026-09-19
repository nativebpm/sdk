// Native Code-First Workflow definitions for NativeBPM (NB-282)
import {
  WorkflowAST,
  NodeAST,
  FlowAST,
} from './schemas/workflow-ast.js';
import {
  serializeFormSchema,
  generateBPMNXML,
  getDefaultSigningKey,
  computeWorkflowHash,
  computeWorkflowSignature,
} from './builder.js';
import { getDefaultClient, Client } from './client.js';
import type { ExecuteProcessResponse } from './api/src/models/ExecuteProcessResponse.js';

export interface FormOptions {
  form?: any;
  schema?: any;
  title?: string;
  formId?: string;
  [key: string]: any;
}

export interface StepOptions {
  [key: string]: any;
}

export interface WorkflowStepContext {
  step<T = any>(topic: string, payload?: StepOptions): Promise<T>;
  form<T = any>(formId: string, options?: FormOptions): Promise<T>;
  parallel?<T extends readonly unknown[] | []>(steps: T): Promise<T>;
}

export type AsyncWorkflowFn = (ctx: WorkflowStepContext) => Promise<any> | any;

export interface StateConfig {
  type: 'user_task' | 'service_task' | 'end';
  topic?: string;
  form?: any;
  on?: Record<string, string>;
}

export interface StateMachineConfig {
  id: string;
  name?: string;
  initial: string;
  states: Record<string, StateConfig>;
}

interface ParsedAction {
  kind: 'form' | 'step' | 'return';
  value: string;
}

function extractBracedBlock(src: string, startIndex: number): { content: string; endIndex: number } {
  let depth = 0;
  let blockStart = -1;
  for (let i = startIndex; i < src.length; i++) {
    if (src[i] === '{') {
      if (depth === 0) blockStart = i + 1;
      depth++;
    } else if (src[i] === '}') {
      depth--;
      if (depth === 0) {
        return { content: src.slice(blockStart, i), endIndex: i + 1 };
      }
    }
  }
  return { content: '', endIndex: src.length };
}

function extractActions(block: string): ParsedAction[] {
  if (!block) return [];
  const regex = /(?:(?:await\s+)?(?:ctx\.)?form\s*\(\s*['"`]([^'"`]+)['"`]|(?:await\s+)?(?:ctx\.)?step\s*\(\s*['"`]([^'"`]+)['"`]|return\b)/g;
  const actions: ParsedAction[] = [];
  let match;
  while ((match = regex.exec(block)) !== null) {
    if (match[1]) {
      actions.push({ kind: 'form', value: match[1] });
    } else if (match[2]) {
      actions.push({ kind: 'step', value: match[2] });
    } else {
      actions.push({ kind: 'return', value: 'return' });
    }
  }
  return actions;
}

export class NativeWorkflow {
  public id: string;
  public name: string;
  private fn?: AsyncWorkflowFn;
  private smConfig?: StateMachineConfig;
  private cachedAST?: WorkflowAST;

  constructor(id: string, name: string, fnOrConfig: AsyncWorkflowFn | StateMachineConfig) {
    this.id = id;
    this.name = name;
    if (typeof fnOrConfig === 'function') {
      this.fn = fnOrConfig;
    } else {
      this.smConfig = fnOrConfig;
    }
  }

  public toAST(): WorkflowAST {
    if (this.cachedAST) {
      return this.cachedAST;
    }

    if (this.smConfig) {
      this.cachedAST = this.compileStateMachine(this.smConfig);
      return this.cachedAST;
    }

    if (this.fn) {
      this.cachedAST = this.compileAsyncFunction(this.id, this.name, this.fn);
      return this.cachedAST;
    }

    throw new Error(`NativeWorkflow ${this.id} has no function or state machine configuration`);
  }

  public toBPMN(): string {
    return generateBPMNXML(this.toAST());
  }

  public compile(): WorkflowAST {
    return this.toAST();
  }

  public async run(
    variables: Record<string, any> = {},
    options?: { client?: Client; signingKey?: string }
  ): Promise<ExecuteProcessResponse> {
    const ast = this.toAST();
    const signingKey = options?.signingKey || getDefaultSigningKey();

    let signature: string | undefined;
    const headers: Record<string, string> = {};

    if (signingKey) {
      const hash = computeWorkflowHash(ast);
      signature = computeWorkflowSignature(hash, signingKey);
      headers['X-NativeBPM-Signature'] = signature;
    }

    const client = options?.client || getDefaultClient();
    if (!client) {
      throw new Error('NativeBPM Client not configured. Pass { client } to run() or setDefaultClient(client).');
    }

    if (typeof (client as any).executeProcess === 'function') {
      return (client as any).executeProcess(
        {
          definitionId: this.id,
          variables,
          ast,
          signature,
        },
        { headers }
      );
    }

    if (typeof (client as any).execute === 'function') {
      return (client as any).execute({
        definitionId: this.id,
        variables,
        ast,
        signature,
      });
    }

    throw new Error('Client must provide execute or executeProcess method');
  }

  private compileStateMachine(config: StateMachineConfig): WorkflowAST {
    const nodes: NodeAST[] = [];
    const flows: FlowAST[] = [];

    // 1. Start event
    nodes.push({
      type: 'startEvent',
      id: 'start',
      name: 'Start',
      topic: '',
      wasmPath: '',
      provider: '',
      model: '',
      prompt: '',
      systemInstruction: '',
      responseSchema: '',
      temperature: null,
      resultVar: '',
      assignee: '',
      candidateGroups: '',
      dueDate: '',
      calledElement: '',
      decisionRef: '',
      mapDecisionResult: '',
      inVariables: [],
      outVariables: [],
      hitPolicy: '',
      inputs: [],
      outputs: [],
      rules: [],
    });

    if (config.initial) {
      flows.push({
        id: `flow_start_to_${config.initial}`,
        source: 'start',
        target: config.initial,
        condition: '',
      });
    }

    // 2. States
    for (const [stateId, state] of Object.entries(config.states)) {
      let nodeType: NodeAST['type'] = 'serviceTask';
      let formId = '';
      let inputSchema = '';
      let topic = '';

      if (state.type === 'user_task') {
        nodeType = 'userTask';
        formId = typeof state.form === 'string' ? state.form : stateId;
        inputSchema = serializeFormSchema(state.form);
      } else if (state.type === 'service_task') {
        nodeType = 'serviceTask';
        topic = state.topic || stateId;
      } else if (state.type === 'end') {
        nodeType = 'endEvent';
      }

      nodes.push({
        type: nodeType,
        id: stateId,
        name: stateId,
        topic,
        formId,
        inputSchema,
        wasmPath: '',
        provider: '',
        model: '',
        prompt: '',
        systemInstruction: '',
        responseSchema: '',
        temperature: null,
        resultVar: '',
        assignee: '',
        candidateGroups: '',
        dueDate: '',
        calledElement: '',
        decisionRef: '',
        mapDecisionResult: '',
        inVariables: [],
        outVariables: [],
        hitPolicy: '',
        inputs: [],
        outputs: [],
        rules: [],
      });

      // Flows from transitions
      if (state.on) {
        for (const [event, targetState] of Object.entries(state.on)) {
          flows.push({
            id: `flow_${stateId}_to_${targetState}`,
            source: stateId,
            target: targetState,
            condition: event === 'DEFAULT' || event === 'DONE' || event === 'SUBMIT' ? '' : event,
          });
        }
      }
    }

    return {
      id: config.id,
      name: config.name || config.id,
      nodes,
      flows,
    };
  }

  private compileAsyncFunction(id: string, name: string, fn: AsyncWorkflowFn): WorkflowAST {
    const recordedForms: Record<string, FormOptions> = {};
    const recordedSteps: Record<string, StepOptions> = {};

    // 1. Dry run trace to capture runtime object references (Zod schemas, etc.)
    const tracerCtx: WorkflowStepContext = {
      form: async <T = any>(formId: string, options?: FormOptions): Promise<T> => {
        recordedForms[formId] = options || {};
        return new Proxy({}, { get: (_, prop) => prop }) as unknown as T;
      },
      step: async <T = any>(topic: string, payload?: StepOptions): Promise<T> => {
        recordedSteps[topic] = payload || {};
        return new Proxy({}, { get: (_, prop) => prop }) as unknown as T;
      },
    };

    try {
      const res = fn(tracerCtx);
      if (res && typeof res.then === 'function') {
        // synchronously drain if already resolved
        res.catch(() => {});
      }
    } catch {
      // ignore runtime mock errors during dry run
    }

    // 2. Parse function string
    const fnStr = fn.toString();
    const nodes: NodeAST[] = [];
    const flows: FlowAST[] = [];

    // Helper to create empty NodeAST
    const createNode = (opts: Partial<NodeAST> & { id: string; type: NodeAST['type'] }): NodeAST => ({
      name: opts.name || opts.id,
      topic: '',
      wasmPath: '',
      provider: '',
      model: '',
      prompt: '',
      systemInstruction: '',
      responseSchema: '',
      temperature: null,
      resultVar: '',
      assignee: '',
      candidateGroups: '',
      dueDate: '',
      calledElement: '',
      decisionRef: '',
      mapDecisionResult: '',
      inVariables: [],
      outVariables: [],
      hitPolicy: '',
      inputs: [],
      outputs: [],
      rules: [],
      ...opts,
    });

    // Start Event
    nodes.push(createNode({ id: 'start', name: 'Start', type: 'startEvent' }));

    const ifRegex = /\bif\s*\(/;
    const ifMatch = fnStr.match(ifRegex);

    let preBlock = fnStr;
    let hasIf = false;
    let conditionExpr = '';
    let thenBlock = '';
    let elseBlock = '';

    if (ifMatch && ifMatch.index !== undefined) {
      hasIf = true;
      const ifIdx = ifMatch.index;
      preBlock = fnStr.slice(0, ifIdx);

      const condStart = fnStr.indexOf('(', ifIdx);
      const condEnd = fnStr.indexOf(')', condStart);
      conditionExpr = fnStr.slice(condStart + 1, condEnd).trim();

      const thenRes = extractBracedBlock(fnStr, condEnd);
      thenBlock = thenRes.content;

      const afterThen = fnStr.slice(thenRes.endIndex);
      const elseMatch = afterThen.match(/^\s*else\s*/);
      if (elseMatch) {
        const elseRes = extractBracedBlock(afterThen, elseMatch[0].length - 1);
        elseBlock = elseRes.content;
      }
    }

    let seqIdx = 1;

    // Helper to process actions list into nodes and flows (REFACTOR-001)
    const processActionList = (
      actions: ParsedAction[],
      startFromId: string,
      firstCondition: string,
      endName: string
    ): string => {
      let currentId = startFromId;
      for (const action of actions) {
        if (action.kind === 'step') {
          const topic = action.value;
          const nodeId = `step_${topic.replace(/[^a-zA-Z0-9_]/g, '_')}`;
          nodes.push(createNode({
            id: nodeId,
            name: topic,
            type: 'serviceTask',
            topic,
          }));

          flows.push({
            id: `flow_${currentId}_to_${nodeId}`,
            source: currentId,
            target: nodeId,
            condition: currentId === startFromId ? firstCondition : '',
          });
          currentId = nodeId;
        } else if (action.kind === 'form') {
          const formId = action.value;
          const nodeId = `form_${formId}`;
          const formOpts = recordedForms[formId] || {};
          const inputSchema = serializeFormSchema(formOpts.form || formOpts.schema);

          nodes.push(createNode({
            id: nodeId,
            name: formId,
            type: 'userTask',
            formId,
            inputSchema,
          }));

          flows.push({
            id: `flow_${currentId}_to_${nodeId}`,
            source: currentId,
            target: nodeId,
            condition: currentId === startFromId ? firstCondition : '',
          });
          currentId = nodeId;
        } else if (action.kind === 'return') {
          const endId = `end_${endName.toLowerCase()}_${seqIdx++}`;
          nodes.push(createNode({ id: endId, name: endName, type: 'endEvent' }));
          flows.push({
            id: `flow_${currentId}_to_${endId}`,
            source: currentId,
            target: endId,
            condition: currentId === startFromId ? firstCondition : '',
          });
          currentId = endId;
        }
      }
      return currentId;
    };

    const preActions = extractActions(preBlock);
    const lastPreId = processActionList(preActions, 'start', '', 'End');

    // Handle If / Else branching
    if (hasIf) {
      const gatewayId = `gateway_${seqIdx++}`;
      nodes.push(createNode({
        id: gatewayId,
        name: conditionExpr,
        type: 'exclusiveGateway',
      }));

      flows.push({
        id: `flow_${lastPreId}_to_${gatewayId}`,
        source: lastPreId,
        target: gatewayId,
        condition: '',
      });

      // THEN Branch
      const thenActions = extractActions(thenBlock);
      processActionList(thenActions, gatewayId, `\${${conditionExpr}}`, 'Success');

      // ELSE Branch
      if (elseBlock) {
        const elseActions = extractActions(elseBlock);
        processActionList(elseActions, gatewayId, `!\${${conditionExpr}}`, 'Rejected');
      }
    }

    return {
      id,
      name,
      nodes,
      flows,
    };
  }
}

export function defineWorkflow(
  idOrConfig: string | StateMachineConfig,
  fnOrConfig?: AsyncWorkflowFn | StateMachineConfig
): NativeWorkflow {
  if (typeof idOrConfig === 'object') {
    return new NativeWorkflow(idOrConfig.id, idOrConfig.name || idOrConfig.id, idOrConfig);
  }

  if (typeof fnOrConfig === 'function') {
    return new NativeWorkflow(idOrConfig, idOrConfig, fnOrConfig);
  }

  if (typeof fnOrConfig === 'object') {
    return new NativeWorkflow(idOrConfig, fnOrConfig.name || idOrConfig, fnOrConfig);
  }

  throw new Error(`Invalid defineWorkflow arguments for ${idOrConfig}`);
}
