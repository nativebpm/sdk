// Zero-dependency AST Workflow builder with native Zod 4 & BPMN 2.0 XML generation
import { z } from 'zod';
import { WorkflowASTSchema, } from './schemas/workflow-ast.js';
export { WorkflowASTSchema };
export function serializeFormSchema(form) {
    if (!form)
        return '';
    if (typeof form === 'string')
        return form;
    // Native Zod 4 support
    if (form instanceof z.ZodType || (form && typeof form === 'object' && ('_def' in form || '_zod' in form || typeof form.safeParse === 'function'))) {
        try {
            const compiled = z.toJSONSchema(form);
            return JSON.stringify(compiled);
        }
        catch {
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
function populateNodeProperties(node, opts) {
    if (!opts)
        return;
    for (const [key, val] of Object.entries(opts)) {
        let targetKey = key;
        if (key === 'wasm')
            targetKey = 'wasmPath';
        if (key === 'resultVariable')
            targetKey = 'resultVar';
        if (key === 'form') {
            node.inputSchema = serializeFormSchema(val);
            if (!node.formId && opts.formId) {
                node.formId = opts.formId;
            }
            else if (!node.formId && opts.form_id) {
                node.formId = opts.form_id;
            }
            else if (!node.formId && node.id) {
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
            if (targetKey === 'wasm')
                targetKey = 'wasmPath';
            if (targetKey === 'resultVariable')
                targetKey = 'resultVar';
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
    hasEnded = false;
    workflow;
    gatewayID;
    currentNodeID;
    isConditional;
    condition;
    constructor(workflow, gatewayID, currentNodeID, isConditional, condition) {
        this.workflow = workflow;
        this.gatewayID = gatewayID;
        this.currentNodeID = currentNodeID;
        this.isConditional = isConditional;
        this.condition = condition;
    }
    connectNode(id) {
        if (this.hasEnded)
            return;
        const merges = this.workflow.pendingMerges;
        if (merges && merges.length > 0) {
            for (const sourceID of merges) {
                this.workflow.sequenceFlow(sourceID, id);
            }
            this.workflow.pendingMerges = [];
            this.currentNodeID = id;
            return;
        }
        if (this.currentNodeID === this.gatewayID) {
            if (this.isConditional) {
                this.workflow.sequenceFlowWithCondition(this.gatewayID, id, this.condition || '');
            }
            else {
                this.workflow.sequenceFlow(this.gatewayID, id);
            }
        }
        else if (this.currentNodeID && this.currentNodeID !== id) {
            this.workflow.sequenceFlow(this.currentNodeID, id);
        }
        this.currentNodeID = id;
    }
    user(id, name, opts) {
        this.workflow.userTask(id, name, opts);
        this.connectNode(id);
        return this;
    }
    service(id, name, topic, opts) {
        this.workflow.serviceTask(id, name, topic, opts);
        this.connectNode(id);
        return this;
    }
    ai(id, name, opts) {
        this.workflow.aiTask(id, name, opts);
        this.connectNode(id);
        return this;
    }
    call(id, name, calledElement, opts) {
        this.workflow.callActivity(id, name, calledElement, opts);
        this.connectNode(id);
        return this;
    }
    businessRule(id, name, decisionRef, opts) {
        this.workflow.businessRuleTask(id, name, decisionRef, opts);
        this.connectNode(id);
        return this;
    }
    boundaryTimer(id, name, attachedToRef, timeDuration, cancelActivity = true) {
        this.workflow.boundaryTimerEvent(id, name, attachedToRef, timeDuration, cancelActivity);
        this.connectNode(id);
        return this;
    }
    end(id, name) {
        this.workflow.endEvent(id, name);
        this.connectNode(id);
        this.hasEnded = true;
        return this;
    }
    when(condition) {
        const gwID = `gw_${this.currentNodeID}_decision`;
        this.workflow.exclusiveGateway(gwID, 'Decision Gateway');
        this.connectNode(gwID);
        const condStr = typeof condition === 'string' ? condition : condition.toString();
        return new WhenBranchBuilder(this, gwID, condStr);
    }
}
export class WhenBuilder {
    workflow;
    gatewayID;
    condition;
    constructor(workflow, gatewayID, condition) {
        this.workflow = workflow;
        this.gatewayID = gatewayID;
        this.condition = condition;
    }
    then(targetOrFn) {
        if (typeof targetOrFn === 'function') {
            const thenBranch = new Branch(this.workflow, this.gatewayID, this.gatewayID, true, this.condition);
            targetOrFn(thenBranch);
            if (!thenBranch.hasEnded && thenBranch.currentNodeID !== this.gatewayID) {
                this.workflow.pendingMerges.push(thenBranch.currentNodeID);
            }
            return new ThenBuilder(this.workflow, this.gatewayID);
        }
        else {
            const targetId = targetOrFn;
            this.workflow.sequenceFlowWithCondition(this.gatewayID, targetId, this.condition);
            this.workflow.currentNodeID = targetId;
            return new ThenBuilder(this.workflow, this.gatewayID, targetId);
        }
    }
}
export class ThenBuilder {
    workflow;
    gatewayID;
    currentTargetID;
    constructor(workflow, gatewayID, currentTargetID) {
        this.workflow = workflow;
        this.gatewayID = gatewayID;
        this.currentTargetID = currentTargetID;
    }
    userTask(id, name, opts) {
        this.workflow.userTask(id, name, opts);
        this.workflow.currentNodeID = id;
        this.currentTargetID = id;
        return this;
    }
    serviceTask(id, name, topic, opts) {
        this.workflow.serviceTask(id, name, topic, opts);
        this.workflow.currentNodeID = id;
        this.currentTargetID = id;
        return this;
    }
    aiTask(id, name, opts) {
        this.workflow.aiTask(id, name, opts);
        this.workflow.currentNodeID = id;
        this.currentTargetID = id;
        return this;
    }
    businessRuleTask(id, name, decisionRef, opts) {
        this.workflow.businessRuleTask(id, name, decisionRef, opts);
        this.workflow.currentNodeID = id;
        this.currentTargetID = id;
        return this;
    }
    else(targetOrFn) {
        if (this.currentTargetID) {
            this.workflow.pendingMerges.push(this.currentTargetID);
        }
        if (typeof targetOrFn === 'function') {
            const elseBranch = new Branch(this.workflow, this.gatewayID, this.gatewayID, false);
            targetOrFn(elseBranch);
            if (!elseBranch.hasEnded && elseBranch.currentNodeID !== this.gatewayID) {
                this.workflow.pendingMerges.push(elseBranch.currentNodeID);
            }
            return this.workflow;
        }
        else {
            const targetId = targetOrFn;
            this.workflow.sequenceFlow(this.gatewayID, targetId);
            this.workflow.currentNodeID = targetId;
            return this.workflow;
        }
    }
    otherwise(targetOrFn) {
        return this.else(targetOrFn);
    }
    when(condition) {
        return new WhenBuilder(this.workflow, this.gatewayID, String(condition));
    }
}
export class WhenBranchBuilder {
    branch;
    gatewayID;
    condition;
    constructor(branch, gatewayID, condition) {
        this.branch = branch;
        this.gatewayID = gatewayID;
        this.condition = condition;
    }
    then(targetOrFn) {
        if (typeof targetOrFn === 'function') {
            const thenBranch = new Branch(this.branch.workflow, this.gatewayID, this.gatewayID, true, this.condition);
            targetOrFn(thenBranch);
            if (!thenBranch.hasEnded && thenBranch.currentNodeID !== this.gatewayID) {
                this.branch.workflow.pendingMerges.push(thenBranch.currentNodeID);
            }
            return new ThenBranchBuilder(this.branch, this.gatewayID);
        }
        else {
            const targetId = targetOrFn;
            this.branch.workflow.sequenceFlowWithCondition(this.gatewayID, targetId, this.condition);
            this.branch.currentNodeID = targetId;
            return new ThenBranchBuilder(this.branch, this.gatewayID, targetId);
        }
    }
}
export class ThenBranchBuilder {
    branch;
    gatewayID;
    currentTargetID;
    constructor(branch, gatewayID, currentTargetID) {
        this.branch = branch;
        this.gatewayID = gatewayID;
        this.currentTargetID = currentTargetID;
    }
    userTask(id, name, opts) {
        this.branch.workflow.userTask(id, name, opts);
        this.branch.currentNodeID = id;
        this.currentTargetID = id;
        return this;
    }
    serviceTask(id, name, topic, opts) {
        this.branch.workflow.serviceTask(id, name, topic, opts);
        this.branch.currentNodeID = id;
        this.currentTargetID = id;
        return this;
    }
    else(targetOrFn) {
        if (typeof targetOrFn === 'function') {
            const elseBranch = new Branch(this.branch.workflow, this.gatewayID, this.gatewayID, false);
            targetOrFn(elseBranch);
            if (!elseBranch.hasEnded && elseBranch.currentNodeID !== this.gatewayID) {
                this.branch.workflow.pendingMerges.push(elseBranch.currentNodeID);
            }
            return this.branch;
        }
        else {
            const targetId = targetOrFn;
            this.branch.workflow.sequenceFlow(this.gatewayID, targetId);
            this.branch.currentNodeID = targetId;
            return this.branch;
        }
    }
    otherwise(targetOrFn) {
        return this.else(targetOrFn);
    }
    when(condition) {
        return new WhenBranchBuilder(this.branch, this.gatewayID, String(condition));
    }
}
export class Workflow {
    id;
    name;
    inputSchema;
    nodes = [];
    flows = [];
    currentNodeID = '';
    pendingMerges = [];
    constructor(id, name) {
        this.id = id;
        this.name = name;
    }
    variables(schema) {
        this.inputSchema = serializeFormSchema(schema);
        return this;
    }
    connectNode(id) {
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
        }
        else if (this.currentNodeID && this.currentNodeID !== id) {
            this.sequenceFlow(this.currentNodeID, id);
        }
        this.currentNodeID = id;
    }
    start(id = 'start') {
        this.startEvent(id);
        this.connectNode(id);
        return this;
    }
    end(id = 'end', name = 'End') {
        this.endEvent(id, name);
        this.connectNode(id);
        this.currentNodeID = '';
        return this;
    }
    user(id, name, opts) {
        this.userTask(id, name, opts);
        this.connectNode(id);
        return this;
    }
    service(id, name, topic, opts) {
        this.serviceTask(id, name, topic, opts);
        this.connectNode(id);
        return this;
    }
    ai(id, name, opts) {
        this.aiTask(id, name, opts);
        this.connectNode(id);
        return this;
    }
    call(id, name, calledElement, opts) {
        this.callActivity(id, name, calledElement, opts);
        this.connectNode(id);
        return this;
    }
    businessRule(id, name, decisionRef, opts) {
        this.businessRuleTask(id, name, decisionRef, opts);
        this.connectNode(id);
        return this;
    }
    when(condition) {
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
    else(targetOrFn) {
        if (typeof targetOrFn === 'function') {
            const elseBranch = new Branch(this, this.currentNodeID, this.currentNodeID, false);
            targetOrFn(elseBranch);
            if (!elseBranch.hasEnded && elseBranch.currentNodeID !== this.currentNodeID) {
                this.pendingMerges.push(elseBranch.currentNodeID);
            }
            return this;
        }
        else {
            const targetId = targetOrFn;
            this.sequenceFlow(this.currentNodeID, targetId);
            this.currentNodeID = targetId;
            return this;
        }
    }
    startEvent(id = 'start') {
        if (this.findNode(id)) {
            this.currentNodeID = id;
            return this;
        }
        const node = { type: 'startEvent', id, name: 'Start' };
        this.nodes.push(node);
        this.currentNodeID = id;
        return this;
    }
    endEvent(id, name) {
        if (this.findNode(id))
            return this;
        const node = { type: 'endEvent', id, name };
        this.nodes.push(node);
        return this;
    }
    serviceTask(id, name, topic, opts) {
        let node = this.findNode(id);
        if (!node) {
            node = { type: 'serviceTask', id, name, topic };
            populateNodeProperties(node, opts);
            this.nodes.push(node);
        }
        this.connectNode(id);
        return this;
    }
    aiTask(id, name, opts) {
        let node = this.findNode(id);
        if (!node) {
            node = { type: 'aiServiceTask', id, name };
            populateNodeProperties(node, opts);
            this.nodes.push(node);
        }
        this.connectNode(id);
        return this;
    }
    userTask(id, name, opts) {
        let node = this.findNode(id);
        if (!node) {
            node = { type: 'userTask', id, name };
            populateNodeProperties(node, opts);
            this.nodes.push(node);
        }
        this.connectNode(id);
        return this;
    }
    exclusiveGateway(id, name) {
        let node = this.findNode(id);
        if (!node) {
            node = { type: 'exclusiveGateway', id, name: name || id };
            this.nodes.push(node);
        }
        this.connectNode(id);
        return this;
    }
    parallelGateway(id, name) {
        let node = this.findNode(id);
        if (!node) {
            node = { type: 'parallelGateway', id, name: name || id };
            this.nodes.push(node);
        }
        this.connectNode(id);
        return this;
    }
    eventBasedGateway(id, name) {
        let node = this.findNode(id);
        if (!node) {
            node = { type: 'eventBasedGateway', id, name: name || id };
            this.nodes.push(node);
        }
        this.connectNode(id);
        return this;
    }
    boundaryTimer(id, name, attachedToRefOrDuration, timeDuration, cancelActivity = true) {
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
                z.iso.duration().parse(duration);
            }
            catch (err) {
                throw new Error(`Invalid ISO 8601 duration "${duration}": ${err.message || err}`);
            }
        }
        return this.boundaryTimerEvent(id, name, attached, duration, cancel);
    }
    boundaryTimerEvent(id, name, attachedToRef, timeDuration, cancelActivity = true) {
        if (timeDuration) {
            try {
                z.iso.duration().parse(timeDuration);
            }
            catch (err) {
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
    callActivity(id, name, calledElement, opts) {
        let node = this.findNode(id);
        if (!node) {
            node = { type: 'callActivity', id, name, calledElement };
            populateNodeProperties(node, opts);
            this.nodes.push(node);
        }
        this.connectNode(id);
        return this;
    }
    businessRuleTask(id, name, decisionRef, opts) {
        let node = this.findNode(id);
        if (!node) {
            node = { type: 'businessRuleTask', id, name, decisionRef };
            populateNodeProperties(node, opts);
            this.nodes.push(node);
        }
        this.connectNode(id);
        return this;
    }
    sequenceFlow(source, target, condition) {
        if (condition) {
            return this.sequenceFlowWithCondition(source, target, condition);
        }
        const existing = this.flows.find(f => f.source === source && f.target === target && !f.condition);
        if (existing)
            return this;
        this.flows.push({ id: `flow-${source}-${target}`, source, target, condition: '' });
        return this;
    }
    sequenceFlowWithCondition(source, target, condition) {
        const existing = this.flows.find(f => f.source === source && f.target === target && f.condition === condition);
        if (existing)
            return this;
        this.flows.push({ id: `flow-${source}-${target}`, source, target, condition });
        return this;
    }
    findNode(id) {
        return this.nodes.find(n => n.id === id);
    }
    toAST() {
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
    toJSON() {
        return JSON.stringify(this.toAST());
    }
    toBPMN() {
        return generateBPMNXML(this.toAST());
    }
    extractForms() {
        const forms = {};
        for (const node of this.nodes) {
            if (node.type === 'userTask') {
                const formId = node.formId || node.id;
                if (node.inputSchema) {
                    try {
                        forms[formId] = typeof node.inputSchema === 'string'
                            ? JSON.parse(node.inputSchema)
                            : node.inputSchema;
                    }
                    catch {
                        forms[formId] = node.inputSchema;
                    }
                }
            }
        }
        return forms;
    }
}
export class WorkflowBuilder extends Workflow {
}
function escapeXml(str) {
    if (!str)
        return '';
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&apos;');
}
export function generateBPMNXML(ast) {
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
    if (ast.inputSchema) {
        const pSchema = typeof ast.inputSchema === 'string' ? ast.inputSchema : JSON.stringify(ast.inputSchema);
        processAttrs += ` nativebpm:inputSchema="${escapeXml(pSchema)}"`;
    }
    xml += `  <bpmn:process${processAttrs}>\n`;
    const incomingMap = {};
    const outgoingMap = {};
    for (const f of ast.flows) {
        if (!outgoingMap[f.source])
            outgoingMap[f.source] = [];
        outgoingMap[f.source].push(f.id);
        if (!incomingMap[f.target])
            incomingMap[f.target] = [];
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
                if (n.assignee)
                    attrs += ` camunda:assignee="${escapeXml(n.assignee)}"`;
                if (n.candidateGroups)
                    attrs += ` camunda:candidateGroups="${escapeXml(n.candidateGroups)}"`;
                if (n.formId)
                    attrs += ` camunda:formKey="${escapeXml(n.formId)}"`;
                else if (n.formKey)
                    attrs += ` camunda:formKey="${escapeXml(n.formKey)}"`;
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
                if (n.topic)
                    attrs += ` camunda:type="external" camunda:topic="${escapeXml(n.topic)}" nativebpm:topic="${escapeXml(n.topic)}"`;
                if (n.wasmPath)
                    attrs += ` nativebpm:wasmPath="${escapeXml(n.wasmPath)}"`;
                xml += `    <bpmn:serviceTask${attrs}>\n`;
                xml += incomings;
                xml += outgoings;
                xml += `    </bpmn:serviceTask>\n`;
                break;
            }
            case 'aiServiceTask': {
                let attrs = ` id="${escapeXml(n.id)}" name="${escapeXml(n.name)}"`;
                if (n.provider)
                    attrs += ` nativebpm:provider="${escapeXml(n.provider)}"`;
                if (n.model)
                    attrs += ` nativebpm:model="${escapeXml(n.model)}"`;
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
                if (n.decisionRef)
                    attrs += ` camunda:decisionRef="${escapeXml(n.decisionRef)}"`;
                if (n.resultVar)
                    attrs += ` camunda:resultVariable="${escapeXml(n.resultVar)}"`;
                if (n.mapDecisionResult)
                    attrs += ` camunda:mapDecisionResult="${escapeXml(n.mapDecisionResult)}"`;
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
        }
        else {
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
export function evaluateDMNRule(node, variables) {
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
            }
            else if (expected.startsWith('>=') || expected.startsWith('<=') || expected.startsWith('>') || expected.startsWith('<') || expected.startsWith('==') || expected.startsWith('!=')) {
                const numVal = Number(actualValue);
                const opMatch = expected.match(/^([><=!]+)\s*(.+)$/);
                if (opMatch) {
                    const op = opMatch[1];
                    const threshold = Number(opMatch[2]);
                    if (op === '>=' && !(numVal >= threshold)) {
                        matches = false;
                        break;
                    }
                    if (op === '<=' && !(numVal <= threshold)) {
                        matches = false;
                        break;
                    }
                    if (op === '>' && !(numVal > threshold)) {
                        matches = false;
                        break;
                    }
                    if (op === '<' && !(numVal < threshold)) {
                        matches = false;
                        break;
                    }
                    if (op === '==' && !(numVal === threshold)) {
                        matches = false;
                        break;
                    }
                    if (op === '!=' && !(numVal !== threshold)) {
                        matches = false;
                        break;
                    }
                }
            }
            else {
                if (String(actualValue) !== expected) {
                    matches = false;
                    break;
                }
            }
        }
        if (matches) {
            const result = {};
            for (let o = 0; o < rule.outputs.length; o++) {
                const outName = outputs[o]?.name || node.resultVar || `output_${o}`;
                const outType = outputs[o]?.type || 'string';
                let rawOut = rule.outputs[o];
                if (rawOut.startsWith('"') && rawOut.endsWith('"')) {
                    rawOut = rawOut.slice(1, -1);
                }
                if (outType === 'number') {
                    result[outName] = Number(rawOut);
                }
                else if (outType === 'boolean') {
                    result[outName] = rawOut === 'true';
                }
                else {
                    result[outName] = rawOut;
                }
            }
            return result;
        }
    }
    return {};
}
export class Expression extends String {
}
export class Variable {
    name;
    constructor(name) {
        this.name = name;
    }
    eq(value) {
        const valStr = typeof value === 'boolean' ? (value ? 'true' : 'false') : String(value);
        return new Expression(`${this.name} == ${valStr}`);
    }
    ne(value) {
        const valStr = typeof value === 'boolean' ? (value ? 'true' : 'false') : String(value);
        return new Expression(`${this.name} != ${valStr}`);
    }
    gt(value) {
        return new Expression(`${this.name} > ${value}`);
    }
    gte(value) {
        return new Expression(`${this.name} >= ${value}`);
    }
    lt(value) {
        return new Expression(`${this.name} < ${value}`);
    }
    lte(value) {
        return new Expression(`${this.name} <= ${value}`);
    }
}
export function V(name) {
    return new Variable(name);
}
export function v(name) {
    return new Variable(name);
}
