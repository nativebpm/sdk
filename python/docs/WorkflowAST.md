# WorkflowAST


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **str** | Unique process definition identifier | 
**name** | **str** | Human-readable process name | 
**input_schema** | **str** | JSON Schema contract for process start variables | [optional] 
**nodes** | [**List[NodeAST]**](NodeAST.md) |  | 
**flows** | [**List[FlowAST]**](FlowAST.md) |  | 

## Example

```python
from nativebpm_client.models.workflow_ast import WorkflowAST

# TODO update the JSON string below
json = "{}"
# create an instance of WorkflowAST from a JSON string
workflow_ast_instance = WorkflowAST.from_json(json)
# print the JSON string representation of the object
print(WorkflowAST.to_json())

# convert the object into a dict
workflow_ast_dict = workflow_ast_instance.to_dict()
# create an instance of WorkflowAST from a dict
workflow_ast_from_dict = WorkflowAST.from_dict(workflow_ast_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


