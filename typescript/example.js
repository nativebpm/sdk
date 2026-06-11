import { run as runWasm } from './examples/workflow_wasm.js';
import { run as runNative } from './examples/workflow_native.js';

async function main() {
  console.log("==================================================");
  console.log("🚀 RUNNING NATIVEBPM TS SDK EXAMPLES");
  console.log("==================================================");
  
  // 1. Run WASM-based Workflow example
  await runWasm();
  
  console.log("\n--------------------------------------------------\n");
  
  // 2. Run Native-based Workflow example
  await runNative();
  
  console.log("==================================================");
  console.log("🎉 ALL EXAMPLES FINISHED");
  console.log("==================================================");
}

main().catch(err => {
  console.error("Fatal error running examples:", err);
  process.exit(1);
});
