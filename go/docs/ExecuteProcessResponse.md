# ExecuteProcessResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InstanceId** | **string** | Process instance identifier | 
**DefinitionId** | **string** | Process definition identifier | 
**Version** | **int32** | Version number of the process definition | 
**IsNewVersionDeployed** | **bool** | Indicates whether a new process definition version was deployed in this call | 
**Status** | **string** | Status of the created instance (e.g., ACTIVE, COMPLETED) | 
**State** | Pointer to **map[string]interface{}** | Current execution state variables and data | [optional] 
**CurrentTasks** | Pointer to [**[]TaskRecord**](TaskRecord.md) | Active tasks currently waiting for human intervention or external completion | [optional] 

## Methods

### NewExecuteProcessResponse

`func NewExecuteProcessResponse(instanceId string, definitionId string, version int32, isNewVersionDeployed bool, status string, ) *ExecuteProcessResponse`

NewExecuteProcessResponse instantiates a new ExecuteProcessResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExecuteProcessResponseWithDefaults

`func NewExecuteProcessResponseWithDefaults() *ExecuteProcessResponse`

NewExecuteProcessResponseWithDefaults instantiates a new ExecuteProcessResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInstanceId

`func (o *ExecuteProcessResponse) GetInstanceId() string`

GetInstanceId returns the InstanceId field if non-nil, zero value otherwise.

### GetInstanceIdOk

`func (o *ExecuteProcessResponse) GetInstanceIdOk() (*string, bool)`

GetInstanceIdOk returns a tuple with the InstanceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceId

`func (o *ExecuteProcessResponse) SetInstanceId(v string)`

SetInstanceId sets InstanceId field to given value.


### GetDefinitionId

`func (o *ExecuteProcessResponse) GetDefinitionId() string`

GetDefinitionId returns the DefinitionId field if non-nil, zero value otherwise.

### GetDefinitionIdOk

`func (o *ExecuteProcessResponse) GetDefinitionIdOk() (*string, bool)`

GetDefinitionIdOk returns a tuple with the DefinitionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefinitionId

`func (o *ExecuteProcessResponse) SetDefinitionId(v string)`

SetDefinitionId sets DefinitionId field to given value.


### GetVersion

`func (o *ExecuteProcessResponse) GetVersion() int32`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ExecuteProcessResponse) GetVersionOk() (*int32, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ExecuteProcessResponse) SetVersion(v int32)`

SetVersion sets Version field to given value.


### GetIsNewVersionDeployed

`func (o *ExecuteProcessResponse) GetIsNewVersionDeployed() bool`

GetIsNewVersionDeployed returns the IsNewVersionDeployed field if non-nil, zero value otherwise.

### GetIsNewVersionDeployedOk

`func (o *ExecuteProcessResponse) GetIsNewVersionDeployedOk() (*bool, bool)`

GetIsNewVersionDeployedOk returns a tuple with the IsNewVersionDeployed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsNewVersionDeployed

`func (o *ExecuteProcessResponse) SetIsNewVersionDeployed(v bool)`

SetIsNewVersionDeployed sets IsNewVersionDeployed field to given value.


### GetStatus

`func (o *ExecuteProcessResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ExecuteProcessResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ExecuteProcessResponse) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetState

`func (o *ExecuteProcessResponse) GetState() map[string]interface{}`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *ExecuteProcessResponse) GetStateOk() (*map[string]interface{}, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *ExecuteProcessResponse) SetState(v map[string]interface{})`

SetState sets State field to given value.

### HasState

`func (o *ExecuteProcessResponse) HasState() bool`

HasState returns a boolean if a field has been set.

### GetCurrentTasks

`func (o *ExecuteProcessResponse) GetCurrentTasks() []TaskRecord`

GetCurrentTasks returns the CurrentTasks field if non-nil, zero value otherwise.

### GetCurrentTasksOk

`func (o *ExecuteProcessResponse) GetCurrentTasksOk() (*[]TaskRecord, bool)`

GetCurrentTasksOk returns a tuple with the CurrentTasks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentTasks

`func (o *ExecuteProcessResponse) SetCurrentTasks(v []TaskRecord)`

SetCurrentTasks sets CurrentTasks field to given value.

### HasCurrentTasks

`func (o *ExecuteProcessResponse) HasCurrentTasks() bool`

HasCurrentTasks returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


