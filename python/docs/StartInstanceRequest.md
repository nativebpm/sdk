# StartInstanceRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**instance_id** | **UUID** | Optional user-generated UUID to enforce idempotency | [optional] 
**business_key** | **str** | Business tracking keyword | [optional] 
**variables** | **Dict[str, object]** |  | [optional] 

## Example

```python
from nativebpm_client.models.start_instance_request import StartInstanceRequest

# TODO update the JSON string below
json = "{}"
# create an instance of StartInstanceRequest from a JSON string
start_instance_request_instance = StartInstanceRequest.from_json(json)
# print the JSON string representation of the object
print(StartInstanceRequest.to_json())

# convert the object into a dict
start_instance_request_dict = start_instance_request_instance.to_dict()
# create an instance of StartInstanceRequest from a dict
start_instance_request_from_dict = StartInstanceRequest.from_dict(start_instance_request_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


