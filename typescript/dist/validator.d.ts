export interface UIWidgetSpec {
    name: string;
    label: string;
    type: string;
    widget: string;
    required: boolean;
    value: any;
    options: string[];
    minLength?: number;
    maxLength?: number;
    minimum?: number;
    maximum?: number;
    pattern?: string;
    format?: string;
    xStep?: number;
}
export interface ValidationResult {
    valid: boolean;
    errors: string[];
}
export declare function validateJSON(schema: string | object, data: string | object, wasmPath?: string): Promise<ValidationResult>;
export declare function parseSchema(schema: string | object, variables?: object, wasmPath?: string): Promise<UIWidgetSpec[]>;
