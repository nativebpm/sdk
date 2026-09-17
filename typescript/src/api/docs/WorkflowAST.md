
# WorkflowAST


## Properties

Name | Type
------------ | -------------
`id` | string
`name` | string
`inputSchema` | string
`nodes` | [Array&lt;NodeAST&gt;](NodeAST.md)
`flows` | [Array&lt;FlowAST&gt;](FlowAST.md)

## Example

```typescript
import type { WorkflowAST } from '@nativebpm/client'

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "name": null,
  "inputSchema": null,
  "nodes": null,
  "flows": null,
} satisfies WorkflowAST

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as WorkflowAST
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


