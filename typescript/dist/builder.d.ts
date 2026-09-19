import { WorkflowASTSchema, type InVariableAST as InVariable, type OutVariableAST as OutVariable, type NodeAST, type FlowAST, type WorkflowAST } from './schemas/workflow-ast.js';
export { InVariable, OutVariable, NodeAST, FlowAST, WorkflowAST, WorkflowASTSchema };
export declare function serializeFormSchema(form: any): string;
export declare class Branch {
    hasEnded: boolean;
    workflow: Workflow;
    gatewayID: string;
    currentNodeID: string;
    isConditional: boolean;
    condition?: string;
    constructor(workflow: Workflow, gatewayID: string, currentNodeID: string, isConditional: boolean, condition?: string);
    private connectNode;
    user(id: string, name: string, opts?: Record<string, any>): Branch;
    service(id: string, name: string, topic: string, opts?: Record<string, any>): Branch;
    ai(id: string, name: string, opts?: Record<string, any>): Branch;
    call(id: string, name: string, calledElement: string, opts?: Record<string, any>): Branch;
    businessRule(id: string, name: string, decisionRef: string, opts?: Record<string, any>): Branch;
    boundaryTimer(id: string, name: string, attachedToRef: string, timeDuration: string, cancelActivity?: boolean): Branch;
    end(id: string, name: string): Branch;
    when(condition: string | {
        toString(): string;
    }): WhenBranchBuilder;
}
export declare class WhenBuilder {
    workflow: Workflow;
    gatewayID: string;
    condition: string;
    constructor(workflow: Workflow, gatewayID: string, condition: string);
    then(targetOrFn: string | ((flow: Branch) => void)): ThenBuilder;
}
export declare class ThenBuilder {
    workflow: Workflow;
    gatewayID: string;
    currentTargetID?: string;
    constructor(workflow: Workflow, gatewayID: string, currentTargetID?: string);
    userTask(id: string, name: string, opts?: Record<string, any>): ThenBuilder;
    serviceTask(id: string, name: string, topic: string, opts?: Record<string, any>): ThenBuilder;
    aiTask(id: string, name: string, opts?: Record<string, any>): ThenBuilder;
    businessRuleTask(id: string, name: string, decisionRef: string, opts?: Record<string, any>): ThenBuilder;
    else(targetOrFn: string | ((flow: Branch) => void)): Workflow;
    otherwise(targetOrFn: string | ((flow: Branch) => void)): Workflow;
    when(condition: string | {
        toString(): string;
    }): WhenBuilder;
}
export declare class WhenBranchBuilder {
    branch: Branch;
    gatewayID: string;
    condition: string;
    constructor(branch: Branch, gatewayID: string, condition: string);
    then(targetOrFn: string | ((flow: Branch) => void)): ThenBranchBuilder;
}
export declare class ThenBranchBuilder {
    branch: Branch;
    gatewayID: string;
    currentTargetID?: string;
    constructor(branch: Branch, gatewayID: string, currentTargetID?: string);
    userTask(id: string, name: string, opts?: Record<string, any>): ThenBranchBuilder;
    serviceTask(id: string, name: string, topic: string, opts?: Record<string, any>): ThenBranchBuilder;
    else(targetOrFn: string | ((flow: Branch) => void)): Branch;
    otherwise(targetOrFn: string | ((flow: Branch) => void)): Branch;
    when(condition: string | {
        toString(): string;
    }): WhenBranchBuilder;
}
export declare class Workflow {
    id: string;
    name: string;
    inputSchema?: string | Record<string, any>;
    nodes: NodeAST[];
    flows: FlowAST[];
    currentNodeID: string;
    private pendingMerges;
    constructor(id: string, name: string);
    variables(schema: any): this;
    private connectNode;
    start(id?: string): Workflow;
    end(id?: string, name?: string): Workflow;
    user(id: string, name: string, opts?: Record<string, any>): Workflow;
    service(id: string, name: string, topic: string, opts?: Record<string, any>): Workflow;
    ai(id: string, name: string, opts?: Record<string, any>): Workflow;
    call(id: string, name: string, calledElement: string, opts?: Record<string, any>): Workflow;
    businessRule(id: string, name: string, decisionRef: string, opts?: Record<string, any>): Workflow;
    when(condition: string | {
        toString(): string;
    }): WhenBuilder;
    else(targetOrFn: string | ((flow: Branch) => void)): Workflow;
    startEvent(id?: string): Workflow;
    endEvent(id: string, name: string): Workflow;
    serviceTask(id: string, name: string, topic: string, opts?: Record<string, any>): Workflow;
    aiTask(id: string, name: string, opts?: Record<string, any>): Workflow;
    userTask(id: string, name: string, opts?: Record<string, any>): Workflow;
    exclusiveGateway(id: string, name?: string): Workflow;
    parallelGateway(id: string, name?: string): Workflow;
    eventBasedGateway(id: string, name?: string): Workflow;
    boundaryTimer(id: string, name: string, attachedToRefOrDuration?: string, timeDuration?: string, cancelActivity?: boolean): Workflow;
    boundaryTimerEvent(id: string, name: string, attachedToRef: string, timeDuration: string, cancelActivity?: boolean): Workflow;
    callActivity(id: string, name: string, calledElement: string, opts?: Record<string, any>): Workflow;
    businessRuleTask(id: string, name: string, decisionRef: string, opts?: Record<string, any>): Workflow;
    sequenceFlow(source: string, target: string, condition?: string): Workflow;
    sequenceFlowWithCondition(source: string, target: string, condition: string): Workflow;
    findNode(id: string): NodeAST | undefined;
    toAST(): WorkflowAST;
    toJSON(): string;
    toBPMN(): string;
    extractForms(): Record<string, any>;
    private clientInstance?;
    withClient(client: any): this;
    getClient(): any;
    getContentHash(): string;
    run<TVariables extends Record<string, any> = Record<string, any>>(variables?: TVariables, options?: {
        client?: any;
        businessKey?: string;
        forms?: Record<string, any>;
        baseUrl?: string;
        apiToken?: string;
    }): Promise<ProcessInstanceHandle>;
}
export interface ProcessInstanceHandle {
    instanceId: string;
    definitionId: string;
    version: number;
    isNewVersionDeployed: boolean;
    status: string;
    state?: Record<string, any>;
    currentTasks: any[];
    client?: any;
    claimTask?(taskId: string, assignee: string): Promise<any>;
    completeTask?(taskId: string, variables?: Record<string, any>): Promise<any>;
}
export declare function canonicalJsonStringify(obj: any): string;
export declare function computeWorkflowHash(ast: WorkflowAST | string): string;
export declare function clearDeployedHashCache(): void;
export declare class WorkflowBuilder extends Workflow {
}
export declare function generateBPMNXML(ast: WorkflowAST): string;
export declare function evaluateDMNRule(node: NodeAST, variables: Record<string, any>): Record<string, any>;
export declare class Expression extends String {
}
export declare class Variable {
    name: string;
    constructor(name: string);
    eq(value: any): Expression;
    ne(value: any): Expression;
    gt(value: any): Expression;
    gte(value: any): Expression;
    lt(value: any): Expression;
    lte(value: any): Expression;
}
export declare function V(name: string): Variable;
export declare function v(name: string): Variable;
