import { WorkflowAST } from './schemas/workflow-ast.js';
import { Client } from './client.js';
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
export declare class NativeWorkflow {
    id: string;
    name: string;
    private fn?;
    private smConfig?;
    private cachedAST?;
    constructor(id: string, name: string, fnOrConfig: AsyncWorkflowFn | StateMachineConfig);
    toAST(): WorkflowAST;
    toBPMN(): string;
    compile(): WorkflowAST;
    run(variables?: Record<string, any>, options?: {
        client?: Client;
        signingKey?: string;
    }): Promise<ExecuteProcessResponse>;
    private compileStateMachine;
    private compileAsyncFunction;
}
export declare function defineWorkflow(idOrConfig: string | StateMachineConfig, fnOrConfig?: AsyncWorkflowFn | StateMachineConfig): NativeWorkflow;
