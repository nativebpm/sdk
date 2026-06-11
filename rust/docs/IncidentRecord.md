# IncidentRecord

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **uuid::Uuid** |  | 
**instance_id** | **uuid::Uuid** |  | 
**activity_id** | **String** |  | 
**error_message** | **String** |  | 
**stack_trace** | Option<**String**> |  | [optional]
**attempts_made** | **i32** |  | 
**resolved** | **bool** |  | 
**created_at** | **chrono::DateTime<chrono::FixedOffset>** |  | 
**resolved_at** | Option<**chrono::DateTime<chrono::FixedOffset>**> |  | [optional]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


