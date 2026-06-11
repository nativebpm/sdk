[← Back to Platform Root](../../README.md)

# NativeBPM PHP SDK

Official PHP Client SDK for NativeBPM Cloud-Native engine. Uses `Guzzle` for PSR-18 compliant HTTP transport.

## Installation

Add the local package dependency to your PHP project's `composer.json` or fetch it via composer:

```bash
composer install
```

## Composer Configuration

Since this SDK is structured in a subdirectory (`sdk/php/`), it cannot be published to standard package registries directly. Instead, you can import it securely using Composer's `package` repository type pointing directly to the Git repository tag:

```json
{
    "repositories": [
        {
            "type": "package",
            "package": {
                "name": "nativebpm/client",
                "version": "1.0.0",
                "source": {
                    "url": "https://gitlab.com/nativebpm/platform.git",
                    "type": "git",
                    "reference": "sdk/php/v1.0.0"
                },
                "autoload": {
                    "psr-4": { "NativeBPM\\Client\\": "sdk/php/lib/" }
                }
            }
        }
    ],
    "require": {
        "nativebpm/client": "1.0.0"
    }
}
```

If the GitLab repository is private, you will also need to configure credentials in `auth.json`:
```json
{
    "gitlab-token": {
        "gitlab.com": "your_personal_access_or_deploy_token"
    }
}
```

## Usage Example

Initialize configuration, add API token security header, and query resources:

```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

use NativeBPM\Client\Configuration;
use NativeBPM\Client\Api\DefaultApi;
use NativeBPM\Client\Model\StartInstanceRequest;

// 1. Configure API Client
$config = Configuration::getDefaultConfiguration()
    ->setHost('http://localhost:8080')
    ->setApiKey('Authorization', 'your-api-token')
    ->setApiKeyPrefix('Authorization', 'Bearer');

$apiInstance = new DefaultApi(
    new GuzzleHttp\Client(),
    $config
);

try {
    // 2. List deployed process definitions
    $definitions = $apiInstance->listDefinitions();
    foreach ($definitions as $def) {
        echo "Definition: " . $def->getName() . " (ID: " . $def->getId() . ")\n";
    }

    // 3. Start a process instance with variables
    $startRequest = new StartInstanceRequest([
        'business_key' => 'REQ-4402',
        'variables' => [
            'approved' => true,
            'amount' => 5000
        ]
    ]);

    $instance = $apiInstance->startInstance('expense-approval', $startRequest);
    echo "Started process instance: " . $instance->getId() . "\n";

} catch (Exception $e) {
    echo 'Exception when calling API: ', $e->getMessage(), PHP_EOL;
}
?>
```

## Developer Guide: Releases & Tagging

Since this SDK is located in a subdirectory (`sdk/php/`), it is not published to any package registry. Instead, simply tag the repository to release new versions:

```bash
git tag sdk/php/v1.0.0
git push origin sdk/php/v1.0.0
```

Consumers can then pull this specific version directly from Git using Composer's `package` repository configuration mapping the tag reference.

