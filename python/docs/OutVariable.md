# OutVariable


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**source** | **str** |  | [optional] 
**target** | **str** |  | [optional] 
**variables** | **str** |  | [optional] 

## Example

```python
from nativebpm_client.models.out_variable import OutVariable

# TODO update the JSON string below
json = "{}"
# create an instance of OutVariable from a JSON string
out_variable_instance = OutVariable.from_json(json)
# print the JSON string representation of the object
print(OutVariable.to_json())

# convert the object into a dict
out_variable_dict = out_variable_instance.to_dict()
# create an instance of OutVariable from a dict
out_variable_from_dict = OutVariable.from_dict(out_variable_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


