using System;
using System.IO;
using System.Linq;
using System.Reflection;
using Wasmtime;
using System.Text.Json;
using System.Text;

class Program
{
    static void Main()
    {
        try
        {
            var engine = new Engine();
            var moduleBytes = File.ReadAllBytes("src/NativeBPM.Client/core.wasm");
            if (!(moduleBytes.Length >= 4 && moduleBytes[0] == 0x00 && moduleBytes[1] == 0x61 && moduleBytes[2] == 0x73 && moduleBytes[3] == 0x6d))
            {
                using var ms = new MemoryStream(moduleBytes);
                using var brotli = new System.IO.Compression.BrotliStream(ms, System.IO.Compression.CompressionMode.Decompress);
                using var outMs = new MemoryStream();
                brotli.CopyTo(outMs);
                moduleBytes = outMs.ToArray();
            }

            var module = Wasmtime.Module.FromBytes(engine, "core", moduleBytes);
            using var store = new Store(engine);
            var wasiConfig = new WasiConfiguration();
            store.SetWasiConfiguration(wasiConfig);

            using var linker = new Linker(engine);
            
            // Define proc_exit first
            linker.Define("wasi_snapshot_preview1", "proc_exit", Function.FromCallback(store, (Caller caller, int exitCode) => {
                throw new Exception($"exit_{exitCode}");
            }));

            // Call DefineWasi after
            linker.DefineWasi();

            var instance = linker.Instantiate(store, module);
            var start = instance.GetFunction("_start");
            if (start != null)
            {
                Console.WriteLine("Calling _start...");
                try
                {
                    start.Invoke();
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Caught expected exit exception: {ex.Message}");
                }
            }

            var allocate = instance.GetFunction("allocate");
            var compileWorkflow = instance.GetFunction("compileWorkflow");
            var memory = instance.GetMemory("memory");

            if (allocate == null || compileWorkflow == null || memory == null)
            {
                Console.WriteLine("Exports are null!");
                return;
            }

            Console.WriteLine("Calling allocate...");
            string json = "{\"id\":\"test-process\",\"name\":\"Test Process Schema\",\"nodes\":[],\"flows\":[]}";
            byte[] bytes = Encoding.UTF8.GetBytes(json);
            int ptr = (int)allocate.Invoke(bytes.Length);
            Console.WriteLine($"Allocated at: {ptr}");

            memory.WriteString(ptr, json, Encoding.UTF8);

            Console.WriteLine("Calling compileWorkflow...");
            long resultPacked = (long)compileWorkflow.Invoke(ptr, bytes.Length);
            int resPtr = (int)(resultPacked >> 32);
            int resSize = (int)(resultPacked & 0xFFFFFFFF);

            string resultJson = memory.ReadString(resPtr, resSize, Encoding.UTF8);
            Console.WriteLine($"Result JSON: {resultJson}");
        }
        catch (Exception ex)
        {
            Console.WriteLine(ex.ToString());
        }
    }
}
