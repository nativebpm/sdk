
# ExecuteProcessResponse


## Properties

Name | Type
------------ | -------------
`instanceId` | string
`definitionId` | string
`version` | number
`isNewVersionDeployed` | boolean
`status` | string
`state` | { [key: string]: any; }
`currentTasks` | [Array&lt;TaskRecord&gt;](TaskRecord.md)

## Example

```typescript
import type { ExecuteProcessResponse } from '@nativebpm/client'

// TODO: Update the object below with actual values
const example = {
  "instanceId": null,
  "definitionId": null,
  "version": null,
  "isNewVersionDeployed": null,
  "status": null,
  "state": null,
  "currentTasks": null,
} satisfies ExecuteProcessResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ExecuteProcessResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


