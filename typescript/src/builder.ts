// Zero-dependency AST Workflow builder with native Zod 4 & BPMN 2.0 XML generation
import { z } from 'zod';
import {
  WorkflowASTSchema,
  type InVariableAST as InVariable,
  type OutVariableAST as OutVariable,
  type NodeAST,
  type FlowAST,
  type WorkflowAST,
} from './schemas/workflow-ast.js';

export { InVariable, OutVariable, NodeAST, FlowAST, WorkflowAST, WorkflowASTSchema };

export function serializeFormSchema(form: any): string {
  if (!form) return '';
  if (typeof form === 'string') return form;

  // Native Zod 4 support
  if (form instanceof z.ZodType || (form && typeof form === 'object' && ('_def' in form || '_zod' in form || typeof form.safeParse === 'function'))) {
    try {
      const compiled = (z as any).toJSONSchema(form);
      return JSON.stringify(compiled);
    } catch {
      // fallback
    }
  }

  if (typeof form.toJSONSchema === 'function') {
    const res = form.toJSONSchema();
    return typeof res === 'string' ? res : JSON.stringify(res);
  }
  if (typeof form === 'object') {
    return JSON.stringify(form);
  }
  return String(form);
}

function populateNodeProperties(node: any, opts?: Record<string, any>): void {
  if (!opts) return;
  for (const [key, val] of Object.entries(opts)) {
    let targetKey = key;
    if (key === 'wasm') targetKey = 'wasmPath';
    if (key === 'resultVariable') targetKey = 'resultVar';
    if (key === 'form') {
      node.inputSchema = serializeFormSchema(val);
      if (!node.formId && opts.formId) {
        node.formId = opts.formId;
      } else if (!node.formId && opts.form_id) {
        node.formId = opts.form_id;
      } else if (!node.formId && node.id) {
        node.formId = node.id;
      }
      continue;
    }
    if (key === 'responseSchema' || key === 'response_schema') {
      node.responseSchema = serializeFormSchema(val);
      continue;
    }
    if (key.includes('_')) {
      targetKey = key.replace(/_([a-z])/g, (g) => g[1].toUpperCase());
      if (targetKey === 'wasm') targetKey = 'wasmPath';
      if (targetKey === 'resultVariable') targetKey = 'resultVar';
    }
    node[targetKey] = val;
  }

  if (node.inputSchema && typeof node.inputSchema === 'object') {
    node.inputSchema = serializeFormSchema(node.inputSchema);
  }
  if (node.responseSchema && typeof node.responseSchema === 'object') {
    node.responseSchema = serializeFormSchema(node.responseSchema);
  }
}

export class Branch {
  public hasEnded: boolean = false;
  public workflow: Workflow;
  public gatewayID: string;
  public currentNodeID: string;
  public isConditional: boolean;
  public condition?: string;

  constructor(
    workflow: Workflow,
    gatewayID: string,
    currentNodeID: string,
    isConditional: boolean,
    condition?: string
  ) {
    this.workflow = workflow;
    this.gatewayID = gatewayID;
    this.currentNodeID = currentNodeID;
    this.isConditional = isConditional;
    this.condition = condition;
  }

  private connectNode(id: string): void {
    if (this.hasEnded) return;

    const merges = (this.workflow as any).pendingMerges as string[];
    if (merges && merges.length > 0) {
      for (const sourceID of merges) {
        this.workflow.sequenceFlow(sourceID, id);
      }
      (this.workflow as any).pendingMerges = [];
      this.currentNodeID = id;
      return;
    }

    if (this.currentNodeID === this.gatewayID) {
      if (this.isConditional) {
        this.workflow.sequenceFlowWithCondition(this.gatewayID, id, this.condition || '');
      } else {
        this.workflow.sequenceFlow(this.gatewayID, id);
      }
    } else if (this.currentNodeID && this.currentNodeID !== id) {
      this.workflow.sequenceFlow(this.currentNodeID, id);
    }
    this.currentNodeID = id;
  }

  public user(id: string, name: string, opts?: Record<string, any>): Branch {
    this.workflow.userTask(id, name, opts);
    this.connectNode(id);
    return this;
  }

  public service(id: string, name: string, topic: string, opts?: Record<string, any>): Branch {
    this.workflow.serviceTask(id, name, topic, opts);
    this.connectNode(id);
    return this;
  }

  public ai(id: string, name: string, opts?: Record<string, any>): Branch {
    this.workflow.aiTask(id, name, opts);
    this.connectNode(id);
    return this;
  }

  public call(id: string, name: string, calledElement: string, opts?: Record<string, any>): Branch {
    this.workflow.callActivity(id, name, calledElement, opts);
    this.connectNode(id);
    return this;
  }

  public businessRule(id: string, name: string, decisionRef: string, opts?: Record<string, any>): Branch {
    this.workflow.businessRuleTask(id, name, decisionRef, opts);
    this.connectNode(id);
    return this;
  }

  public boundaryTimer(
    id: string,
    name: string,
    attachedToRef: string,
    timeDuration: string,
    cancelActivity: boolean = true
  ): Branch {
    this.workflow.boundaryTimerEvent(id, name, attachedToRef, timeDuration, cancelActivity);
    this.connectNode(id);
    return this;
  }

  public end(id: string, name: string): Branch {
    this.workflow.endEvent(id, name);
    this.connectNode(id);
    this.hasEnded = true;
    return this;
  }

  public when(condition: string | { toString(): string }): WhenBranchBuilder {
    const gwID = `gw_${this.currentNodeID}_decision`;
    this.workflow.exclusiveGateway(gwID, 'Decision Gateway');
    this.connectNode(gwID);

    const condStr = typeof condition === 'string' ? condition : condition.toString();
    return new WhenBranchBuilder(this, gwID, condStr);
  }
}

export class WhenBuilder {
  public workflow: Workflow;
  public gatewayID: string;
  public condition: string;

  constructor(workflow: Workflow, gatewayID: string, condition: string) {
    this.workflow = workflow;
    this.gatewayID = gatewayID;
    this.condition = condition;
  }

  public then(targetOrFn: string | ((flow: Branch) => void)): ThenBuilder {
    if (typeof targetOrFn === 'function') {
      const thenBranch = new Branch(this.workflow, this.gatewayID, this.gatewayID, true, this.condition);
      targetOrFn(thenBranch);

      if (!thenBranch.hasEnded && thenBranch.currentNodeID !== this.gatewayID) {
        ((this.workflow as any).pendingMerges as string[]).push(thenBranch.currentNodeID);
      }

      return new ThenBuilder(this.workflow, this.gatewayID);
    } else {
      const targetId = targetOrFn;
      this.workflow.sequenceFlowWithCondition(this.gatewayID, targetId, this.condition);
      (this.workflow as any).currentNodeID = targetId;
      return new ThenBuilder(this.workflow, this.gatewayID, targetId);
    }
  }
}

export class ThenBuilder {
  public workflow: Workflow;
  public gatewayID: string;
  public currentTargetID?: string;

  constructor(workflow: Workflow, gatewayID: string, currentTargetID?: string) {
    this.workflow = workflow;
    this.gatewayID = gatewayID;
    this.currentTargetID = currentTargetID;
  }

  public userTask(id: string, name: string, opts?: Record<string, any>): ThenBuilder {
    this.workflow.userTask(id, name, opts);
    (this.workflow as any).currentNodeID = id;
    this.currentTargetID = id;
    return this;
  }

  public serviceTask(id: string, name: string, topic: string, opts?: Record<string, any>): ThenBuilder {
    this.workflow.serviceTask(id, name, topic, opts);
    (this.workflow as any).currentNodeID = id;
    this.currentTargetID = id;
    return this;
  }

  public aiTask(id: string, name: string, opts?: Record<string, any>): ThenBuilder {
    this.workflow.aiTask(id, name, opts);
    (this.workflow as any).currentNodeID = id;
    this.currentTargetID = id;
    return this;
  }

  public businessRuleTask(id: string, name: string, decisionRef: string, opts?: Record<string, any>): ThenBuilder {
    this.workflow.businessRuleTask(id, name, decisionRef, opts);
    (this.workflow as any).currentNodeID = id;
    this.currentTargetID = id;
    return this;
  }

  public else(targetOrFn: string | ((flow: Branch) => void)): Workflow {
    if (this.currentTargetID) {
      ((this.workflow as any).pendingMerges as string[]).push(this.currentTargetID);
    }
    if (typeof targetOrFn === 'function') {
      const elseBranch = new Branch(this.workflow, this.gatewayID, this.gatewayID, false);
      targetOrFn(elseBranch);

      if (!elseBranch.hasEnded && elseBranch.currentNodeID !== this.gatewayID) {
        ((this.workflow as any).pendingMerges as string[]).push(elseBranch.currentNodeID);
      }

      return this.workflow;
    } else {
      const targetId = targetOrFn;
      this.workflow.sequenceFlow(this.gatewayID, targetId);
      (this.workflow as any).currentNodeID = targetId;
      return this.workflow;
    }
  }

  public otherwise(targetOrFn: string | ((flow: Branch) => void)): Workflow {
    return this.else(targetOrFn);
  }

  public when(condition: string | { toString(): string }): WhenBuilder {
    return new WhenBuilder(this.workflow, this.gatewayID, String(condition));
  }
}

export class WhenBranchBuilder {
  public branch: Branch;
  public gatewayID: string;
  public condition: string;

  constructor(branch: Branch, gatewayID: string, condition: string) {
    this.branch = branch;
    this.gatewayID = gatewayID;
    this.condition = condition;
  }

  public then(targetOrFn: string | ((flow: Branch) => void)): ThenBranchBuilder {
    if (typeof targetOrFn === 'function') {
      const thenBranch = new Branch(this.branch.workflow, this.gatewayID, this.gatewayID, true, this.condition);
      targetOrFn(thenBranch);

      if (!thenBranch.hasEnded && thenBranch.currentNodeID !== this.gatewayID) {
        ((this.branch.workflow as any).pendingMerges as string[]).push(thenBranch.currentNodeID);
      }

      return new ThenBranchBuilder(this.branch, this.gatewayID);
    } else {
      const targetId = targetOrFn;
      this.branch.workflow.sequenceFlowWithCondition(this.gatewayID, targetId, this.condition);
      this.branch.currentNodeID = targetId;
      return new ThenBranchBuilder(this.branch, this.gatewayID, targetId);
    }
  }
}

export class ThenBranchBuilder {
  public branch: Branch;
  public gatewayID: string;
  public currentTargetID?: string;

  constructor(branch: Branch, gatewayID: string, currentTargetID?: string) {
    this.branch = branch;
    this.gatewayID = gatewayID;
    this.currentTargetID = currentTargetID;
  }

  public userTask(id: string, name: string, opts?: Record<string, any>): ThenBranchBuilder {
    this.branch.workflow.userTask(id, name, opts);
    this.branch.currentNodeID = id;
    this.currentTargetID = id;
    return this;
  }

  public serviceTask(id: string, name: string, topic: string, opts?: Record<string, any>): ThenBranchBuilder {
    this.branch.workflow.serviceTask(id, name, topic, opts);
    this.branch.currentNodeID = id;
    this.currentTargetID = id;
    return this;
  }

  public else(targetOrFn: string | ((flow: Branch) => void)): Branch {
    if (typeof targetOrFn === 'function') {
      const elseBranch = new Branch(this.branch.workflow, this.gatewayID, this.gatewayID, false);
      targetOrFn(elseBranch);

      if (!elseBranch.hasEnded && elseBranch.currentNodeID !== this.gatewayID) {
        ((this.branch.workflow as any).pendingMerges as string[]).push(elseBranch.currentNodeID);
      }

      return this.branch;
    } else {
      const targetId = targetOrFn;
      this.branch.workflow.sequenceFlow(this.gatewayID, targetId);
      this.branch.currentNodeID = targetId;
      return this.branch;
    }
  }

  public otherwise(targetOrFn: string | ((flow: Branch) => void)): Branch {
    return this.else(targetOrFn);
  }

  public when(condition: string | { toString(): string }): WhenBranchBuilder {
    return new WhenBranchBuilder(this.branch, this.gatewayID, String(condition));
  }
}

export class Workflow {
  public id: string;
  public name: string;
  public inputSchema?: string | Record<string, any>;
  public nodes: NodeAST[] = [];
  public flows: FlowAST[] = [];
  public currentNodeID: string = '';
  private pendingMerges: string[] = [];

  constructor(id: string, name: string) {
    this.id = id;
    this.name = name;
  }

  public variables(schema: any): this {
    this.inputSchema = serializeFormSchema(schema);
    return this;
  }

  private connectNode(id: string): void {
    const node = this.findNode(id);
    const hasStart = this.nodes.some(n => n.type === 'startEvent');
    if (!hasStart && node && node.type !== 'startEvent') {
      this.startEvent('start');
      this.sequenceFlow('start', id);
      this.currentNodeID = id;
      return;
    }

    if (this.pendingMerges.length > 0) {
      for (const sourceID of this.pendingMerges) {
        this.sequenceFlow(sourceID, id);
      }
      this.pendingMerges = [];
    } else if (this.currentNodeID && this.currentNodeID !== id) {
      this.sequenceFlow(this.currentNodeID, id);
    }
    this.currentNodeID = id;
  }

  public start(id: string = 'start'): Workflow {
    this.startEvent(id);
    this.connectNode(id);
    return this;
  }

  public end(id: string = 'end', name: string = 'End'): Workflow {
    this.endEvent(id, name);
    this.connectNode(id);
    this.currentNodeID = '';
    return this;
  }

  public user(id: string, name: string, opts?: Record<string, any>): Workflow {
    this.userTask(id, name, opts);
    this.connectNode(id);
    return this;
  }

  public service(id: string, name: string, topic: string, opts?: Record<string, any>): Workflow {
    this.serviceTask(id, name, topic, opts);
    this.connectNode(id);
    return this;
  }

  public ai(id: string, name: string, opts?: Record<string, any>): Workflow {
    this.aiTask(id, name, opts);
    this.connectNode(id);
    return this;
  }

  public call(id: string, name: string, calledElement: string, opts?: Record<string, any>): Workflow {
    this.callActivity(id, name, calledElement, opts);
    this.connectNode(id);
    return this;
  }

  public businessRule(id: string, name: string, decisionRef: string, opts?: Record<string, any>): Workflow {
    this.businessRuleTask(id, name, decisionRef, opts);
    this.connectNode(id);
    return this;
  }

  public when(condition: string | { toString(): string }): WhenBuilder {
    const isCurrentGateway = this.nodes.some(n => n.id === this.currentNodeID && n.type === 'exclusiveGateway');
    const gwID = isCurrentGateway
      ? this.currentNodeID
      : `gw_${this.currentNodeID}_decision`;

    if (!isCurrentGateway) {
      this.exclusiveGateway(gwID, 'Decision Gateway');
      this.connectNode(gwID);
    }

    const condStr = typeof condition === 'string' ? condition : condition.toString();
    return new WhenBuilder(this, gwID, condStr);
  }

  public else(targetOrFn: string | ((flow: Branch) => void)): Workflow {
    if (typeof targetOrFn === 'function') {
      const elseBranch = new Branch(this, this.currentNodeID, this.currentNodeID, false);
      targetOrFn(elseBranch);
      if (!elseBranch.hasEnded && elseBranch.currentNodeID !== this.currentNodeID) {
        this.pendingMerges.push(elseBranch.currentNodeID);
      }
      return this;
    } else {
      const targetId = targetOrFn;
      this.sequenceFlow(this.currentNodeID, targetId);
      this.currentNodeID = targetId;
      return this;
    }
  }

  public startEvent(id: string = 'start'): Workflow {
    if (this.findNode(id)) {
      this.currentNodeID = id;
      return this;
    }
    const node: NodeAST = { type: 'startEvent', id, name: 'Start' };
    this.nodes.push(node);
    this.currentNodeID = id;
    return this;
  }

  public endEvent(id: string, name: string): Workflow {
    if (this.findNode(id)) return this;
    const node: NodeAST = { type: 'endEvent', id, name };
    this.nodes.push(node);
    return this;
  }

  public serviceTask(id: string, name: string, topic: string, opts?: Record<string, any>): Workflow {
    let node = this.findNode(id);
    if (!node) {
      node = { type: 'serviceTask', id, name, topic };
      populateNodeProperties(node, opts);
      this.nodes.push(node);
    }
    this.connectNode(id);
    return this;
  }

  public aiTask(id: string, name: string, opts?: Record<string, any>): Workflow {
    let node = this.findNode(id);
    if (!node) {
      node = { type: 'aiServiceTask', id, name };
      populateNodeProperties(node, opts);
      this.nodes.push(node);
    }
    this.connectNode(id);
    return this;
  }

  public userTask(id: string, name: string, opts?: Record<string, any>): Workflow {
    let node = this.findNode(id);
    if (!node) {
      node = { type: 'userTask', id, name };
      populateNodeProperties(node, opts);
      this.nodes.push(node);
    }
    this.connectNode(id);
    return this;
  }

  public exclusiveGateway(id: string, name?: string): Workflow {
    let node = this.findNode(id);
    if (!node) {
      node = { type: 'exclusiveGateway', id, name: name || id };
      this.nodes.push(node);
    }
    this.connectNode(id);
    return this;
  }

  public parallelGateway(id: string, name?: string): Workflow {
    let node = this.findNode(id);
    if (!node) {
      node = { type: 'parallelGateway', id, name: name || id };
      this.nodes.push(node);
    }
    this.connectNode(id);
    return this;
  }

  public eventBasedGateway(id: string, name?: string): Workflow {
    let node = this.findNode(id);
    if (!node) {
      node = { type: 'eventBasedGateway', id, name: name || id };
      this.nodes.push(node);
    }
    this.connectNode(id);
    return this;
  }

  public boundaryTimer(
    id: string,
    name: string,
    attachedToRefOrDuration?: string,
    timeDuration?: string,
    cancelActivity: boolean = true
  ): Workflow {
    let attached = attachedToRefOrDuration || '';
    let duration = timeDuration || '';
    let cancel = cancelActivity;

    if (timeDuration === undefined && attachedToRefOrDuration !== undefined) {
      duration = attachedToRefOrDuration;
      attached = this.currentNodeID;
      cancel = true;
    }

    if (duration) {
      try {
        (z as any).iso.duration().parse(duration);
      } catch (err: any) {
        throw new Error(`Invalid ISO 8601 duration "${duration}": ${err.message || err}`);
      }
    }

    return this.boundaryTimerEvent(id, name, attached, duration, cancel);
  }

  public boundaryTimerEvent(
    id: string,
    name: string,
    attachedToRef: string,
    timeDuration: string,
    cancelActivity: boolean = true
  ): Workflow {
    if (timeDuration) {
      try {
        (z as any).iso.duration().parse(timeDuration);
      } catch (err: any) {
        throw new Error(`Invalid ISO 8601 duration "${timeDuration}": ${err.message || err}`);
      }
    }
    let node = this.findNode(id);
    if (!node) {
      node = {
        type: 'boundaryTimerEvent',
        id,
        name,
        attachedToRef,
        timeDuration,
        cancelActivity
      };
      this.nodes.push(node);
    }
    this.currentNodeID = id;
    return this;
  }

  public callActivity(id: string, name: string, calledElement: string, opts?: Record<string, any>): Workflow {
    let node = this.findNode(id);
    if (!node) {
      node = { type: 'callActivity', id, name, calledElement };
      populateNodeProperties(node, opts);
      this.nodes.push(node);
    }
    this.connectNode(id);
    return this;
  }

  public businessRuleTask(id: string, name: string, decisionRef: string, opts?: Record<string, any>): Workflow {
    let node = this.findNode(id);
    if (!node) {
      node = { type: 'businessRuleTask', id, name, decisionRef };
      populateNodeProperties(node, opts);
      this.nodes.push(node);
    }
    this.connectNode(id);
    return this;
  }

  public sequenceFlow(source: string, target: string, condition?: string): Workflow {
    if (condition) {
      return this.sequenceFlowWithCondition(source, target, condition);
    }
    const existing = this.flows.find(f => f.source === source && f.target === target && !f.condition);
    if (existing) return this;
    this.flows.push({ id: `flow-${source}-${target}`, source, target, condition: '' });
    return this;
  }

  public sequenceFlowWithCondition(source: string, target: string, condition: string): Workflow {
    const existing = this.flows.find(f => f.source === source && f.target === target && f.condition === condition);
    if (existing) return this;
    this.flows.push({ id: `flow-${source}-${target}`, source, target, condition });
    return this;
  }

  public findNode(id: string): NodeAST | undefined {
    return this.nodes.find(n => n.id === id);
  }

  public toAST(): WorkflowAST {
    const nodes = [...this.nodes];
    const flows = [...this.flows];

    const sourceIDs = new Set(flows.map(f => f.source));

    for (const node of this.nodes) {
      if (node.type === 'endEvent' || node.type === 'startEvent') {
        continue;
      }
      if (!sourceIDs.has(node.id)) {
        const endID = `end_${node.id}`;
        nodes.push({ type: 'endEvent', id: endID, name: 'Process Finished' });
        flows.push({ id: `flow-${node.id}-${endID}`, source: node.id, target: endID, condition: '' });
      }
    }

    return {
      id: this.id,
      name: this.name,
      ...(this.inputSchema ? { inputSchema: this.inputSchema } : {}),
      nodes,
      flows
    };
  }

  public toJSON(): string {
    return JSON.stringify(this.toAST());
  }

  public toBPMN(): string {
    return generateBPMNXML(this.toAST());
  }

  public extractForms(): Record<string, any> {
    const forms: Record<string, any> = {};
    for (const node of this.nodes) {
      if (node.type === 'userTask') {
        const formId = node.formId || node.id;
        if (node.inputSchema) {
          try {
            forms[formId] = typeof node.inputSchema === 'string'
              ? JSON.parse(node.inputSchema)
              : node.inputSchema;
          } catch {
            forms[formId] = node.inputSchema;
          }
        }
      }
    }
    return forms;
  }
}

export class WorkflowBuilder extends Workflow {}

function escapeXml(str: string): string {
  if (!str) return '';
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&apos;');
}

export function generateBPMNXML(ast: WorkflowAST): string {
  let xml = `<?xml version="1.0" encoding="UTF-8"?>\n`;
  xml += `<bpmn:definitions xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL"\n`;
  xml += `  xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI"\n`;
  xml += `  xmlns:dc="http://www.omg.org/spec/DD/20100524/DC"\n`;
  xml += `  xmlns:di="http://www.omg.org/spec/DD/20100524/DI"\n`;
  xml += `  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"\n`;
  xml += `  xmlns:camunda="http://camunda.org/schema/1.0/bpmn"\n`;
  xml += `  xmlns:nativebpm="https://nativebpm.com/schema/1.0/bpmn"\n`;
  xml += `  id="Definitions_${escapeXml(ast.id)}"\n`;
  xml += `  targetNamespace="http://bpmn.io/schema/bpmn"\n`;
  xml += `  exporter="NativeBPM TypeScript SDK"\n`;
  xml += `  exporterVersion="1.0.0">\n`;
  let processAttrs = ` id="${escapeXml(ast.id)}" name="${escapeXml(ast.name)}" isExecutable="true"`;
  if ((ast as any).inputSchema) {
    const pSchema = typeof (ast as any).inputSchema === 'string' ? (ast as any).inputSchema : JSON.stringify((ast as any).inputSchema);
    processAttrs += ` nativebpm:inputSchema="${escapeXml(pSchema)}"`;
  }
  xml += `  <bpmn:process${processAttrs}>\n`;

  const incomingMap: Record<string, string[]> = {};
  const outgoingMap: Record<string, string[]> = {};

  for (const f of ast.flows) {
    if (!outgoingMap[f.source]) outgoingMap[f.source] = [];
    outgoingMap[f.source].push(f.id);

    if (!incomingMap[f.target]) incomingMap[f.target] = [];
    incomingMap[f.target].push(f.id);
  }

  for (const n of ast.nodes) {
    const incomings = (incomingMap[n.id] || []).map(id => `      <bpmn:incoming>${escapeXml(id)}</bpmn:incoming>\n`).join('');
    const outgoings = (outgoingMap[n.id] || []).map(id => `      <bpmn:outgoing>${escapeXml(id)}</bpmn:outgoing>\n`).join('');

    switch (n.type) {
      case 'startEvent':
        xml += `    <bpmn:startEvent id="${escapeXml(n.id)}" name="${escapeXml(n.name)}">\n`;
        xml += outgoings;
        xml += `    </bpmn:startEvent>\n`;
        break;

      case 'endEvent':
        xml += `    <bpmn:endEvent id="${escapeXml(n.id)}" name="${escapeXml(n.name)}">\n`;
        xml += incomings;
        xml += `    </bpmn:endEvent>\n`;
        break;

      case 'userTask': {
        let attrs = ` id="${escapeXml(n.id)}" name="${escapeXml(n.name)}"`;
        if (n.assignee) attrs += ` camunda:assignee="${escapeXml(n.assignee)}"`;
        if (n.candidateGroups) attrs += ` camunda:candidateGroups="${escapeXml(n.candidateGroups)}"`;
        if (n.formId) attrs += ` camunda:formKey="${escapeXml(n.formId)}"`;
        else if (n.formKey) attrs += ` camunda:formKey="${escapeXml(n.formKey)}"`;
        if (n.inputSchema) {
          const schemaStr = typeof n.inputSchema === 'string' ? n.inputSchema : JSON.stringify(n.inputSchema);
          attrs += ` inputSchema="${escapeXml(schemaStr)}"`;
        }
        xml += `    <bpmn:userTask${attrs}>\n`;
        xml += incomings;
        xml += outgoings;
        xml += `    </bpmn:userTask>\n`;
        break;
      }

      case 'serviceTask': {
        let attrs = ` id="${escapeXml(n.id)}" name="${escapeXml(n.name)}"`;
        if (n.topic) attrs += ` camunda:type="external" camunda:topic="${escapeXml(n.topic)}" nativebpm:topic="${escapeXml(n.topic)}"`;
        if (n.wasmPath) attrs += ` nativebpm:wasmPath="${escapeXml(n.wasmPath)}"`;
        xml += `    <bpmn:serviceTask${attrs}>\n`;
        xml += incomings;
        xml += outgoings;
        xml += `    </bpmn:serviceTask>\n`;
        break;
      }

      case 'aiServiceTask': {
        let attrs = ` id="${escapeXml(n.id)}" name="${escapeXml(n.name)}"`;
        if (n.provider) attrs += ` nativebpm:provider="${escapeXml(n.provider)}"`;
        if (n.model) attrs += ` nativebpm:model="${escapeXml(n.model)}"`;
        if (n.responseSchema) {
          const respStr = typeof n.responseSchema === 'string' ? n.responseSchema : JSON.stringify(n.responseSchema);
          attrs += ` nativebpm:responseSchema="${escapeXml(respStr)}"`;
        }
        xml += `    <bpmn:serviceTask${attrs}>\n`;
        xml += incomings;
        xml += outgoings;
        xml += `    </bpmn:serviceTask>\n`;
        break;
      }

      case 'businessRuleTask': {
        let attrs = ` id="${escapeXml(n.id)}" name="${escapeXml(n.name)}"`;
        if (n.decisionRef) attrs += ` camunda:decisionRef="${escapeXml(n.decisionRef)}"`;
        if (n.resultVar) attrs += ` camunda:resultVariable="${escapeXml(n.resultVar)}"`;
        if (n.mapDecisionResult) attrs += ` camunda:mapDecisionResult="${escapeXml(n.mapDecisionResult)}"`;
        xml += `    <bpmn:businessRuleTask${attrs}>\n`;
        xml += incomings;
        xml += outgoings;
        xml += `    </bpmn:businessRuleTask>\n`;
        break;
      }

      case 'exclusiveGateway':
        xml += `    <bpmn:exclusiveGateway id="${escapeXml(n.id)}" name="${escapeXml(n.name)}">\n`;
        xml += incomings;
        xml += outgoings;
        xml += `    </bpmn:exclusiveGateway>\n`;
        break;

      case 'parallelGateway':
        xml += `    <bpmn:parallelGateway id="${escapeXml(n.id)}" name="${escapeXml(n.name)}">\n`;
        xml += incomings;
        xml += outgoings;
        xml += `    </bpmn:parallelGateway>\n`;
        break;

      case 'eventBasedGateway':
        xml += `    <bpmn:eventBasedGateway id="${escapeXml(n.id)}" name="${escapeXml(n.name)}">\n`;
        xml += incomings;
        xml += outgoings;
        xml += `    </bpmn:eventBasedGateway>\n`;
        break;

      case 'callActivity':
        xml += `    <bpmn:callActivity id="${escapeXml(n.id)}" name="${escapeXml(n.name)}" calledElement="${escapeXml(n.calledElement || '')}">\n`;
        xml += incomings;
        xml += outgoings;
        xml += `    </bpmn:callActivity>\n`;
        break;

      case 'boundaryTimerEvent': {
        const cancelAct = n.cancelActivity !== false;
        xml += `    <bpmn:boundaryEvent id="${escapeXml(n.id)}" name="${escapeXml(n.name)}" attachedToRef="${escapeXml(n.attachedToRef || '')}" cancelActivity="${cancelAct}">\n`;
        xml += outgoings;
        xml += `      <bpmn:timerEventDefinition id="timerDef_${escapeXml(n.id)}">\n`;
        xml += `        <bpmn:timeDuration xsi:type="bpmn:tFormalExpression">${escapeXml(n.timeDuration || '')}</bpmn:timeDuration>\n`;
        xml += `      </bpmn:timerEventDefinition>\n`;
        xml += `    </bpmn:boundaryEvent>\n`;
        break;
      }
    }
  }

  for (const f of ast.flows) {
    if (f.condition) {
      xml += `    <bpmn:sequenceFlow id="${escapeXml(f.id)}" sourceRef="${escapeXml(f.source)}" targetRef="${escapeXml(f.target)}">\n`;
      xml += `      <bpmn:conditionExpression xsi:type="bpmn:tFormalExpression">${escapeXml(f.condition)}</bpmn:conditionExpression>\n`;
      xml += `    </bpmn:sequenceFlow>\n`;
    } else {
      xml += `    <bpmn:sequenceFlow id="${escapeXml(f.id)}" sourceRef="${escapeXml(f.source)}" targetRef="${escapeXml(f.target)}" />\n`;
    }
  }

  xml += `  </bpmn:process>\n`;
  xml += `  <bpmndi:BPMNDiagram id="BPMNDiagram_1">\n`;
  xml += `    <bpmndi:BPMNPlane id="BPMNPlane_1" bpmnElement="${escapeXml(ast.id)}">\n`;
  xml += `    </bpmndi:BPMNPlane>\n`;
  xml += `  </bpmndi:BPMNDiagram>\n`;
  xml += `</bpmn:definitions>\n`;
  return xml;
}

export function evaluateDMNRule(
  node: NodeAST,
  variables: Record<string, any>
): Record<string, any> {
  if (!node.rules || node.rules.length === 0) {
    return {};
  }
  const inputs = node.inputs || [];
  const outputs = node.outputs || [];

  for (const rule of node.rules) {
    let matches = true;
    for (let i = 0; i < rule.inputs.length; i++) {
      const expr = inputs[i]?.expression;
      const rawExpected = rule.inputs[i].trim();
      const expected = (rawExpected.startsWith('"') && rawExpected.endsWith('"'))
        ? rawExpected.slice(1, -1)
        : rawExpected;
      if (!expr || expected === '-' || expected === '') {
        continue;
      }

      const actualValue = variables[expr];
      if (rawExpected.startsWith('"') && rawExpected.endsWith('"')) {
        if (String(actualValue) !== expected) {
          matches = false;
          break;
        }
      } else if (expected.startsWith('>=') || expected.startsWith('<=') || expected.startsWith('>') || expected.startsWith('<') || expected.startsWith('==') || expected.startsWith('!=')) {
        const numVal = Number(actualValue);
        const opMatch = expected.match(/^([><=!]+)\s*(.+)$/);
        if (opMatch) {
          const op = opMatch[1];
          const threshold = Number(opMatch[2]);
          if (op === '>=' && !(numVal >= threshold)) { matches = false; break; }
          if (op === '<=' && !(numVal <= threshold)) { matches = false; break; }
          if (op === '>' && !(numVal > threshold)) { matches = false; break; }
          if (op === '<' && !(numVal < threshold)) { matches = false; break; }
          if (op === '==' && !(numVal === threshold)) { matches = false; break; }
          if (op === '!=' && !(numVal !== threshold)) { matches = false; break; }
        }
      } else {
        if (String(actualValue) !== expected) {
          matches = false;
          break;
        }
      }
    }

    if (matches) {
      const result: Record<string, any> = {};
      for (let o = 0; o < rule.outputs.length; o++) {
        const outName = outputs[o]?.name || node.resultVar || `output_${o}`;
        const outType = outputs[o]?.type || 'string';
        let rawOut = rule.outputs[o];
        if (rawOut.startsWith('"') && rawOut.endsWith('"')) {
          rawOut = rawOut.slice(1, -1);
        }
        if (outType === 'number') {
          result[outName] = Number(rawOut);
        } else if (outType === 'boolean') {
          result[outName] = rawOut === 'true';
        } else {
          result[outName] = rawOut;
        }
      }
      return result;
    }
  }

  return {};
}

export class Expression extends String {}

export class Variable {
  public name: string;

  constructor(name: string) {
    this.name = name;
  }

  public eq(value: any): Expression {
    const valStr = typeof value === 'boolean' ? (value ? 'true' : 'false') : String(value);
    return new Expression(`${this.name} == ${valStr}`);
  }
  public ne(value: any): Expression {
    const valStr = typeof value === 'boolean' ? (value ? 'true' : 'false') : String(value);
    return new Expression(`${this.name} != ${valStr}`);
  }
  public gt(value: any): Expression {
    return new Expression(`${this.name} > ${value}`);
  }
  public gte(value: any): Expression {
    return new Expression(`${this.name} >= ${value}`);
  }
  public lt(value: any): Expression {
    return new Expression(`${this.name} < ${value}`);
  }
  public lte(value: any): Expression {
    return new Expression(`${this.name} <= ${value}`);
  }
}

export function V(name: string): Variable {
  return new Variable(name);
}

export function v(name: string): Variable {
  return new Variable(name);
}
