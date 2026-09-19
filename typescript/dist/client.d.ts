import * as api from "./api/src/index.js";
import { Workflow } from "./builder.js";
export type ProcessDefinition = api.ProcessDefinition;
export type ProcessInstance = api.ProcessInstance;
export type HistoryRecord = api.HistoryRecord;
export type IncidentRecord = api.IncidentRecord;
export type TaskRecord = api.TaskRecord;
export type WebhookRecord = api.WebhookRecord;
export type WebhookDeliveryRecord = api.WebhookDeliveryRecord;
export type VisualizationData = api.VisualizationData;
export type ExecuteProcessRequest = api.ExecuteProcessRequest;
export type ExecuteProcessResponse = api.ExecuteProcessResponse;
export declare function setDefaultClient(client: Client): void;
export declare function getDefaultClient(): Client | undefined;
export declare class Client {
    private baseUrl;
    private apiToken;
    constructor(baseUrl: string, apiToken: string);
    getHeaders(): Record<string, string>;
    getBaseUrl(): string;
    deploy(workflow: Workflow): Promise<ProcessDefinition>;
    definitions(): DefinitionsService;
    instances(): InstancesService;
    tasks(): TasksService;
    webhooks(): WebhooksService;
    execute(request: api.ExecuteProcessRequest): Promise<api.ExecuteProcessResponse>;
    process(): ProcessService;
}
export declare class ProcessService {
    private client;
    constructor(client: Client);
    execute(): ExecuteProcessBuilder;
}
export declare class ExecuteProcessBuilder {
    private client;
    private req;
    constructor(client: Client);
    withRequest(req: api.ExecuteProcessRequest): this;
    withAST(ast: any): this;
    withForms(forms: Record<string, any>): this;
    withContentHash(hash: string): this;
    withDefinitionId(id: string): this;
    withBusinessKey(key: string): this;
    withVariables(vars: Record<string, any>): this;
    send(): Promise<api.ExecuteProcessResponse>;
}
export declare class DefinitionsService {
    private client;
    constructor(client: Client);
    list(): ListDefinitionsBuilder;
    get(id: string): GetDefinitionBuilder;
    delete(id: string): DeleteDefinitionBuilder;
    deploy(): DeployDefinitionBuilder;
}
export declare class ListDefinitionsBuilder {
    private client;
    constructor(client: Client);
    send(): Promise<ProcessDefinition[]>;
}
export declare class GetDefinitionBuilder {
    private client;
    private id;
    constructor(client: Client, id: string);
    send(): Promise<ProcessDefinition>;
}
export declare class DeleteDefinitionBuilder {
    private client;
    private id;
    constructor(client: Client, id: string);
    send(): Promise<void>;
}
export declare class DeployDefinitionBuilder {
    private client;
    private id?;
    private name?;
    private bpmnXML?;
    private workflow?;
    constructor(client: Client);
    withID(id: string): this;
    withName(name: string): this;
    withBPMN(xml: string | Blob | Uint8Array | Buffer): this;
    withWorkflow(workflow: Workflow): this;
    send(): Promise<ProcessDefinition>;
}
export declare class InstancesService {
    private client;
    constructor(client: Client);
    list(): ListInstancesBuilder;
    get(id: string): GetInstanceBuilder;
    start(processID: string): StartInstanceBuilder;
    cancel(id: string): CancelInstanceBuilder;
    complete(id: string): CompleteInstanceTaskBuilder;
    resume(id: string): ResumeInstanceBuilder;
    history(id: string): GetInstanceHistoryBuilder;
    incidents(id: string): ListIncidentsBuilder;
    resolveIncident(id: string, incidentID: string): ResolveIncidentBuilder;
    getVisualization(instanceID: string): Promise<VisualizationData>;
    getVisualizationHTML(instanceID: string): Promise<string>;
    subscribe(instanceID: string, onUpdate: () => void): () => void;
}
export declare class ListInstancesBuilder {
    private client;
    constructor(client: Client);
    send(): Promise<ProcessInstance[]>;
}
export declare class GetInstanceBuilder {
    private client;
    private id;
    constructor(client: Client, id: string);
    send(): Promise<ProcessInstance>;
}
export declare class CancelInstanceBuilder {
    private client;
    private id;
    constructor(client: Client, id: string);
    send(): Promise<ProcessInstance>;
}
export declare class StartInstanceBuilder {
    private client;
    private processID;
    private instanceID?;
    private businessKey?;
    private variables;
    constructor(client: Client, processID: string);
    withInstanceID(id: string): this;
    withBusinessKey(key: string): this;
    withVariable(name: string, value: any): this;
    withVariables(variables: Record<string, any>): this;
    send(): Promise<ProcessInstance>;
}
export declare class CompleteInstanceTaskBuilder {
    private client;
    private instanceID;
    private nodeID?;
    private variables;
    constructor(client: Client, instanceID: string);
    withNodeID(nodeID: string): this;
    withVariable(name: string, value: any): this;
    withVariables(variables: Record<string, any>): this;
    send(): Promise<ProcessInstance>;
}
export declare class ResumeInstanceBuilder {
    private client;
    private id;
    constructor(client: Client, id: string);
    send(): Promise<ProcessInstance>;
}
export declare class GetInstanceHistoryBuilder {
    private client;
    private id;
    constructor(client: Client, id: string);
    send(): Promise<HistoryRecord[]>;
}
export declare class ListIncidentsBuilder {
    private client;
    private id;
    constructor(client: Client, id: string);
    send(): Promise<IncidentRecord[]>;
}
export declare class ResolveIncidentBuilder {
    private client;
    private id;
    private incidentID;
    constructor(client: Client, id: string, incidentID: string);
    send(): Promise<void>;
}
export declare class TasksService {
    private client;
    constructor(client: Client);
    list(): ListTasksBuilder;
    get(id: string): GetTaskBuilder;
    claim(id: string): ClaimTaskBuilder;
    complete(id: string): CompleteTaskBuilder;
    resolve(id: string): ResolveTaskBuilder;
}
export declare class ListTasksBuilder {
    private client;
    private assignee?;
    private candidateGroup?;
    private status?;
    constructor(client: Client);
    withAssignee(assignee: string): this;
    withCandidateGroup(candidateGroup: string): this;
    withStatus(status: string): this;
    send(): Promise<TaskRecord[]>;
}
export declare class GetTaskBuilder {
    private client;
    private id;
    constructor(client: Client, id: string);
    send(): Promise<TaskRecord>;
}
export declare class ClaimTaskBuilder {
    private client;
    private id;
    private assignee?;
    constructor(client: Client, id: string);
    withAssignee(assignee: string): this;
    send(): Promise<TaskRecord>;
}
export declare class CompleteTaskBuilder {
    private client;
    private id;
    private variables;
    constructor(client: Client, id: string);
    withVariable(name: string, value: any): this;
    withVariables(variables: Record<string, any>): this;
    send(): Promise<ProcessInstance>;
}
export declare class ResolveTaskBuilder {
    private client;
    private id;
    private variables;
    constructor(client: Client, id: string);
    withVariable(name: string, value: any): this;
    withVariables(variables: Record<string, any>): this;
    send(): Promise<ProcessInstance>;
}
export declare class WebhooksService {
    private client;
    constructor(client: Client);
    list(): ListWebhooksBuilder;
    create(): CreateWebhookBuilder;
    update(id: string): UpdateWebhookBuilder;
    delete(id: string): DeleteWebhookBuilder;
    test(id: string): TestWebhookBuilder;
    deliveries(id: string): ListWebhookDeliveriesBuilder;
}
export declare class ListWebhooksBuilder {
    private client;
    constructor(client: Client);
    send(): Promise<WebhookRecord[]>;
}
export declare class CreateWebhookBuilder {
    private client;
    private url?;
    private secret?;
    private events;
    private processID?;
    private isActive?;
    private enableAudit?;
    constructor(client: Client);
    withURL(url: string): this;
    withSecret(secret: string): this;
    withEvents(events: string[]): this;
    withProcessID(id: string): this;
    withActive(active: boolean): this;
    withAudit(audit: boolean): this;
    send(): Promise<WebhookRecord>;
}
export declare class UpdateWebhookBuilder {
    private client;
    private id;
    private url?;
    private secret?;
    private events;
    private processID?;
    private isActive?;
    private enableAudit?;
    constructor(client: Client, id: string);
    withURL(url: string): this;
    withSecret(secret: string): this;
    withEvents(events: string[]): this;
    withProcessID(id: string): this;
    withActive(active: boolean): this;
    withAudit(audit: boolean): this;
    send(): Promise<WebhookRecord>;
}
export declare class DeleteWebhookBuilder {
    private client;
    private id;
    constructor(client: Client, id: string);
    send(): Promise<void>;
}
export declare class TestWebhookBuilder {
    private client;
    private id;
    constructor(client: Client, id: string);
    send(): Promise<void>;
}
export declare class ListWebhookDeliveriesBuilder {
    private client;
    private id;
    constructor(client: Client, id: string);
    send(): Promise<WebhookDeliveryRecord[]>;
}
export interface TaskContext {
    id: string;
    topic: string;
    definitionId?: string;
    instanceId?: string;
    stepId?: string;
    variables: Record<string, any>;
    timeoutMs?: number;
}
export type TaskHandler = (task: TaskContext) => Promise<Record<string, any> | void> | Record<string, any> | void;
export declare class Worker {
    private serverUrl;
    private apiToken;
    private workerId;
    private maxConcurrency;
    private pollIntervalMs;
    private maxRetries;
    private handlers;
    private idempotencyCache;
    private client;
    private running;
    private loopPromise?;
    constructor(serverUrl: string, apiToken: string);
    withWorkerId(id: string): this;
    withMaxConcurrency(n: number): this;
    withPollInterval(ms: number): this;
    withRetries(n: number): this;
    withTopic(topic: string, handler: TaskHandler): this;
    getClient(): Client;
    getWorkerId(): string;
    getRegisteredTopics(): string[];
    processTask(task: {
        id: string;
        topic?: string;
        activity_id?: string;
        name?: string;
        instance_id?: string;
        variables?: Record<string, any>;
        draft_variables?: Record<string, any>;
    }): Promise<{
        status: "completed" | "incident";
        result?: any;
        error?: string;
    }>;
    pollOnce(): Promise<number>;
    start(): Promise<void>;
    stop(): Promise<void>;
}
