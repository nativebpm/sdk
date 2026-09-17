

# NodeAST


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**id** | **String** |  |  |
|**name** | **String** |  |  |
|**type** | [**TypeEnum**](#TypeEnum) |  |  |
|**topic** | **String** |  |  [optional] |
|**wasmPath** | **String** |  |  [optional] |
|**provider** | **String** |  |  [optional] |
|**model** | **String** |  |  [optional] |
|**prompt** | **String** |  |  [optional] |
|**systemInstruction** | **String** |  |  [optional] |
|**responseSchema** | **String** |  |  [optional] |
|**temperature** | **BigDecimal** |  |  [optional] |
|**resultVar** | **String** |  |  [optional] |
|**assignee** | **String** |  |  [optional] |
|**candidateGroups** | **String** |  |  [optional] |
|**dueDate** | **String** |  |  [optional] |
|**inputSchema** | **String** |  |  [optional] |
|**formId** | **String** |  |  [optional] |
|**formKey** | **String** |  |  [optional] |
|**calledElement** | **String** |  |  [optional] |
|**inVariables** | [**List&lt;InVariable&gt;**](InVariable.md) |  |  [optional] |
|**outVariables** | [**List&lt;OutVariable&gt;**](OutVariable.md) |  |  [optional] |
|**decisionRef** | **String** |  |  [optional] |
|**mapDecisionResult** | **String** |  |  [optional] |
|**hitPolicy** | **String** |  |  [optional] |
|**inputs** | [**List&lt;DMNInputAST&gt;**](DMNInputAST.md) |  |  [optional] |
|**outputs** | [**List&lt;DMNOutputAST&gt;**](DMNOutputAST.md) |  |  [optional] |
|**rules** | [**List&lt;DMNRuleAST&gt;**](DMNRuleAST.md) |  |  [optional] |
|**attachedToRef** | **String** |  |  [optional] |
|**timeDuration** | **String** |  |  [optional] |
|**timeDate** | **String** |  |  [optional] |
|**timeCycle** | **String** |  |  [optional] |
|**cancelActivity** | **Boolean** |  |  [optional] |



## Enum: TypeEnum

| Name | Value |
|---- | -----|
| START_EVENT | &quot;startEvent&quot; |
| END_EVENT | &quot;endEvent&quot; |
| SERVICE_TASK | &quot;serviceTask&quot; |
| USER_TASK | &quot;userTask&quot; |
| AI_SERVICE_TASK | &quot;aiServiceTask&quot; |
| AI_TASK | &quot;aiTask&quot; |
| EXCLUSIVE_GATEWAY | &quot;exclusiveGateway&quot; |
| PARALLEL_GATEWAY | &quot;parallelGateway&quot; |
| EVENT_BASED_GATEWAY | &quot;eventBasedGateway&quot; |
| CALL_ACTIVITY | &quot;callActivity&quot; |
| BUSINESS_RULE_TASK | &quot;businessRuleTask&quot; |
| BOUNDARY_TIMER_EVENT | &quot;boundaryTimerEvent&quot; |



