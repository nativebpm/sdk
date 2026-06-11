# TaskRecord

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **uuid::Uuid** |  | 
**instance_id** | **uuid::Uuid** |  | 
**activity_id** | **String** |  | 
**name** | **String** |  | 
**assignee** | **String** |  | 
**candidate_groups** | **String** |  | 
**status** | **String** |  | 
**due_date** | Option<**chrono::DateTime<chrono::FixedOffset>**> |  | [optional]
**input_schema** | Option<**String**> | JSON schema definition of form widgets | [optional]
**created_at** | **chrono::DateTime<chrono::FixedOffset>** |  | 
**claimed_at** | Option<**chrono::DateTime<chrono::FixedOffset>**> |  | [optional]
**completed_at** | Option<**chrono::DateTime<chrono::FixedOffset>**> |  | [optional]
**current_step** | Option<**i32**> |  | [optional]
**draft_variables** | Option<**std::collections::HashMap<String, serde_json::Value>**> |  | [optional]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


