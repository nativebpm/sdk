# WorkflowAST

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Unique process definition identifier | 
**Name** | **string** | Human-readable process name | 
**InputSchema** | Pointer to **string** | JSON Schema contract for process start variables | [optional] 
**Nodes** | [**[]NodeAST**](NodeAST.md) |  | 
**Flows** | [**[]FlowAST**](FlowAST.md) |  | 

## Methods

### NewWorkflowAST

`func NewWorkflowAST(id string, name string, nodes []NodeAST, flows []FlowAST, ) *WorkflowAST`

NewWorkflowAST instantiates a new WorkflowAST object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowASTWithDefaults

`func NewWorkflowASTWithDefaults() *WorkflowAST`

NewWorkflowASTWithDefaults instantiates a new WorkflowAST object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *WorkflowAST) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *WorkflowAST) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *WorkflowAST) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *WorkflowAST) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *WorkflowAST) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *WorkflowAST) SetName(v string)`

SetName sets Name field to given value.


### GetInputSchema

`func (o *WorkflowAST) GetInputSchema() string`

GetInputSchema returns the InputSchema field if non-nil, zero value otherwise.

### GetInputSchemaOk

`func (o *WorkflowAST) GetInputSchemaOk() (*string, bool)`

GetInputSchemaOk returns a tuple with the InputSchema field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInputSchema

`func (o *WorkflowAST) SetInputSchema(v string)`

SetInputSchema sets InputSchema field to given value.

### HasInputSchema

`func (o *WorkflowAST) HasInputSchema() bool`

HasInputSchema returns a boolean if a field has been set.

### GetNodes

`func (o *WorkflowAST) GetNodes() []NodeAST`

GetNodes returns the Nodes field if non-nil, zero value otherwise.

### GetNodesOk

`func (o *WorkflowAST) GetNodesOk() (*[]NodeAST, bool)`

GetNodesOk returns a tuple with the Nodes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNodes

`func (o *WorkflowAST) SetNodes(v []NodeAST)`

SetNodes sets Nodes field to given value.


### GetFlows

`func (o *WorkflowAST) GetFlows() []FlowAST`

GetFlows returns the Flows field if non-nil, zero value otherwise.

### GetFlowsOk

`func (o *WorkflowAST) GetFlowsOk() (*[]FlowAST, bool)`

GetFlowsOk returns a tuple with the Flows field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFlows

`func (o *WorkflowAST) SetFlows(v []FlowAST)`

SetFlows sets Flows field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


