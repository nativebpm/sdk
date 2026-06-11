# CompleteInstanceTaskRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**node_id** | **str** | BPMN task element identifier | 
**variables** | **Dict[str, object]** |  | [optional] 

## Example

```python
from nativebpm_client.models.complete_instance_task_request import CompleteInstanceTaskRequest

# TODO update the JSON string below
json = "{}"
# create an instance of CompleteInstanceTaskRequest from a JSON string
complete_instance_task_request_instance = CompleteInstanceTaskRequest.from_json(json)
# print the JSON string representation of the object
print(CompleteInstanceTaskRequest.to_json())

# convert the object into a dict
complete_instance_task_request_dict = complete_instance_task_request_instance.to_dict()
# create an instance of CompleteInstanceTaskRequest from a dict
complete_instance_task_request_from_dict = CompleteInstanceTaskRequest.from_dict(complete_instance_task_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


