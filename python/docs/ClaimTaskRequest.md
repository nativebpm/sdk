# ClaimTaskRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**assignee** | **str** |  | 

## Example

```python
from nativebpm_client.models.claim_task_request import ClaimTaskRequest

# TODO update the JSON string below
json = "{}"
# create an instance of ClaimTaskRequest from a JSON string
claim_task_request_instance = ClaimTaskRequest.from_json(json)
# print the JSON string representation of the object
print(ClaimTaskRequest.to_json())

# convert the object into a dict
claim_task_request_dict = claim_task_request_instance.to_dict()
# create an instance of ClaimTaskRequest from a dict
claim_task_request_from_dict = ClaimTaskRequest.from_dict(claim_task_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


