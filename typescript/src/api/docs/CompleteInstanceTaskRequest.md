
# CompleteInstanceTaskRequest


## Properties

Name | Type
------------ | -------------
`nodeId` | string
`variables` | { [key: string]: any; }

## Example

```typescript
import type { CompleteInstanceTaskRequest } from '@nativebpm/client'

// TODO: Update the object below with actual values
const example = {
  "nodeId": null,
  "variables": null,
} satisfies CompleteInstanceTaskRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CompleteInstanceTaskRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


