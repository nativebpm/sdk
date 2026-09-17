# nativebpm_client.model.WorkflowAST

## Load the model package
```dart
import 'package:nativebpm_client/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique process definition identifier | 
**name** | **String** | Human-readable process name | 
**inputSchema** | **String** | JSON Schema contract for process start variables | [optional] 
**nodes** | [**List<NodeAST>**](NodeAST.md) |  | [default to const []]
**flows** | [**List<FlowAST>**](FlowAST.md) |  | [default to const []]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


