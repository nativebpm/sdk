export interface NativeBPMBundlerOptions {
    /**
     * Secret key used to sign workflow content hashes during build.
     * If not provided, will read from process.env.NATIVEBPM_SIGNING_SECRET.
     */
    secretKey?: string;
    /**
     * Output directory where the signatures manifest should be saved.
     * Defaults to 'dist' or the bundler's outDir.
     */
    outDir?: string;
    /**
     * Filename for the signatures manifest.
     * Defaults to '__workflow_signatures.json'.
     */
    manifestFile?: string;
    /**
     * Known workflows to pre-sign at build time.
     */
    workflows?: Array<{
        id: string;
        getContentHash: () => string;
    }>;
}
export interface WorkflowSignatureManifest {
    version: string;
    generatedAt: string;
    signatures: Record<string, {
        definitionId: string;
        contentHash: string;
        signature: string;
    }>;
}
/**
 * Computes an HMAC-SHA256 signature for a workflow content hash or definition identifier.
 */
export declare function computeWorkflowSignature(contentHash: string, secretKey: string): string;
/**
 * Verifies an HMAC-SHA256 signature against a workflow content hash.
 */
export declare function verifyWorkflowSignature(contentHash: string, signature: string, secretKey: string): boolean;
/**
 * Generates a workflow signatures manifest object from a set of workflows.
 */
export declare function generateWorkflowSignaturesManifest(workflows: Array<{
    id: string;
    getContentHash: () => string;
}>, secretKey: string): WorkflowSignatureManifest;
/**
 * NativeBPM Vite / Rollup bundler plugin.
 * Automatically computes workflow signatures at build time and writes __workflow_signatures.json.
 */
export declare function nativebpmVitePlugin(options?: NativeBPMBundlerOptions): {
    name: string;
    configResolved(config: any): void;
    generateBundle(): void;
};
/**
 * Rollup plugin alias for nativebpmVitePlugin.
 */
export declare const nativebpmRollupPlugin: typeof nativebpmVitePlugin;
