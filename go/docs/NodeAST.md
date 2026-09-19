# NodeAST

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Name** | **string** |  | 
**Type** | **string** |  | 
**Topic** | Pointer to **string** |  | [optional] 
**WasmPath** | Pointer to **string** |  | [optional] 
**Provider** | Pointer to **string** |  | [optional] 
**Model** | Pointer to **string** |  | [optional] 
**Prompt** | Pointer to **string** |  | [optional] 
**SystemInstruction** | Pointer to **string** |  | [optional] 
**ResponseSchema** | Pointer to **string** |  | [optional] 
**Temperature** | Pointer to **float32** |  | [optional] 
**ResultVar** | Pointer to **string** |  | [optional] 
**Assignee** | Pointer to **string** |  | [optional] 
**CandidateGroups** | Pointer to **string** |  | [optional] 
**DueDate** | Pointer to **string** |  | [optional] 
**InputSchema** | Pointer to **string** |  | [optional] 
**FormId** | Pointer to **string** |  | [optional] 
**FormKey** | Pointer to **string** |  | [optional] 
**CalledElement** | Pointer to **string** |  | [optional] 
**InVariables** | Pointer to [**[]InVariable**](InVariable.md) |  | [optional] 
**OutVariables** | Pointer to [**[]OutVariable**](OutVariable.md) |  | [optional] 
**DecisionRef** | Pointer to **string** |  | [optional] 
**MapDecisionResult** | Pointer to **string** |  | [optional] 
**HitPolicy** | Pointer to **string** |  | [optional] 
**Inputs** | Pointer to [**[]DMNInputAST**](DMNInputAST.md) |  | [optional] 
**Outputs** | Pointer to [**[]DMNOutputAST**](DMNOutputAST.md) |  | [optional] 
**Rules** | Pointer to [**[]DMNRuleAST**](DMNRuleAST.md) |  | [optional] 
**AttachedToRef** | Pointer to **string** |  | [optional] 
**TimeDuration** | Pointer to **string** |  | [optional] 
**TimeDate** | Pointer to **string** |  | [optional] 
**TimeCycle** | Pointer to **string** |  | [optional] 
**CancelActivity** | Pointer to **bool** |  | [optional] 

## Methods

### NewNodeAST

`func NewNodeAST(id string, name string, type_ string, ) *NodeAST`

NewNodeAST instantiates a new NodeAST object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNodeASTWithDefaults

`func NewNodeASTWithDefaults() *NodeAST`

NewNodeASTWithDefaults instantiates a new NodeAST object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *NodeAST) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *NodeAST) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *NodeAST) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *NodeAST) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NodeAST) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NodeAST) SetName(v string)`

SetName sets Name field to given value.


### GetType

`func (o *NodeAST) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *NodeAST) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *NodeAST) SetType(v string)`

SetType sets Type field to given value.


### GetTopic

`func (o *NodeAST) GetTopic() string`

GetTopic returns the Topic field if non-nil, zero value otherwise.

### GetTopicOk

`func (o *NodeAST) GetTopicOk() (*string, bool)`

GetTopicOk returns a tuple with the Topic field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTopic

`func (o *NodeAST) SetTopic(v string)`

SetTopic sets Topic field to given value.

### HasTopic

`func (o *NodeAST) HasTopic() bool`

HasTopic returns a boolean if a field has been set.

### GetWasmPath

`func (o *NodeAST) GetWasmPath() string`

GetWasmPath returns the WasmPath field if non-nil, zero value otherwise.

### GetWasmPathOk

`func (o *NodeAST) GetWasmPathOk() (*string, bool)`

GetWasmPathOk returns a tuple with the WasmPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWasmPath

`func (o *NodeAST) SetWasmPath(v string)`

SetWasmPath sets WasmPath field to given value.

### HasWasmPath

`func (o *NodeAST) HasWasmPath() bool`

HasWasmPath returns a boolean if a field has been set.

### GetProvider

`func (o *NodeAST) GetProvider() string`

GetProvider returns the Provider field if non-nil, zero value otherwise.

### GetProviderOk

`func (o *NodeAST) GetProviderOk() (*string, bool)`

GetProviderOk returns a tuple with the Provider field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvider

`func (o *NodeAST) SetProvider(v string)`

SetProvider sets Provider field to given value.

### HasProvider

`func (o *NodeAST) HasProvider() bool`

HasProvider returns a boolean if a field has been set.

### GetModel

`func (o *NodeAST) GetModel() string`

GetModel returns the Model field if non-nil, zero value otherwise.

### GetModelOk

`func (o *NodeAST) GetModelOk() (*string, bool)`

GetModelOk returns a tuple with the Model field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetModel

`func (o *NodeAST) SetModel(v string)`

SetModel sets Model field to given value.

### HasModel

`func (o *NodeAST) HasModel() bool`

HasModel returns a boolean if a field has been set.

### GetPrompt

`func (o *NodeAST) GetPrompt() string`

GetPrompt returns the Prompt field if non-nil, zero value otherwise.

### GetPromptOk

`func (o *NodeAST) GetPromptOk() (*string, bool)`

GetPromptOk returns a tuple with the Prompt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrompt

`func (o *NodeAST) SetPrompt(v string)`

SetPrompt sets Prompt field to given value.

### HasPrompt

`func (o *NodeAST) HasPrompt() bool`

HasPrompt returns a boolean if a field has been set.

### GetSystemInstruction

`func (o *NodeAST) GetSystemInstruction() string`

GetSystemInstruction returns the SystemInstruction field if non-nil, zero value otherwise.

### GetSystemInstructionOk

`func (o *NodeAST) GetSystemInstructionOk() (*string, bool)`

GetSystemInstructionOk returns a tuple with the SystemInstruction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSystemInstruction

`func (o *NodeAST) SetSystemInstruction(v string)`

SetSystemInstruction sets SystemInstruction field to given value.

### HasSystemInstruction

`func (o *NodeAST) HasSystemInstruction() bool`

HasSystemInstruction returns a boolean if a field has been set.

### GetResponseSchema

`func (o *NodeAST) GetResponseSchema() string`

GetResponseSchema returns the ResponseSchema field if non-nil, zero value otherwise.

### GetResponseSchemaOk

`func (o *NodeAST) GetResponseSchemaOk() (*string, bool)`

GetResponseSchemaOk returns a tuple with the ResponseSchema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResponseSchema

`func (o *NodeAST) SetResponseSchema(v string)`

SetResponseSchema sets ResponseSchema field to given value.

### HasResponseSchema

`func (o *NodeAST) HasResponseSchema() bool`

HasResponseSchema returns a boolean if a field has been set.

### GetTemperature

`func (o *NodeAST) GetTemperature() float32`

GetTemperature returns the Temperature field if non-nil, zero value otherwise.

### GetTemperatureOk

`func (o *NodeAST) GetTemperatureOk() (*float32, bool)`

GetTemperatureOk returns a tuple with the Temperature field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTemperature

`func (o *NodeAST) SetTemperature(v float32)`

SetTemperature sets Temperature field to given value.

### HasTemperature

`func (o *NodeAST) HasTemperature() bool`

HasTemperature returns a boolean if a field has been set.

### GetResultVar

`func (o *NodeAST) GetResultVar() string`

GetResultVar returns the ResultVar field if non-nil, zero value otherwise.

### GetResultVarOk

`func (o *NodeAST) GetResultVarOk() (*string, bool)`

GetResultVarOk returns a tuple with the ResultVar field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResultVar

`func (o *NodeAST) SetResultVar(v string)`

SetResultVar sets ResultVar field to given value.

### HasResultVar

`func (o *NodeAST) HasResultVar() bool`

HasResultVar returns a boolean if a field has been set.

### GetAssignee

`func (o *NodeAST) GetAssignee() string`

GetAssignee returns the Assignee field if non-nil, zero value otherwise.

### GetAssigneeOk

`func (o *NodeAST) GetAssigneeOk() (*string, bool)`

GetAssigneeOk returns a tuple with the Assignee field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAssignee

`func (o *NodeAST) SetAssignee(v string)`

SetAssignee sets Assignee field to given value.

### HasAssignee

`func (o *NodeAST) HasAssignee() bool`

HasAssignee returns a boolean if a field has been set.

### GetCandidateGroups

`func (o *NodeAST) GetCandidateGroups() string`

GetCandidateGroups returns the CandidateGroups field if non-nil, zero value otherwise.

### GetCandidateGroupsOk

`func (o *NodeAST) GetCandidateGroupsOk() (*string, bool)`

GetCandidateGroupsOk returns a tuple with the CandidateGroups field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCandidateGroups

`func (o *NodeAST) SetCandidateGroups(v string)`

SetCandidateGroups sets CandidateGroups field to given value.

### HasCandidateGroups

`func (o *NodeAST) HasCandidateGroups() bool`

HasCandidateGroups returns a boolean if a field has been set.

### GetDueDate

`func (o *NodeAST) GetDueDate() string`

GetDueDate returns the DueDate field if non-nil, zero value otherwise.

### GetDueDateOk

`func (o *NodeAST) GetDueDateOk() (*string, bool)`

GetDueDateOk returns a tuple with the DueDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDueDate

`func (o *NodeAST) SetDueDate(v string)`

SetDueDate sets DueDate field to given value.

### HasDueDate

`func (o *NodeAST) HasDueDate() bool`

HasDueDate returns a boolean if a field has been set.

### GetInputSchema

`func (o *NodeAST) GetInputSchema() string`

GetInputSchema returns the InputSchema field if non-nil, zero value otherwise.

### GetInputSchemaOk

`func (o *NodeAST) GetInputSchemaOk() (*string, bool)`

GetInputSchemaOk returns a tuple with the InputSchema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInputSchema

`func (o *NodeAST) SetInputSchema(v string)`

SetInputSchema sets InputSchema field to given value.

### HasInputSchema

`func (o *NodeAST) HasInputSchema() bool`

HasInputSchema returns a boolean if a field has been set.

### GetFormId

`func (o *NodeAST) GetFormId() string`

GetFormId returns the FormId field if non-nil, zero value otherwise.

### GetFormIdOk

`func (o *NodeAST) GetFormIdOk() (*string, bool)`

GetFormIdOk returns a tuple with the FormId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFormId

`func (o *NodeAST) SetFormId(v string)`

SetFormId sets FormId field to given value.

### HasFormId

`func (o *NodeAST) HasFormId() bool`

HasFormId returns a boolean if a field has been set.

### GetFormKey

`func (o *NodeAST) GetFormKey() string`

GetFormKey returns the FormKey field if non-nil, zero value otherwise.

### GetFormKeyOk

`func (o *NodeAST) GetFormKeyOk() (*string, bool)`

GetFormKeyOk returns a tuple with the FormKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFormKey

`func (o *NodeAST) SetFormKey(v string)`

SetFormKey sets FormKey field to given value.

### HasFormKey

`func (o *NodeAST) HasFormKey() bool`

HasFormKey returns a boolean if a field has been set.

### GetCalledElement

`func (o *NodeAST) GetCalledElement() string`

GetCalledElement returns the CalledElement field if non-nil, zero value otherwise.

### GetCalledElementOk

`func (o *NodeAST) GetCalledElementOk() (*string, bool)`

GetCalledElementOk returns a tuple with the CalledElement field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCalledElement

`func (o *NodeAST) SetCalledElement(v string)`

SetCalledElement sets CalledElement field to given value.

### HasCalledElement

`func (o *NodeAST) HasCalledElement() bool`

HasCalledElement returns a boolean if a field has been set.

### GetInVariables

`func (o *NodeAST) GetInVariables() []InVariable`

GetInVariables returns the InVariables field if non-nil, zero value otherwise.

### GetInVariablesOk

`func (o *NodeAST) GetInVariablesOk() (*[]InVariable, bool)`

GetInVariablesOk returns a tuple with the InVariables field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInVariables

`func (o *NodeAST) SetInVariables(v []InVariable)`

SetInVariables sets InVariables field to given value.

### HasInVariables

`func (o *NodeAST) HasInVariables() bool`

HasInVariables returns a boolean if a field has been set.

### GetOutVariables

`func (o *NodeAST) GetOutVariables() []OutVariable`

GetOutVariables returns the OutVariables field if non-nil, zero value otherwise.

### GetOutVariablesOk

`func (o *NodeAST) GetOutVariablesOk() (*[]OutVariable, bool)`

GetOutVariablesOk returns a tuple with the OutVariables field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutVariables

`func (o *NodeAST) SetOutVariables(v []OutVariable)`

SetOutVariables sets OutVariables field to given value.

### HasOutVariables

`func (o *NodeAST) HasOutVariables() bool`

HasOutVariables returns a boolean if a field has been set.

### GetDecisionRef

`func (o *NodeAST) GetDecisionRef() string`

GetDecisionRef returns the DecisionRef field if non-nil, zero value otherwise.

### GetDecisionRefOk

`func (o *NodeAST) GetDecisionRefOk() (*string, bool)`

GetDecisionRefOk returns a tuple with the DecisionRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDecisionRef

`func (o *NodeAST) SetDecisionRef(v string)`

SetDecisionRef sets DecisionRef field to given value.

### HasDecisionRef

`func (o *NodeAST) HasDecisionRef() bool`

HasDecisionRef returns a boolean if a field has been set.

### GetMapDecisionResult

`func (o *NodeAST) GetMapDecisionResult() string`

GetMapDecisionResult returns the MapDecisionResult field if non-nil, zero value otherwise.

### GetMapDecisionResultOk

`func (o *NodeAST) GetMapDecisionResultOk() (*string, bool)`

GetMapDecisionResultOk returns a tuple with the MapDecisionResult field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMapDecisionResult

`func (o *NodeAST) SetMapDecisionResult(v string)`

SetMapDecisionResult sets MapDecisionResult field to given value.

### HasMapDecisionResult

`func (o *NodeAST) HasMapDecisionResult() bool`

HasMapDecisionResult returns a boolean if a field has been set.

### GetHitPolicy

`func (o *NodeAST) GetHitPolicy() string`

GetHitPolicy returns the HitPolicy field if non-nil, zero value otherwise.

### GetHitPolicyOk

`func (o *NodeAST) GetHitPolicyOk() (*string, bool)`

GetHitPolicyOk returns a tuple with the HitPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHitPolicy

`func (o *NodeAST) SetHitPolicy(v string)`

SetHitPolicy sets HitPolicy field to given value.

### HasHitPolicy

`func (o *NodeAST) HasHitPolicy() bool`

HasHitPolicy returns a boolean if a field has been set.

### GetInputs

`func (o *NodeAST) GetInputs() []DMNInputAST`

GetInputs returns the Inputs field if non-nil, zero value otherwise.

### GetInputsOk

`func (o *NodeAST) GetInputsOk() (*[]DMNInputAST, bool)`

GetInputsOk returns a tuple with the Inputs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInputs

`func (o *NodeAST) SetInputs(v []DMNInputAST)`

SetInputs sets Inputs field to given value.

### HasInputs

`func (o *NodeAST) HasInputs() bool`

HasInputs returns a boolean if a field has been set.

### GetOutputs

`func (o *NodeAST) GetOutputs() []DMNOutputAST`

GetOutputs returns the Outputs field if non-nil, zero value otherwise.

### GetOutputsOk

`func (o *NodeAST) GetOutputsOk() (*[]DMNOutputAST, bool)`

GetOutputsOk returns a tuple with the Outputs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutputs

`func (o *NodeAST) SetOutputs(v []DMNOutputAST)`

SetOutputs sets Outputs field to given value.

### HasOutputs

`func (o *NodeAST) HasOutputs() bool`

HasOutputs returns a boolean if a field has been set.

### GetRules

`func (o *NodeAST) GetRules() []DMNRuleAST`

GetRules returns the Rules field if non-nil, zero value otherwise.

### GetRulesOk

`func (o *NodeAST) GetRulesOk() (*[]DMNRuleAST, bool)`

GetRulesOk returns a tuple with the Rules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRules

`func (o *NodeAST) SetRules(v []DMNRuleAST)`

SetRules sets Rules field to given value.

### HasRules

`func (o *NodeAST) HasRules() bool`

HasRules returns a boolean if a field has been set.

### GetAttachedToRef

`func (o *NodeAST) GetAttachedToRef() string`

GetAttachedToRef returns the AttachedToRef field if non-nil, zero value otherwise.

### GetAttachedToRefOk

`func (o *NodeAST) GetAttachedToRefOk() (*string, bool)`

GetAttachedToRefOk returns a tuple with the AttachedToRef field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttachedToRef

`func (o *NodeAST) SetAttachedToRef(v string)`

SetAttachedToRef sets AttachedToRef field to given value.

### HasAttachedToRef

`func (o *NodeAST) HasAttachedToRef() bool`

HasAttachedToRef returns a boolean if a field has been set.

### GetTimeDuration

`func (o *NodeAST) GetTimeDuration() string`

GetTimeDuration returns the TimeDuration field if non-nil, zero value otherwise.

### GetTimeDurationOk

`func (o *NodeAST) GetTimeDurationOk() (*string, bool)`

GetTimeDurationOk returns a tuple with the TimeDuration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeDuration

`func (o *NodeAST) SetTimeDuration(v string)`

SetTimeDuration sets TimeDuration field to given value.

### HasTimeDuration

`func (o *NodeAST) HasTimeDuration() bool`

HasTimeDuration returns a boolean if a field has been set.

### GetTimeDate

`func (o *NodeAST) GetTimeDate() string`

GetTimeDate returns the TimeDate field if non-nil, zero value otherwise.

### GetTimeDateOk

`func (o *NodeAST) GetTimeDateOk() (*string, bool)`

GetTimeDateOk returns a tuple with the TimeDate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeDate

`func (o *NodeAST) SetTimeDate(v string)`

SetTimeDate sets TimeDate field to given value.

### HasTimeDate

`func (o *NodeAST) HasTimeDate() bool`

HasTimeDate returns a boolean if a field has been set.

### GetTimeCycle

`func (o *NodeAST) GetTimeCycle() string`

GetTimeCycle returns the TimeCycle field if non-nil, zero value otherwise.

### GetTimeCycleOk

`func (o *NodeAST) GetTimeCycleOk() (*string, bool)`

GetTimeCycleOk returns a tuple with the TimeCycle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeCycle

`func (o *NodeAST) SetTimeCycle(v string)`

SetTimeCycle sets TimeCycle field to given value.

### HasTimeCycle

`func (o *NodeAST) HasTimeCycle() bool`

HasTimeCycle returns a boolean if a field has been set.

### GetCancelActivity

`func (o *NodeAST) GetCancelActivity() bool`

GetCancelActivity returns the CancelActivity field if non-nil, zero value otherwise.

### GetCancelActivityOk

`func (o *NodeAST) GetCancelActivityOk() (*bool, bool)`

GetCancelActivityOk returns a tuple with the CancelActivity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCancelActivity

`func (o *NodeAST) SetCancelActivity(v bool)`

SetCancelActivity sets CancelActivity field to given value.

### HasCancelActivity

`func (o *NodeAST) HasCancelActivity() bool`

HasCancelActivity returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


