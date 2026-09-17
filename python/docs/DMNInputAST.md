# DMNInputAST


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**expression** | **str** |  | 
**type** | **str** |  | 

## Example

```python
from nativebpm_client.models.dmn_input_ast import DMNInputAST

# TODO update the JSON string below
json = "{}"
# create an instance of DMNInputAST from a JSON string
dmn_input_ast_instance = DMNInputAST.from_json(json)
# print the JSON string representation of the object
print(DMNInputAST.to_json())

# convert the object into a dict
dmn_input_ast_dict = dmn_input_ast_instance.to_dict()
# create an instance of DMNInputAST from a dict
dmn_input_ast_from_dict = DMNInputAST.from_dict(dmn_input_ast_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


