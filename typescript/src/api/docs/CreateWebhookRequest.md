
# CreateWebhookRequest


## Properties

Name | Type
------------ | -------------
`url` | string
`secret` | string
`events` | Array&lt;string&gt;
`processId` | string
`isActive` | boolean
`enableAudit` | boolean

## Example

```typescript
import type { CreateWebhookRequest } from '@nativebpm/client'

// TODO: Update the object below with actual values
const example = {
  "url": null,
  "secret": null,
  "events": null,
  "processId": null,
  "isActive": null,
  "enableAudit": null,
} satisfies CreateWebhookRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as CreateWebhookRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


