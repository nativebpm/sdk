import { z, toJSONSchema } from 'zod';

export { z, toJSONSchema };

export const InVariableSchema = z.object({
  source: z.string().optional().describe('Source variable name'),
  target: z.string().optional().describe('Target variable name in child context'),
  variables: z.string().optional().describe('Special mapping, e.g. "all"'),
  local: z.boolean().optional().describe('Whether variable is scoped locally'),
});

export const OutVariableSchema = z.object({
  source: z.string().optional().describe('Source variable name in child context'),
  target: z.string().optional().describe('Target variable name in parent context'),
  variables: z.string().optional().describe('Special mapping, e.g. "all"'),
});

export const DMNInputSchema = z.object({
  expression: z.string().describe('Input expression (variable name)'),
  type: z.string().describe('Data type (string, number, boolean)'),
});

export const DMNOutputSchema = z.object({
  name: z.string().describe('Output variable name'),
  type: z.string().describe('Data type (string, number, boolean)'),
});

export const DMNRuleSchema = z.object({
  inputs: z.array(z.string()).describe('Input match conditions per column'),
  outputs: z.array(z.string()).describe('Output assigned values per column'),
});

export const NodeASTSchema = z.object({
  id: z.string().min(1).describe('Unique element identifier'),
  name: z.string().describe('Human-readable element label'),
  type: z.enum([
    'startEvent',
    'endEvent',
    'serviceTask',
    'userTask',
    'aiServiceTask',
    'aiTask',
    'exclusiveGateway',
    'parallelGateway',
    'eventBasedGateway',
    'callActivity',
    'businessRuleTask',
    'boundaryTimerEvent',
  ]).describe('BPMN 2.0 node type'),
  topic: z.string().optional().describe('Worker subscription topic for service tasks'),
  wasmPath: z.string().optional().describe('Local or OCI path to WASM module for Wasmee worker'),
  provider: z.string().optional().describe('AI Provider (openai, gemini, anthropic)'),
  model: z.string().optional().describe('Model identifier (gpt-4o, gemini-1.5-pro)'),
  prompt: z.string().optional().describe('Prompt template with variable interpolation'),
  systemInstruction: z.string().optional().describe('System instruction for AI assistant'),
  responseSchema: z.union([z.string(), z.record(z.string(), z.any())]).optional().describe('Expected JSON Schema for AI response'),
  temperature: z.number().optional().describe('Model temperature (0.0 to 2.0)'),
  resultVar: z.string().optional().describe('Variable name to store task result'),
  assignee: z.string().optional().describe('Assigned user identifier for UserTask'),
  candidateGroups: z.string().optional().describe('Authorized candidate groups for UserTask'),
  dueDate: z.string().optional().describe('ISO 8601 due date for task completion'),
  inputSchema: z.union([z.string(), z.record(z.string(), z.any())]).optional().describe('JSON Schema defining the UserTask form fields'),
  formId: z.string().optional().describe('Unique form definition identifier'),
  formKey: z.string().optional().describe('Camunda-compatible form key reference'),
  calledElement: z.string().optional().describe('Called subprocess definition ID'),
  inVariables: z.array(InVariableSchema).optional().describe('Input variable mappings for callActivity'),
  outVariables: z.array(OutVariableSchema).optional().describe('Output variable mappings for callActivity'),
  decisionRef: z.string().optional().describe('DMN decision table reference ID'),
  mapDecisionResult: z.string().optional().describe('Result mapping policy (singleEntry, collectEntries)'),
  hitPolicy: z.string().optional().describe('DMN hit policy (UNIQUE, FIRST, COLLECT)'),
  inputs: z.array(DMNInputSchema).optional().describe('DMN input column declarations'),
  outputs: z.array(DMNOutputSchema).optional().describe('DMN output column declarations'),
  rules: z.array(DMNRuleSchema).optional().describe('DMN decision matrix rows'),
  attachedToRef: z.string().optional().describe('Host task ID for boundaryTimerEvent'),
  timeDuration: z.string().optional().describe('ISO 8601 duration (e.g. PT15M, P1D)'),
  timeDate: z.string().optional().describe('ISO 8601 target timestamp'),
  timeCycle: z.string().optional().describe('ISO 8601 recurrence cycle (e.g. R3/PT10M)'),
  cancelActivity: z.boolean().optional().describe('Whether boundary timer interrupts the host activity'),
});

export const FlowASTSchema = z.object({
  id: z.string().min(1).describe('Unique sequence flow identifier'),
  source: z.string().min(1).describe('Source node ID'),
  target: z.string().min(1).describe('Target node ID'),
  condition: z.string().optional().describe('JUEL/FEEL condition expression (e.g. ${approved == true})'),
});

export const WorkflowASTSchema = z.object({
  id: z.string().min(1).describe('Unique process definition identifier'),
  name: z.string().min(1).describe('Human-readable process name'),
  inputSchema: z.union([z.string(), z.record(z.string(), z.any())]).optional().describe('JSON Schema contract for process start variables'),
  nodes: z.array(NodeASTSchema).describe('List of BPMN nodes'),
  flows: z.array(FlowASTSchema).describe('List of sequence flows connecting nodes'),
});

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
export function exportWorkflowOpenAPISchema(): Record<string, any> {
  return toJSONSchema(WorkflowASTSchema, { target: 'openapi-3.0' }) as Record<string, any>;
}

/**
 * Export the complete standalone OpenAPI 3.0.3 specification document for Workflow AST.
 * This document is fully compatible with OpenAPI 3.0 parsers and Swagger UI previews (e.g. GitLab/GitHub).
 */
export function exportWorkflowOpenAPIDocument(): Record<string, any> {
  return {
    openapi: '3.0.3',
    info: {
      title: 'NativeBPM Workflow AST Specification',
      version: '1.0.0',
      description: 'Declarative Zod 4-compiled OpenAPI 3.0 specification for NativeBPM Workflow AST graph models, nodes, flows, and DMN decision tables.',
    },
    paths: {
      '/api/deploy': {
        post: {
          summary: 'Deploy Workflow AST process definition',
          description: 'Deploy a programmatic Workflow AST process definition directly to NativeBPM engine.',
          operationId: 'deployWorkflowAST',
          requestBody: {
            required: true,
            content: {
              'application/json': {
                schema: {
                  $ref: '#/components/schemas/WorkflowAST',
                },
              },
            },
          },
          responses: {
            '200': {
              description: 'Process definition deployed successfully',
            },
          },
        },
      },
    },
    components: {
      schemas: {
        WorkflowAST: exportWorkflowOpenAPISchema(),
        NodeAST: toJSONSchema(NodeASTSchema, { target: 'openapi-3.0' }) as Record<string, any>,
        FlowAST: toJSONSchema(FlowASTSchema, { target: 'openapi-3.0' }) as Record<string, any>,
        DMNRule: toJSONSchema(DMNRuleSchema, { target: 'openapi-3.0' }) as Record<string, any>,
        DMNInput: toJSONSchema(DMNInputSchema, { target: 'openapi-3.0' }) as Record<string, any>,
        DMNOutput: toJSONSchema(DMNOutputSchema, { target: 'openapi-3.0' }) as Record<string, any>,
        InVariable: toJSONSchema(InVariableSchema, { target: 'openapi-3.0' }) as Record<string, any>,
        OutVariable: toJSONSchema(OutVariableSchema, { target: 'openapi-3.0' }) as Record<string, any>,
      },
    },
  };
}

/**
 * Export the WorkflowAST schema as a standard JSON Schema Draft 2020-12 document.
 */
export function exportWorkflowJSONSchema(): Record<string, any> {
  return toJSONSchema(WorkflowASTSchema, { target: 'draft-2020-12' }) as Record<string, any>;
}

