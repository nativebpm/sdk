# CompleteTaskRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**variables** | **Dict[str, object]** |  | [optional] 

## Example

```python
from nativebpm_client.models.complete_task_request import CompleteTaskRequest

# TODO update the JSON string below
json = "{}"
# create an instance of CompleteTaskRequest from a JSON string
complete_task_request_instance = CompleteTaskRequest.from_json(json)
# print the JSON string representation of the object
print(CompleteTaskRequest.to_json())

# convert the object into a dict
complete_task_request_dict = complete_task_request_instance.to_dict()
# create an instance of CompleteTaskRequest from a dict
complete_task_request_from_dict = CompleteTaskRequest.from_dict(complete_task_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


