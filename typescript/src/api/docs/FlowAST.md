
# FlowAST


## Properties

Name | Type
------------ | -------------
`id` | string
`source` | string
`target` | string
`condition` | string

## Example

```typescript
import type { FlowAST } from '@nativebpm/client'

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "source": null,
  "target": null,
  "condition": null,
} satisfies FlowAST

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as FlowAST
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


