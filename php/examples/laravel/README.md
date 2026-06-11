[← Back to PHP SDK](../../README.md) | [← Back to Platform Root](../../../../../README.md)

# NativeBPM PHP (Laravel) SDK Integration Example

This example demonstrates how to integrate the generated PHP SDK client into a Laravel application.

## 1. Pulling the SDK Package from GitLab Tags

Since the generated PHP SDK resides inside a subdirectory (`sdk/php/`) of the `platform` repository, standard Composer VCS discovery won't automatically find it. To resolve this, this example specifies a `"package"` type repository in `composer.json` that references the `platform` Git repository at a specific tag reference (e.g. `sdk/php/v1.0.0`) and defines the subdirectory path autoloading.

Here is the excerpt from `composer.json`:

```json
{
    "repositories": [
        {
            "type": "package",
            "package": {
                "name": "nativebpm/client",
                "version": "1.0.0",
                "source": {
                    "url": "https://gitlab.com/nativebpm/sdk.git",
                    "type": "git",
                    "reference": "sdk/php/v1.0.0"
                },
                "autoload": {
                    "psr-4": {
                        "NativeBPM\\Client\\": "sdk/php/lib/"
                    }
                },
                "require": {
                    "php": "^8.1",
                    "ext-curl": "*",
                    "ext-json": "*",
                    "ext-mbstring": "*",
                    "guzzlehttp/guzzle": "^7.3",
                    "guzzlehttp/psr7": "^1.7 || ^2.0"
                }
            }
        }
    ],
    "require": {
        "nativebpm/client": "1.0.0"
    }
}
```

Simply run `composer install` to pull the dependencies.

## 2. Configuration (`config/nativebpm.php`)

Publish or copy `config/nativebpm.php` to your configuration folder:

```php
<?php

return [
    'host'  => env('NATIVEBPM_HOST', 'http://localhost:8080'),
    'token' => env('NATIVEBPM_TOKEN', 'test-bearer-token'),
];
```

Make sure to define these environment variables in your `.env` file:

```env
NATIVEBPM_HOST=http://localhost:8080
NATIVEBPM_TOKEN=your-nativebpm-api-token
```

## 3. Service Provider Registration

The package requires registering a service provider in your Laravel application (e.g. inside `config/app.php` providers array or inside `bootstrap/providers.php` depending on Laravel version):

```php
App\Providers\NativeBPMServiceProvider::class,
```

`NativeBPMServiceProvider` resolves `DefaultApi` into the Laravel container, setting up host configuration and injecting a custom Guzzle HTTP client that automatically includes the `Authorization: Bearer <API_TOKEN>` request header:

```php
$this->app->singleton(DefaultApi::class, function ($app) {
    $config = new Configuration();
    $config->setHost(config('nativebpm.host'));

    $client = new GuzzleClient([
        'headers' => [
            'Authorization' => "Bearer " . config('nativebpm.token'),
            'Accept'        => 'application/json',
        ],
    ]);

    return new DefaultApi($client, $config);
});
```

## 4. Usage in Controller

Inject the `NativeBPM\Client\Api\DefaultApi` directly into your controller:

```php
use NativeBPM\Client\Api\DefaultApi;

class WorkflowController extends Controller
{
    protected DefaultApi $api;

    public function __construct(DefaultApi $api)
    {
        $this->api = $api;
    }
    
    // ...
}
```

The example `WorkflowController` provides routes/methods to:
- **List definitions**: `GET /api/workflows`
- **Deploy BPMN file**: `POST /api/workflows/deploy` (multipart file upload with key `bpmn_file`)
- **Start instance**: `POST /api/workflows/{definitionId}/start` (with optional JSON key `custom_value`)
