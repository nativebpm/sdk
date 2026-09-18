import { z, toJSONSchema } from 'zod';
export { z, toJSONSchema };
export declare const InVariableSchema: z.ZodObject<{
    source: z.ZodOptional<z.ZodString>;
    target: z.ZodOptional<z.ZodString>;
    variables: z.ZodOptional<z.ZodString>;
    local: z.ZodOptional<z.ZodBoolean>;
}, z.core.$strip>;
export declare const OutVariableSchema: z.ZodObject<{
    source: z.ZodOptional<z.ZodString>;
    target: z.ZodOptional<z.ZodString>;
    variables: z.ZodOptional<z.ZodString>;
}, z.core.$strip>;
export declare const DMNInputSchema: z.ZodObject<{
    expression: z.ZodString;
    type: z.ZodString;
}, z.core.$strip>;
export declare const DMNOutputSchema: z.ZodObject<{
    name: z.ZodString;
    type: z.ZodString;
}, z.core.$strip>;
export declare const DMNRuleSchema: z.ZodObject<{
    inputs: z.ZodArray<z.ZodString>;
    outputs: z.ZodArray<z.ZodString>;
}, z.core.$strip>;
export declare const NodeASTSchema: z.ZodObject<{
    id: z.ZodString;
    name: z.ZodString;
    type: z.ZodEnum<{
        startEvent: "startEvent";
        endEvent: "endEvent";
        serviceTask: "serviceTask";
        userTask: "userTask";
        aiServiceTask: "aiServiceTask";
        aiTask: "aiTask";
        exclusiveGateway: "exclusiveGateway";
        parallelGateway: "parallelGateway";
        eventBasedGateway: "eventBasedGateway";
        callActivity: "callActivity";
        businessRuleTask: "businessRuleTask";
        boundaryTimerEvent: "boundaryTimerEvent";
    }>;
    topic: z.ZodOptional<z.ZodString>;
    wasmPath: z.ZodOptional<z.ZodString>;
    provider: z.ZodOptional<z.ZodString>;
    model: z.ZodOptional<z.ZodString>;
    prompt: z.ZodOptional<z.ZodString>;
    systemInstruction: z.ZodOptional<z.ZodString>;
    responseSchema: z.ZodOptional<z.ZodUnion<readonly [z.ZodString, z.ZodRecord<z.ZodString, z.ZodAny>]>>;
    temperature: z.ZodOptional<z.ZodNumber>;
    resultVar: z.ZodOptional<z.ZodString>;
    assignee: z.ZodOptional<z.ZodString>;
    candidateGroups: z.ZodOptional<z.ZodString>;
    dueDate: z.ZodOptional<z.ZodString>;
    inputSchema: z.ZodOptional<z.ZodUnion<readonly [z.ZodString, z.ZodRecord<z.ZodString, z.ZodAny>]>>;
    formId: z.ZodOptional<z.ZodString>;
    formKey: z.ZodOptional<z.ZodString>;
    calledElement: z.ZodOptional<z.ZodString>;
    inVariables: z.ZodOptional<z.ZodArray<z.ZodObject<{
        source: z.ZodOptional<z.ZodString>;
        target: z.ZodOptional<z.ZodString>;
        variables: z.ZodOptional<z.ZodString>;
        local: z.ZodOptional<z.ZodBoolean>;
    }, z.core.$strip>>>;
    outVariables: z.ZodOptional<z.ZodArray<z.ZodObject<{
        source: z.ZodOptional<z.ZodString>;
        target: z.ZodOptional<z.ZodString>;
        variables: z.ZodOptional<z.ZodString>;
    }, z.core.$strip>>>;
    decisionRef: z.ZodOptional<z.ZodString>;
    mapDecisionResult: z.ZodOptional<z.ZodString>;
    hitPolicy: z.ZodOptional<z.ZodString>;
    inputs: z.ZodOptional<z.ZodArray<z.ZodObject<{
        expression: z.ZodString;
        type: z.ZodString;
    }, z.core.$strip>>>;
    outputs: z.ZodOptional<z.ZodArray<z.ZodObject<{
        name: z.ZodString;
        type: z.ZodString;
    }, z.core.$strip>>>;
    rules: z.ZodOptional<z.ZodArray<z.ZodObject<{
        inputs: z.ZodArray<z.ZodString>;
        outputs: z.ZodArray<z.ZodString>;
    }, z.core.$strip>>>;
    attachedToRef: z.ZodOptional<z.ZodString>;
    timeDuration: z.ZodOptional<z.ZodString>;
    timeDate: z.ZodOptional<z.ZodString>;
    timeCycle: z.ZodOptional<z.ZodString>;
    cancelActivity: z.ZodOptional<z.ZodBoolean>;
}, z.core.$strip>;
export declare const FlowASTSchema: z.ZodObject<{
    id: z.ZodString;
    source: z.ZodString;
    target: z.ZodString;
    condition: z.ZodOptional<z.ZodString>;
}, z.core.$strip>;
export declare const WorkflowASTSchema: z.ZodObject<{
    id: z.ZodString;
    name: z.ZodString;
    inputSchema: z.ZodOptional<z.ZodUnion<readonly [z.ZodString, z.ZodRecord<z.ZodString, z.ZodAny>]>>;
    nodes: z.ZodArray<z.ZodObject<{
        id: z.ZodString;
        name: z.ZodString;
        type: z.ZodEnum<{
            startEvent: "startEvent";
            endEvent: "endEvent";
            serviceTask: "serviceTask";
            userTask: "userTask";
            aiServiceTask: "aiServiceTask";
            aiTask: "aiTask";
            exclusiveGateway: "exclusiveGateway";
            parallelGateway: "parallelGateway";
            eventBasedGateway: "eventBasedGateway";
            callActivity: "callActivity";
            businessRuleTask: "businessRuleTask";
            boundaryTimerEvent: "boundaryTimerEvent";
        }>;
        topic: z.ZodOptional<z.ZodString>;
        wasmPath: z.ZodOptional<z.ZodString>;
        provider: z.ZodOptional<z.ZodString>;
        model: z.ZodOptional<z.ZodString>;
        prompt: z.ZodOptional<z.ZodString>;
        systemInstruction: z.ZodOptional<z.ZodString>;
        responseSchema: z.ZodOptional<z.ZodUnion<readonly [z.ZodString, z.ZodRecord<z.ZodString, z.ZodAny>]>>;
        temperature: z.ZodOptional<z.ZodNumber>;
        resultVar: z.ZodOptional<z.ZodString>;
        assignee: z.ZodOptional<z.ZodString>;
        candidateGroups: z.ZodOptional<z.ZodString>;
        dueDate: z.ZodOptional<z.ZodString>;
        inputSchema: z.ZodOptional<z.ZodUnion<readonly [z.ZodString, z.ZodRecord<z.ZodString, z.ZodAny>]>>;
        formId: z.ZodOptional<z.ZodString>;
        formKey: z.ZodOptional<z.ZodString>;
        calledElement: z.ZodOptional<z.ZodString>;
        inVariables: z.ZodOptional<z.ZodArray<z.ZodObject<{
            source: z.ZodOptional<z.ZodString>;
            target: z.ZodOptional<z.ZodString>;
            variables: z.ZodOptional<z.ZodString>;
            local: z.ZodOptional<z.ZodBoolean>;
        }, z.core.$strip>>>;
        outVariables: z.ZodOptional<z.ZodArray<z.ZodObject<{
            source: z.ZodOptional<z.ZodString>;
            target: z.ZodOptional<z.ZodString>;
            variables: z.ZodOptional<z.ZodString>;
        }, z.core.$strip>>>;
        decisionRef: z.ZodOptional<z.ZodString>;
        mapDecisionResult: z.ZodOptional<z.ZodString>;
        hitPolicy: z.ZodOptional<z.ZodString>;
        inputs: z.ZodOptional<z.ZodArray<z.ZodObject<{
            expression: z.ZodString;
            type: z.ZodString;
        }, z.core.$strip>>>;
        outputs: z.ZodOptional<z.ZodArray<z.ZodObject<{
            name: z.ZodString;
            type: z.ZodString;
        }, z.core.$strip>>>;
        rules: z.ZodOptional<z.ZodArray<z.ZodObject<{
            inputs: z.ZodArray<z.ZodString>;
            outputs: z.ZodArray<z.ZodString>;
        }, z.core.$strip>>>;
        attachedToRef: z.ZodOptional<z.ZodString>;
        timeDuration: z.ZodOptional<z.ZodString>;
        timeDate: z.ZodOptional<z.ZodString>;
        timeCycle: z.ZodOptional<z.ZodString>;
        cancelActivity: z.ZodOptional<z.ZodBoolean>;
    }, z.core.$strip>>;
    flows: z.ZodArray<z.ZodObject<{
        id: z.ZodString;
        source: z.ZodString;
        target: z.ZodString;
        condition: z.ZodOptional<z.ZodString>;
    }, z.core.$strip>>;
}, z.core.$strip>;
export type InVariableAST = z.infer<typeof InVariableSchema>;
export type OutVariableAST = z.infer<typeof OutVariableSchema>;
export type DMNInputAST = z.infer<typeof DMNInputSchema>;
export type DMNOutputAST = z.infer<typeof DMNOutputSchema>;
export type DMNRuleAST = z.infer<typeof DMNRuleSchema>;
export type NodeAST = z.infer<typeof NodeASTSchema>;
export type FlowAST = z.infer<typeof FlowASTSchema>;
export type WorkflowAST = z.infer<typeof WorkflowASTSchema>;
/**
 * Export the WorkflowAST schema as an OpenAPI 3.0 Schema Object.
 */
export declare function exportWorkflowOpenAPISchema(): Record<string, any>;
/**
 * Export the complete standalone OpenAPI 3.0.3 specification document for Workflow AST.
 * This document is fully compatible with OpenAPI 3.0 parsers and Swagger UI previews (e.g. GitLab/GitHub).
 */
export declare function exportWorkflowOpenAPIDocument(): Record<string, any>;
/**
 * Export the WorkflowAST schema as a standard JSON Schema Draft 2020-12 document.
 */
export declare function exportWorkflowJSONSchema(): Record<string, any>;
