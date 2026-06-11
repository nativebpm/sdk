# WebhookRecord


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **UUID** |  | 
**tenant_id** | **str** |  | 
**url** | **str** |  | 
**secret** | **str** |  | [optional] 
**events** | **List[str]** |  | 
**process_id** | **str** |  | [optional] 
**is_active** | **bool** |  | 
**enable_audit** | **bool** |  | 
**status** | **str** |  | 
**created_at** | **datetime** |  | 

## Example

```python
from nativebpm_client.models.webhook_record import WebhookRecord

# TODO update the JSON string below
json = "{}"
# create an instance of WebhookRecord from a JSON string
webhook_record_instance = WebhookRecord.from_json(json)
# print the JSON string representation of the object
print(WebhookRecord.to_json())

# convert the object into a dict
webhook_record_dict = webhook_record_instance.to_dict()
# create an instance of WebhookRecord from a dict
webhook_record_from_dict = WebhookRecord.from_dict(webhook_record_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


