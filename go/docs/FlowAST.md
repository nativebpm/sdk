# FlowAST

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Source** | **string** |  | 
**Target** | **string** |  | 
**Condition** | Pointer to **string** |  | [optional] 

## Methods

### NewFlowAST

`func NewFlowAST(id string, source string, target string, ) *FlowAST`

NewFlowAST instantiates a new FlowAST object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewFlowASTWithDefaults

`func NewFlowASTWithDefaults() *FlowAST`

NewFlowASTWithDefaults instantiates a new FlowAST object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *FlowAST) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *FlowAST) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *FlowAST) SetId(v string)`

SetId sets Id field to given value.


### GetSource

`func (o *FlowAST) GetSource() string`

GetSource returns the Source field if non-nil, zero value otherwise.

### GetSourceOk

`func (o *FlowAST) GetSourceOk() (*string, bool)`

GetSourceOk returns a tuple with the Source field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSource

`func (o *FlowAST) SetSource(v string)`

SetSource sets Source field to given value.


### GetTarget

`func (o *FlowAST) GetTarget() string`

GetTarget returns the Target field if non-nil, zero value otherwise.

### GetTargetOk

`func (o *FlowAST) GetTargetOk() (*string, bool)`

GetTargetOk returns a tuple with the Target field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTarget

`func (o *FlowAST) SetTarget(v string)`

SetTarget sets Target field to given value.


### GetCondition

`func (o *FlowAST) GetCondition() string`

GetCondition returns the Condition field if non-nil, zero value otherwise.

### GetConditionOk

`func (o *FlowAST) GetConditionOk() (*string, bool)`

GetConditionOk returns a tuple with the Condition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCondition

`func (o *FlowAST) SetCondition(v string)`

SetCondition sets Condition field to given value.

### HasCondition

`func (o *FlowAST) HasCondition() bool`

HasCondition returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


