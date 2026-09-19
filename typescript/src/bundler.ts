import * as crypto from 'node:crypto';
import * as fs from 'node:fs';
import * as path from 'node:path';

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
  workflows?: Array<{ id: string; getContentHash: () => string }>;
}

export interface WorkflowSignatureManifest {
  version: string;
  generatedAt: string;
  signatures: Record<string, { definitionId: string; contentHash: string; signature: string }>;
}

/**
 * Computes an HMAC-SHA256 signature for a workflow content hash or definition identifier.
 */
export function computeWorkflowSignature(contentHash: string, secretKey: string): string {
  const cleanHash = contentHash.replace(/^sha256:/, '');
  const hmac = crypto.createHmac('sha256', secretKey);
  hmac.update(cleanHash);
  return hmac.digest('hex');
}

/**
 * Verifies an HMAC-SHA256 signature against a workflow content hash.
 */
export function verifyWorkflowSignature(contentHash: string, signature: string, secretKey: string): boolean {
  const expected = computeWorkflowSignature(contentHash, secretKey);
  if (signature.length !== expected.length) {
    return false;
  }
  return crypto.timingSafeEqual(Buffer.from(signature), Buffer.from(expected));
}

/**
 * Generates a workflow signatures manifest object from a set of workflows.
 */
export function generateWorkflowSignaturesManifest(
  workflows: Array<{ id: string; getContentHash: () => string }>,
  secretKey: string
): WorkflowSignatureManifest {
  const signatures: Record<string, { definitionId: string; contentHash: string; signature: string }> = {};

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
export function nativebpmVitePlugin(options: NativeBPMBundlerOptions = {}) {
  let resolvedOutDir = options.outDir || 'dist';

  return {
    name: 'nativebpm-workflow-signer',

    configResolved(config: any) {
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
        if (typeof (this as any).emitFile === 'function') {
          (this as any).emitFile({
            type: 'asset',
            fileName,
            source: JSON.stringify(manifest, null, 2),
          });
        } else {
          // Direct file system fallback
          try {
            if (!fs.existsSync(resolvedOutDir)) {
              fs.mkdirSync(resolvedOutDir, { recursive: true });
            }
            fs.writeFileSync(path.join(resolvedOutDir, fileName), JSON.stringify(manifest, null, 2), 'utf-8');
          } catch {
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
