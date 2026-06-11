import os
import sys
import gc
import psutil
from nativebpm import Workflow
from nativebpm.builder import DEFAULT_WASM_PATH

def get_memory_usage():
    process = psutil.Process(os.getpid())
    return process.memory_info().rss / (1024 * 1024)  # RSS memory in MB

def run_load_test():
    print("Starting load test & memory profiling for Python SDK...")
    
    # 1. Warm-up
    workflow = Workflow('load-test', 'Load Test Schema')
    workflow.start_event('start').next('end')
    workflow.end_event('end', 'End')
    xml = workflow.build_xml()
    
    # Measure baseline memory after JIT compilation is complete
    gc.collect()
    baseline = get_memory_usage()
    print(f"Baseline Memory (post JIT init): {baseline:.2f} MB")
    
    # 2. Iterate and check memory growth
    iterations = 200
    for i in range(iterations):
        # We recreate the workflow instance to check engine cleanup
        wf = Workflow(f'load-test-{i}', f'Load Test Schema {i}')
        wf.start_event('start').next('end')
        wf.end_event('end', 'End')
        xml = wf.build_xml()
        
        if (i + 1) % 50 == 0:
            gc.collect()
            current = get_memory_usage()
            print(f"Iteration {i+1}/{iterations} - Current Memory: {current:.2f} MB (Delta: {current - baseline:.2f} MB)")
            
    gc.collect()
    final = get_memory_usage()
    print(f"Final Memory: {final:.2f} MB (Delta: {final - baseline:.2f} MB)")
    
    # Acceptable memory delta should be minimal (typically < 5MB overhead for runtime initialization/cache)
    if final - baseline > 15.0:
        print("WARNING: Possible memory leak detected in Python SDK!")
        sys.exit(1)
    else:
        print("SUCCESS: Memory profile is stable, no leaks detected.")

if __name__ == '__main__':
    run_load_test()
