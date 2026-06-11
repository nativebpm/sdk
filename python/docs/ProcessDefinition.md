# ProcessDefinition


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**hash** | **str** | MD5/SHA256 content hash of the process XML schema definition | 
**id** | **str** | Unique process definition identifier | 
**name** | **str** | Friendly name of the process | 
**xml_data** | **bytes** | Base64-encoded raw BPMN 2.0 XML schema data | 
**deployed_at** | **datetime** |  | 

## Example

```python
from nativebpm_client.models.process_definition import ProcessDefinition

# TODO update the JSON string below
json = "{}"
# create an instance of ProcessDefinition from a JSON string
process_definition_instance = ProcessDefinition.from_json(json)
# print the JSON string representation of the object
print(ProcessDefinition.to_json())

# convert the object into a dict
process_definition_dict = process_definition_instance.to_dict()
# create an instance of ProcessDefinition from a dict
process_definition_from_dict = ProcessDefinition.from_dict(process_definition_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


