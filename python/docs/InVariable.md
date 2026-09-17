# InVariable


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**source** | **str** |  | [optional] 
**target** | **str** |  | [optional] 
**variables** | **str** |  | [optional] 
**local** | **bool** |  | [optional] 

## Example

```python
from nativebpm_client.models.in_variable import InVariable

# TODO update the JSON string below
json = "{}"
# create an instance of InVariable from a JSON string
in_variable_instance = InVariable.from_json(json)
# print the JSON string representation of the object
print(InVariable.to_json())

# convert the object into a dict
in_variable_dict = in_variable_instance.to_dict()
# create an instance of InVariable from a dict
in_variable_from_dict = InVariable.from_dict(in_variable_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


