# NodeAST


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **str** |  | 
**name** | **str** |  | 
**type** | **str** |  | 
**topic** | **str** |  | [optional] 
**wasm_path** | **str** |  | [optional] 
**provider** | **str** |  | [optional] 
**model** | **str** |  | [optional] 
**prompt** | **str** |  | [optional] 
**system_instruction** | **str** |  | [optional] 
**response_schema** | **str** |  | [optional] 
**temperature** | **float** |  | [optional] 
**result_var** | **str** |  | [optional] 
**assignee** | **str** |  | [optional] 
**candidate_groups** | **str** |  | [optional] 
**due_date** | **str** |  | [optional] 
**input_schema** | **str** |  | [optional] 
**form_id** | **str** |  | [optional] 
**form_key** | **str** |  | [optional] 
**called_element** | **str** |  | [optional] 
**in_variables** | [**List[InVariable]**](InVariable.md) |  | [optional] 
**out_variables** | [**List[OutVariable]**](OutVariable.md) |  | [optional] 
**decision_ref** | **str** |  | [optional] 
**map_decision_result** | **str** |  | [optional] 
**hit_policy** | **str** |  | [optional] 
**inputs** | [**List[DMNInputAST]**](DMNInputAST.md) |  | [optional] 
**outputs** | [**List[DMNOutputAST]**](DMNOutputAST.md) |  | [optional] 
**rules** | [**List[DMNRuleAST]**](DMNRuleAST.md) |  | [optional] 
**attached_to_ref** | **str** |  | [optional] 
**time_duration** | **str** |  | [optional] 
**time_date** | **str** |  | [optional] 
**time_cycle** | **str** |  | [optional] 
**cancel_activity** | **bool** |  | [optional] 

## Example

```python
from nativebpm_client.models.node_ast import NodeAST

# TODO update the JSON string below
json = "{}"
# create an instance of NodeAST from a JSON string
node_ast_instance = NodeAST.from_json(json)
# print the JSON string representation of the object
print(NodeAST.to_json())

# convert the object into a dict
node_ast_dict = node_ast_instance.to_dict()
# create an instance of NodeAST from a dict
node_ast_from_dict = NodeAST.from_dict(node_ast_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


