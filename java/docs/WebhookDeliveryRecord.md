

# WebhookDeliveryRecord


## Properties

| Name | Type | Description | Notes |
|------------ | ------------- | ------------- | -------------|
|**id** | **UUID** |  |  |
|**webhookId** | **UUID** |  |  |
|**tenantId** | **String** |  |  |
|**eventType** | **String** |  |  |
|**payload** | **byte[]** |  |  |
|**status** | **String** |  |  |
|**responseCode** | **Integer** |  |  [optional] |
|**responseBody** | **String** |  |  [optional] |
|**attempts** | **Integer** |  |  |
|**nextRetry** | **OffsetDateTime** |  |  [optional] |
|**createdAt** | **OffsetDateTime** |  |  |
|**processedAt** | **OffsetDateTime** |  |  [optional] |



