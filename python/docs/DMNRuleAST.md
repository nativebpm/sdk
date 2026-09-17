# DMNRuleAST


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**inputs** | **List[str]** |  | 
**outputs** | **List[str]** |  | 

## Example

```python
from nativebpm_client.models.dmn_rule_ast import DMNRuleAST

# TODO update the JSON string below
json = "{}"
# create an instance of DMNRuleAST from a JSON string
dmn_rule_ast_instance = DMNRuleAST.from_json(json)
# print the JSON string representation of the object
print(DMNRuleAST.to_json())

# convert the object into a dict
dmn_rule_ast_dict = dmn_rule_ast_instance.to_dict()
# create an instance of DMNRuleAST from a dict
dmn_rule_ast_from_dict = DMNRuleAST.from_dict(dmn_rule_ast_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


