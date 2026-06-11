# TaskRecord


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **UUID** |  | 
**instance_id** | **UUID** |  | 
**activity_id** | **str** |  | 
**name** | **str** |  | 
**assignee** | **str** |  | 
**candidate_groups** | **str** |  | 
**status** | **str** |  | 
**due_date** | **datetime** |  | [optional] 
**input_schema** | **str** | JSON schema definition of form widgets | [optional] 
**created_at** | **datetime** |  | 
**claimed_at** | **datetime** |  | [optional] 
**completed_at** | **datetime** |  | [optional] 
**current_step** | **int** |  | [optional] 
**draft_variables** | **Dict[str, object]** |  | [optional] 

## Example

```python
from nativebpm_client.models.task_record import TaskRecord

# TODO update the JSON string below
json = "{}"
# create an instance of TaskRecord from a JSON string
task_record_instance = TaskRecord.from_json(json)
# print the JSON string representation of the object
print(TaskRecord.to_json())

# convert the object into a dict
task_record_dict = task_record_instance.to_dict()
# create an instance of TaskRecord from a dict
task_record_from_dict = TaskRecord.from_dict(task_record_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


