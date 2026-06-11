
# TaskRecord


## Properties

Name | Type
------------ | -------------
`id` | string
`instanceId` | string
`activityId` | string
`name` | string
`assignee` | string
`candidateGroups` | string
`status` | string
`dueDate` | Date
`inputSchema` | string
`createdAt` | Date
`claimedAt` | Date
`completedAt` | Date
`currentStep` | number
`draftVariables` | { [key: string]: any; }

## Example

```typescript
import type { TaskRecord } from '@nativebpm/client'

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "instanceId": null,
  "activityId": null,
  "name": null,
  "assignee": null,
  "candidateGroups": null,
  "status": CREATED,
  "dueDate": null,
  "inputSchema": null,
  "createdAt": null,
  "claimedAt": null,
  "completedAt": null,
  "currentStep": null,
  "draftVariables": null,
} satisfies TaskRecord

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as TaskRecord
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


