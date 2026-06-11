
# IncidentRecord


## Properties

Name | Type
------------ | -------------
`id` | string
`instanceId` | string
`activityId` | string
`errorMessage` | string
`stackTrace` | string
`attemptsMade` | number
`resolved` | boolean
`createdAt` | Date
`resolvedAt` | Date

## Example

```typescript
import type { IncidentRecord } from '@nativebpm/client'

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "instanceId": null,
  "activityId": null,
  "errorMessage": null,
  "stackTrace": null,
  "attemptsMade": null,
  "resolved": null,
  "createdAt": null,
  "resolvedAt": null,
} satisfies IncidentRecord

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as IncidentRecord
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


