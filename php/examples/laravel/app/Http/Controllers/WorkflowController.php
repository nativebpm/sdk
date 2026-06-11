<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use NativeBPM\Client\Api\DefaultApi;
use NativeBPM\Client\Model\StartInstanceRequest;
use NativeBPM\Client\ApiException;
use Illuminate\Http\JsonResponse;

class WorkflowController extends Controller
{
    protected DefaultApi $api;

    /**
     * WorkflowController constructor.
     *
     * @param DefaultApi $api Inject the configured NativeBPM API client singleton.
     */
    public function __construct(DefaultApi $api)
    {
        $this->api = $api;
    }

    /**
     * List all deployed definitions from the engine.
     *
     * GET /api/workflows
     *
     * @return JsonResponse
     */
    public function index(): JsonResponse
    {
        try {
            $definitions = $this->api->listDefinitions();
            
            $data = array_map(function ($def) {
                return [
                    'id'      => $def->getId(),
                    'name'    => $def->getName(),
                    'version' => $def->getVersion(),
                    'hash'    => $def->getHash(),
                ];
            }, $definitions);

            return response()->json([
                'success' => true,
                'data'    => $data,
            ]);
        } catch (ApiException $e) {
            return response()->json([
                'success' => false,
                'error'   => 'Failed to list workflow definitions: ' . $e->getMessage(),
                'code'    => $e->getCode(),
                'body'    => json_decode($e->getResponseBody(), true) ?? $e->getResponseBody(),
            ], 500);
        }
    }

    /**
     * Deploy a new BPMN 2.0 workflow XML file to the engine.
     *
     * POST /api/workflows/deploy
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function deploy(Request $request): JsonResponse
    {
        $request->validate([
            'bpmn_file' => 'required|file',
        ]);

        try {
            $file = $request->file('bpmn_file');
            
            // Pass the local path to deployDefinition (handles parsing and multipart stream creation)
            $definition = $this->api->deployDefinition($file->getPathname());

            return response()->json([
                'success' => true,
                'data'    => [
                    'id'      => $definition->getId(),
                    'name'    => $definition->getName(),
                    'version' => $definition->getVersion(),
                    'hash'    => $definition->getHash(),
                ],
            ]);
        } catch (ApiException $e) {
            return response()->json([
                'success' => false,
                'error'   => 'BPMN deployment failed: ' . $e->getMessage(),
                'code'    => $e->getCode(),
                'body'    => json_decode($e->getResponseBody(), true) ?? $e->getResponseBody(),
            ], 500);
        }
    }

    /**
     * Start a new workflow process instance for a specific process definition ID.
     *
     * POST /api/workflows/{definitionId}/start
     *
     * @param Request $request
     * @param string $definitionId
     * @return JsonResponse
     */
    public function start(Request $request, string $definitionId): JsonResponse
    {
        try {
            // Generate UUID process instance ID & Business Key
            $instanceId = \Illuminate\Support\Str::uuid()->toString();
            $businessKey = 'bk-php-' . time();

            // Populate StartInstanceRequest object
            $startRequest = new StartInstanceRequest();
            $startRequest->setInstanceId($instanceId);
            $startRequest->setBusinessKey($businessKey);
            $startRequest->setVariables([
                'started_by'   => 'laravel-example-client',
                'laravel_env'  => app()->environment(),
                'custom_value' => $request->input('custom_value', 'default'),
            ]);

            // Call the SDK startInstance method
            $instance = $this->api->startInstance($definitionId, $startRequest);

            return response()->json([
                'success' => true,
                'data'    => [
                    'id'           => $instance->getId(),
                    'process_id'   => $instance->getProcessId(),
                    'business_key' => $instance->getBusinessKey(),
                    'completed'    => $instance->getCompleted(),
                ],
            ]);
        } catch (ApiException $e) {
            return response()->json([
                'success' => false,
                'error'   => 'Failed to start process instance: ' . $e->getMessage(),
                'code'    => $e->getCode(),
                'body'    => json_decode($e->getResponseBody(), true) ?? $e->getResponseBody(),
            ], 500);
        }
    }
}
