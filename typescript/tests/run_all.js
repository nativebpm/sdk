import { spawnSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import { dirname, join } from 'node:path';

const __dirname = dirname(fileURLToPath(import.meta.url));
const testFiles = [
  'case1_linear_zod.test.ts',
  'case2_dmn_decision.test.ts',
  'case3_exclusive_gateway.test.ts',
  'case4_workers.test.ts',
  'case5_boundary_timers.test.ts',
  'case6_bdui_export.test.ts',
  'case7_edge_cases_and_risks.test.ts',
  'case8_zod_openapi_ast.test.ts',
  'case9_polyglot_360_generation.test.ts',
].map(f => join(__dirname, f));

console.log('================================================================');
console.log('🚀 NativeBPM TypeScript SDK: Executing Comprehensive TDD Suite');
console.log('================================================================');

const result = spawnSync('node', ['--test', ...testFiles], {
  stdio: 'inherit',
  env: process.env,
});

if (result.status !== 0) {
  console.error('\n❌ Test suite execution failed!');
  process.exit(result.status || 1);
} else {
  console.log('\n✅ All 7 business-logic & risk patterns passed with 100% PASS!');
}
