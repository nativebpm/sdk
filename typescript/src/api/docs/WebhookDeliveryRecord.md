
# WebhookDeliveryRecord


## Properties

Name | Type
------------ | -------------
`id` | string
`webhookId` | string
`tenantId` | string
`eventType` | string
`payload` | string
`status` | string
`responseCode` | number
`responseBody` | string
`attempts` | number
`nextRetry` | Date
`createdAt` | Date
`processedAt` | Date

## Example

```typescript
import type { WebhookDeliveryRecord } from '@nativebpm/client'

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "webhookId": null,
  "tenantId": null,
  "eventType": null,
  "payload": null,
  "status": PENDING,
  "responseCode": null,
  "responseBody": null,
  "attempts": null,
  "nextRetry": null,
  "createdAt": null,
  "processedAt": null,
} satisfies WebhookDeliveryRecord

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as WebhookDeliveryRecord
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


