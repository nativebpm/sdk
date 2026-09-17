
# NodeAST


## Properties

Name | Type
------------ | -------------
`id` | string
`name` | string
`type` | string
`topic` | string
`wasmPath` | string
`provider` | string
`model` | string
`prompt` | string
`systemInstruction` | string
`responseSchema` | string
`temperature` | number
`resultVar` | string
`assignee` | string
`candidateGroups` | string
`dueDate` | string
`inputSchema` | string
`formId` | string
`formKey` | string
`calledElement` | string
`inVariables` | [Array&lt;InVariable&gt;](InVariable.md)
`outVariables` | [Array&lt;OutVariable&gt;](OutVariable.md)
`decisionRef` | string
`mapDecisionResult` | string
`hitPolicy` | string
`inputs` | [Array&lt;DMNInputAST&gt;](DMNInputAST.md)
`outputs` | [Array&lt;DMNOutputAST&gt;](DMNOutputAST.md)
`rules` | [Array&lt;DMNRuleAST&gt;](DMNRuleAST.md)
`attachedToRef` | string
`timeDuration` | string
`timeDate` | string
`timeCycle` | string
`cancelActivity` | boolean

## Example

```typescript
import type { NodeAST } from '@nativebpm/client'

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "name": null,
  "type": null,
  "topic": null,
  "wasmPath": null,
  "provider": null,
  "model": null,
  "prompt": null,
  "systemInstruction": null,
  "responseSchema": null,
  "temperature": null,
  "resultVar": null,
  "assignee": null,
  "candidateGroups": null,
  "dueDate": null,
  "inputSchema": null,
  "formId": null,
  "formKey": null,
  "calledElement": null,
  "inVariables": null,
  "outVariables": null,
  "decisionRef": null,
  "mapDecisionResult": null,
  "hitPolicy": null,
  "inputs": null,
  "outputs": null,
  "rules": null,
  "attachedToRef": null,
  "timeDuration": null,
  "timeDate": null,
  "timeCycle": null,
  "cancelActivity": null,
} satisfies NodeAST

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as NodeAST
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


