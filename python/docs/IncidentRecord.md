# IncidentRecord


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **UUID** |  | 
**instance_id** | **UUID** |  | 
**activity_id** | **str** |  | 
**error_message** | **str** |  | 
**stack_trace** | **str** |  | [optional] 
**attempts_made** | **int** |  | 
**resolved** | **bool** |  | 
**created_at** | **datetime** |  | 
**resolved_at** | **datetime** |  | [optional] 

## Example

```python
from nativebpm_client.models.incident_record import IncidentRecord

# TODO update the JSON string below
json = "{}"
# create an instance of IncidentRecord from a JSON string
incident_record_instance = IncidentRecord.from_json(json)
# print the JSON string representation of the object
print(IncidentRecord.to_json())

# convert the object into a dict
incident_record_dict = incident_record_instance.to_dict()
# create an instance of IncidentRecord from a dict
incident_record_from_dict = IncidentRecord.from_dict(incident_record_dict)
```
[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


