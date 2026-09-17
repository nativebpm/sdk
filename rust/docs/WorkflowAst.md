# WorkflowAst

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** | Unique process definition identifier | 
**name** | **String** | Human-readable process name | 
**input_schema** | Option<**String**> | JSON Schema contract for process start variables | [optional]
**nodes** | [**Vec<models::NodeAst>**](NodeAST.md) |  | 
**flows** | [**Vec<models::FlowAst>**](FlowAST.md) |  | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


