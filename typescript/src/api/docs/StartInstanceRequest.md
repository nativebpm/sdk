
# StartInstanceRequest


## Properties

Name | Type
------------ | -------------
`instanceId` | string
`businessKey` | string
`variables` | { [key: string]: any; }

## Example

```typescript
import type { StartInstanceRequest } from '@nativebpm/client'

// TODO: Update the object below with actual values
const example = {
  "instanceId": null,
  "businessKey": null,
  "variables": null,
} satisfies StartInstanceRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as StartInstanceRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


