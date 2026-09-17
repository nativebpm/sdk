# FlowAST


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **str** |  | 
**source** | **str** |  | 
**target** | **str** |  | 
**condition** | **str** |  | [optional] 

## Example

```python
from nativebpm_client.models.flow_ast import FlowAST

# TODO update the JSON string below
json = "{}"
# create an instance of FlowAST from a JSON string
flow_ast_instance = FlowAST.from_json(json)
# print the JSON string representation of the object
print(FlowAST.to_json())

# convert the object into a dict
flow_ast_dict = flow_ast_instance.to_dict()
# create an instance of FlowAST from a dict
flow_ast_from_dict = FlowAST.from_dict(flow_ast_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


