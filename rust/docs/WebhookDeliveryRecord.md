# WebhookDeliveryRecord

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **uuid::Uuid** |  | 
**webhook_id** | **uuid::Uuid** |  | 
**tenant_id** | **String** |  | 
**event_type** | **String** |  | 
**payload** | **String** |  | 
**status** | **String** |  | 
**response_code** | Option<**i32**> |  | [optional]
**response_body** | Option<**String**> |  | [optional]
**attempts** | **i32** |  | 
**next_retry** | Option<**chrono::DateTime<chrono::FixedOffset>**> |  | [optional]
**created_at** | **chrono::DateTime<chrono::FixedOffset>** |  | 
**processed_at** | Option<**chrono::DateTime<chrono::FixedOffset>**> |  | [optional]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


