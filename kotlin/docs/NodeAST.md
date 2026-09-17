
# NodeAST

## Properties
| Name | Type | Description | Notes |
| ------------ | ------------- | ------------- | ------------- |
| **id** | **kotlin.String** |  |  |
| **name** | **kotlin.String** |  |  |
| **type** | [**inline**](#Type) |  |  |
| **topic** | **kotlin.String** |  |  [optional] |
| **wasmPath** | **kotlin.String** |  |  [optional] |
| **provider** | **kotlin.String** |  |  [optional] |
| **model** | **kotlin.String** |  |  [optional] |
| **prompt** | **kotlin.String** |  |  [optional] |
| **systemInstruction** | **kotlin.String** |  |  [optional] |
| **responseSchema** | **kotlin.String** |  |  [optional] |
| **temperature** | [**java.math.BigDecimal**](java.math.BigDecimal.md) |  |  [optional] |
| **resultVar** | **kotlin.String** |  |  [optional] |
| **assignee** | **kotlin.String** |  |  [optional] |
| **candidateGroups** | **kotlin.String** |  |  [optional] |
| **dueDate** | **kotlin.String** |  |  [optional] |
| **inputSchema** | **kotlin.String** |  |  [optional] |
| **formId** | **kotlin.String** |  |  [optional] |
| **formKey** | **kotlin.String** |  |  [optional] |
| **calledElement** | **kotlin.String** |  |  [optional] |
| **inVariables** | [**kotlin.collections.List&lt;InVariable&gt;**](InVariable.md) |  |  [optional] |
| **outVariables** | [**kotlin.collections.List&lt;OutVariable&gt;**](OutVariable.md) |  |  [optional] |
| **decisionRef** | **kotlin.String** |  |  [optional] |
| **mapDecisionResult** | **kotlin.String** |  |  [optional] |
| **hitPolicy** | **kotlin.String** |  |  [optional] |
| **inputs** | [**kotlin.collections.List&lt;DMNInputAST&gt;**](DMNInputAST.md) |  |  [optional] |
| **outputs** | [**kotlin.collections.List&lt;DMNOutputAST&gt;**](DMNOutputAST.md) |  |  [optional] |
| **rules** | [**kotlin.collections.List&lt;DMNRuleAST&gt;**](DMNRuleAST.md) |  |  [optional] |
| **attachedToRef** | **kotlin.String** |  |  [optional] |
| **timeDuration** | **kotlin.String** |  |  [optional] |
| **timeDate** | **kotlin.String** |  |  [optional] |
| **timeCycle** | **kotlin.String** |  |  [optional] |
| **cancelActivity** | **kotlin.Boolean** |  |  [optional] |


<a id="Type"></a>
## Enum: type
| Name | Value |
| ---- | ----- |
| type | startEvent, endEvent, serviceTask, userTask, aiServiceTask, aiTask, exclusiveGateway, parallelGateway, eventBasedGateway, callActivity, businessRuleTask, boundaryTimerEvent |



