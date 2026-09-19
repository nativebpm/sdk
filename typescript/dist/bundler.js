import * as crypto from 'node:crypto';
import * as fs from 'node:fs';
import * as path from 'node:path';
/**
 * Computes an HMAC-SHA256 signature for a workflow content hash or definition identifier.
 */
export function computeWorkflowSignature(contentHash, secretKey) {
    const cleanHash = contentHash.replace(/^sha256:/, '');
    const hmac = crypto.createHmac('sha256', secretKey);
    hmac.update(cleanHash);
    return hmac.digest('hex');
}
/**
 * Verifies an HMAC-SHA256 signature against a workflow content hash.
 */
export function verifyWorkflowSignature(contentHash, signature, secretKey) {
    const expected = computeWorkflowSignature(contentHash, secretKey);
    if (signature.length !== expected.length) {
        return false;
    }
    return crypto.timingSafeEqual(Buffer.from(signature), Buffer.from(expected));
}
/**
 * Generates a workflow signatures manifest object from a set of workflows.
 */
export function generateWorkflowSignaturesManifest(workflows, secretKey) {
    const signatures = {};
    for (const wf of workflows) {
        const hash = wf.getContentHash();
        signatures[wf.id] = {
            definitionId: wf.id,
            contentHash: hash,
            signature: computeWorkflowSignature(hash, secretKey),
        };
    }
    return {
        version: '1.0.0',
        generatedAt: new Date().toISOString(),
        signatures,
    };
}
/**
 * NativeBPM Vite / Rollup bundler plugin.
 * Automatically computes workflow signatures at build time and writes __workflow_signatures.json.
 */
export function nativebpmVitePlugin(options = {}) {
    let resolvedOutDir = options.outDir || 'dist';
    return {
        name: 'nativebpm-workflow-signer',
        configResolved(config) {
            if (!options.outDir && config.build && config.build.outDir) {
                resolvedOutDir = config.build.outDir;
            }
        },
        generateBundle() {
            const secret = options.secretKey || process.env.NATIVEBPM_SIGNING_SECRET;
            if (!secret) {
                return;
            }
            if (options.workflows && options.workflows.length > 0) {
                const manifest = generateWorkflowSignaturesManifest(options.workflows, secret);
                const fileName = options.manifestFile || '__workflow_signatures.json';
                // Emit asset if bundler context supports emitFile
                if (typeof this.emitFile === 'function') {
                    this.emitFile({
                        type: 'asset',
                        fileName,
                        source: JSON.stringify(manifest, null, 2),
                    });
                }
                else {
                    // Direct file system fallback
                    try {
                        if (!fs.existsSync(resolvedOutDir)) {
                            fs.mkdirSync(resolvedOutDir, { recursive: true });
                        }
                        fs.writeFileSync(path.join(resolvedOutDir, fileName), JSON.stringify(manifest, null, 2), 'utf-8');
                    }
                    catch {
                        // Ignore if build directory is handled entirely in memory
                    }
                }
            }
        },
    };
}
/**
 * Rollup plugin alias for nativebpmVitePlugin.
 */
export const nativebpmRollupPlugin = nativebpmVitePlugin;
