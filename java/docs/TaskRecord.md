

# TaskRecord


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**id** | **UUID** |  |  |
|**instanceId** | **UUID** |  |  |
|**activityId** | **String** |  |  |
|**name** | **String** |  |  |
|**assignee** | **String** |  |  |
|**candidateGroups** | **String** |  |  |
|**status** | **String** |  |  |
|**dueDate** | **OffsetDateTime** |  |  [optional] |
|**inputSchema** | **String** | JSON schema definition of form widgets |  [optional] |
|**formId** | **String** | Form identifier or Camunda form key for dynamic schema rendering |  [optional] |
|**createdAt** | **OffsetDateTime** |  |  |
|**claimedAt** | **OffsetDateTime** |  |  [optional] |
|**completedAt** | **OffsetDateTime** |  |  [optional] |
|**currentStep** | **Integer** |  |  [optional] |
|**draftVariables** | **Map&lt;String, Object&gt;** |  |  [optional] |



