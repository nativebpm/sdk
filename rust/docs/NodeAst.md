# NodeAst

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **String** |  | 
**name** | **String** |  | 
**r#type** | **Type** |  (enum: startEvent, endEvent, serviceTask, userTask, aiServiceTask, aiTask, exclusiveGateway, parallelGateway, eventBasedGateway, callActivity, businessRuleTask, boundaryTimerEvent) | 
**topic** | Option<**String**> |  | [optional]
**wasm_path** | Option<**String**> |  | [optional]
**provider** | Option<**String**> |  | [optional]
**model** | Option<**String**> |  | [optional]
**prompt** | Option<**String**> |  | [optional]
**system_instruction** | Option<**String**> |  | [optional]
**response_schema** | Option<**String**> |  | [optional]
**temperature** | Option<**f64**> |  | [optional]
**result_var** | Option<**String**> |  | [optional]
**assignee** | Option<**String**> |  | [optional]
**candidate_groups** | Option<**String**> |  | [optional]
**due_date** | Option<**String**> |  | [optional]
**input_schema** | Option<**String**> |  | [optional]
**form_id** | Option<**String**> |  | [optional]
**form_key** | Option<**String**> |  | [optional]
**called_element** | Option<**String**> |  | [optional]
**in_variables** | Option<[**Vec<models::InVariable>**](InVariable.md)> |  | [optional]
**out_variables** | Option<[**Vec<models::OutVariable>**](OutVariable.md)> |  | [optional]
**decision_ref** | Option<**String**> |  | [optional]
**map_decision_result** | Option<**String**> |  | [optional]
**hit_policy** | Option<**String**> |  | [optional]
**inputs** | Option<[**Vec<models::DmnInputAst>**](DMNInputAST.md)> |  | [optional]
**outputs** | Option<[**Vec<models::DmnOutputAst>**](DMNOutputAST.md)> |  | [optional]
**rules** | Option<[**Vec<models::DmnRuleAst>**](DMNRuleAST.md)> |  | [optional]
**attached_to_ref** | Option<**String**> |  | [optional]
**time_duration** | Option<**String**> |  | [optional]
**time_date** | Option<**String**> |  | [optional]
**time_cycle** | Option<**String**> |  | [optional]
**cancel_activity** | Option<**bool**> |  | [optional]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


