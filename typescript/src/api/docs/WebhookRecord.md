
# WebhookRecord


## Properties

Name | Type
------------ | -------------
`id` | string
`tenantId` | string
`url` | string
`secret` | string
`events` | Array&lt;string&gt;
`processId` | string
`isActive` | boolean
`enableAudit` | boolean
`status` | string
`createdAt` | Date

## Example

```typescript
import type { WebhookRecord } from '@nativebpm/client'

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "tenantId": null,
  "url": null,
  "secret": null,
  "events": null,
  "processId": null,
  "isActive": null,
  "enableAudit": null,
  "status": null,
  "createdAt": null,
} satisfies WebhookRecord

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as WebhookRecord
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


