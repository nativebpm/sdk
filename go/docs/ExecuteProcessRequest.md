# ExecuteProcessRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Ast** | Pointer to [**WorkflowAST**](WorkflowAST.md) |  | [optional] 
**Forms** | Pointer to **map[string]map[string]interface{}** | Optional map of formId to JSON Schema definition objects | [optional] 
**ContentHash** | Pointer to **string** | Deterministic SHA-256 hash of the workflow AST and forms | [optional] 
**DefinitionId** | Pointer to **string** | Process definition identifier (optional if specified in AST) | [optional] 
**BusinessKey** | Pointer to **string** | Optional business correlation key | [optional] 
**Variables** | Pointer to **map[string]interface{}** | Initial variables to start the process with | [optional] 
**Signature** | Pointer to **string** | Optional HMAC signature verifying workflow authenticity | [optional] 

## Methods

### NewExecuteProcessRequest

`func NewExecuteProcessRequest() *ExecuteProcessRequest`

NewExecuteProcessRequest instantiates a new ExecuteProcessRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExecuteProcessRequestWithDefaults

`func NewExecuteProcessRequestWithDefaults() *ExecuteProcessRequest`

NewExecuteProcessRequestWithDefaults instantiates a new ExecuteProcessRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAst

`func (o *ExecuteProcessRequest) GetAst() WorkflowAST`

GetAst returns the Ast field if non-nil, zero value otherwise.

### GetAstOk

`func (o *ExecuteProcessRequest) GetAstOk() (*WorkflowAST, bool)`

GetAstOk returns a tuple with the Ast field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAst

`func (o *ExecuteProcessRequest) SetAst(v WorkflowAST)`

SetAst sets Ast field to given value.

### HasAst

`func (o *ExecuteProcessRequest) HasAst() bool`

HasAst returns a boolean if a field has been set.

### GetForms

`func (o *ExecuteProcessRequest) GetForms() map[string]map[string]interface{}`

GetForms returns the Forms field if non-nil, zero value otherwise.

### GetFormsOk

`func (o *ExecuteProcessRequest) GetFormsOk() (*map[string]map[string]interface{}, bool)`

GetFormsOk returns a tuple with the Forms field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetForms

`func (o *ExecuteProcessRequest) SetForms(v map[string]map[string]interface{})`

SetForms sets Forms field to given value.

### HasForms

`func (o *ExecuteProcessRequest) HasForms() bool`

HasForms returns a boolean if a field has been set.

### GetContentHash

`func (o *ExecuteProcessRequest) GetContentHash() string`

GetContentHash returns the ContentHash field if non-nil, zero value otherwise.

### GetContentHashOk

`func (o *ExecuteProcessRequest) GetContentHashOk() (*string, bool)`

GetContentHashOk returns a tuple with the ContentHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContentHash

`func (o *ExecuteProcessRequest) SetContentHash(v string)`

SetContentHash sets ContentHash field to given value.

### HasContentHash

`func (o *ExecuteProcessRequest) HasContentHash() bool`

HasContentHash returns a boolean if a field has been set.

### GetDefinitionId

`func (o *ExecuteProcessRequest) GetDefinitionId() string`

GetDefinitionId returns the DefinitionId field if non-nil, zero value otherwise.

### GetDefinitionIdOk

`func (o *ExecuteProcessRequest) GetDefinitionIdOk() (*string, bool)`

GetDefinitionIdOk returns a tuple with the DefinitionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefinitionId

`func (o *ExecuteProcessRequest) SetDefinitionId(v string)`

SetDefinitionId sets DefinitionId field to given value.

### HasDefinitionId

`func (o *ExecuteProcessRequest) HasDefinitionId() bool`

HasDefinitionId returns a boolean if a field has been set.

### GetBusinessKey

`func (o *ExecuteProcessRequest) GetBusinessKey() string`

GetBusinessKey returns the BusinessKey field if non-nil, zero value otherwise.

### GetBusinessKeyOk

`func (o *ExecuteProcessRequest) GetBusinessKeyOk() (*string, bool)`

GetBusinessKeyOk returns a tuple with the BusinessKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBusinessKey

`func (o *ExecuteProcessRequest) SetBusinessKey(v string)`

SetBusinessKey sets BusinessKey field to given value.

### HasBusinessKey

`func (o *ExecuteProcessRequest) HasBusinessKey() bool`

HasBusinessKey returns a boolean if a field has been set.

### GetVariables

`func (o *ExecuteProcessRequest) GetVariables() map[string]interface{}`

GetVariables returns the Variables field if non-nil, zero value otherwise.

### GetVariablesOk

`func (o *ExecuteProcessRequest) GetVariablesOk() (*map[string]interface{}, bool)`

GetVariablesOk returns a tuple with the Variables field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVariables

`func (o *ExecuteProcessRequest) SetVariables(v map[string]interface{})`

SetVariables sets Variables field to given value.

### HasVariables

`func (o *ExecuteProcessRequest) HasVariables() bool`

HasVariables returns a boolean if a field has been set.

### GetSignature

`func (o *ExecuteProcessRequest) GetSignature() string`

GetSignature returns the Signature field if non-nil, zero value otherwise.

### GetSignatureOk

`func (o *ExecuteProcessRequest) GetSignatureOk() (*string, bool)`

GetSignatureOk returns a tuple with the Signature field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignature

`func (o *ExecuteProcessRequest) SetSignature(v string)`

SetSignature sets Signature field to given value.

### HasSignature

`func (o *ExecuteProcessRequest) HasSignature() bool`

HasSignature returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


