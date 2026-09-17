# nativebpm_client.model.NodeAST

## Load the model package
```dart
import 'package:nativebpm_client/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** |  | 
**name** | **String** |  | 
**type** | **String** |  | 
**topic** | **String** |  | [optional] 
**wasmPath** | **String** |  | [optional] 
**provider** | **String** |  | [optional] 
**model** | **String** |  | [optional] 
**prompt** | **String** |  | [optional] 
**systemInstruction** | **String** |  | [optional] 
**responseSchema** | **String** |  | [optional] 
**temperature** | **num** |  | [optional] 
**resultVar** | **String** |  | [optional] 
**assignee** | **String** |  | [optional] 
**candidateGroups** | **String** |  | [optional] 
**dueDate** | **String** |  | [optional] 
**inputSchema** | **String** |  | [optional] 
**formId** | **String** |  | [optional] 
**formKey** | **String** |  | [optional] 
**calledElement** | **String** |  | [optional] 
**inVariables** | [**List<InVariable>**](InVariable.md) |  | [optional] [default to const []]
**outVariables** | [**List<OutVariable>**](OutVariable.md) |  | [optional] [default to const []]
**decisionRef** | **String** |  | [optional] 
**mapDecisionResult** | **String** |  | [optional] 
**hitPolicy** | **String** |  | [optional] 
**inputs** | [**List<DMNInputAST>**](DMNInputAST.md) |  | [optional] [default to const []]
**outputs** | [**List<DMNOutputAST>**](DMNOutputAST.md) |  | [optional] [default to const []]
**rules** | [**List<DMNRuleAST>**](DMNRuleAST.md) |  | [optional] [default to const []]
**attachedToRef** | **String** |  | [optional] 
**timeDuration** | **String** |  | [optional] 
**timeDate** | **String** |  | [optional] 
**timeCycle** | **String** |  | [optional] 
**cancelActivity** | **bool** |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


