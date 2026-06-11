# WebhookDeliveryRecord


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **UUID** |  | 
**webhook_id** | **UUID** |  | 
**tenant_id** | **str** |  | 
**event_type** | **str** |  | 
**payload** | **bytes** |  | 
**status** | **str** |  | 
**response_code** | **int** |  | [optional] 
**response_body** | **str** |  | [optional] 
**attempts** | **int** |  | 
**next_retry** | **datetime** |  | [optional] 
**created_at** | **datetime** |  | 
**processed_at** | **datetime** |  | [optional] 

## Example

```python
from nativebpm_client.models.webhook_delivery_record import WebhookDeliveryRecord

# TODO update the JSON string below
json = "{}"
# create an instance of WebhookDeliveryRecord from a JSON string
webhook_delivery_record_instance = WebhookDeliveryRecord.from_json(json)
# print the JSON string representation of the object
print(WebhookDeliveryRecord.to_json())

# convert the object into a dict
webhook_delivery_record_dict = webhook_delivery_record_instance.to_dict()
# create an instance of WebhookDeliveryRecord from a dict
webhook_delivery_record_from_dict = WebhookDeliveryRecord.from_dict(webhook_delivery_record_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


