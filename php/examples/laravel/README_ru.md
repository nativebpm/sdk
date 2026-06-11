[← Назад к PHP SDK](../../README_ru.md) | [← Назад к корню платформы](../../../../../README_ru.md)

# Пример интеграции NativeBPM PHP (Laravel) SDK

Этот пример демонстрирует, как интегрировать сгенерированный клиент PHP SDK в приложение Laravel.

## 1. Загрузка пакета SDK из GitLab тегов

Поскольку сгенерированный PHP SDK находится внутри подпапки (`sdk/php/`) репозитория `platform`, стандартное обнаружение Composer VCS автоматически его не найдет. Чтобы решить эту проблему, в данном примере указывается репозиторий типа `"package"` в `composer.json`, который ссылается на Git-репозиторий `platform` по определенному тегу (например, `sdk/php/v1.0.0`) и настраивает автозагрузку файлов из нужной подпапки.

Вот выдержка из `composer.json`:

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

Просто выполните `composer install` для загрузки зависимостей.

## 2. Конфигурация (`config/nativebpm.php`)

Опубликуйте или скопируйте `config/nativebpm.php` в папку конфигурации вашего приложения:

```php
<?php

return [
    'host'  => env('NATIVEBPM_HOST', 'http://localhost:8080'),
    'token' => env('NATIVEBPM_TOKEN', 'test-bearer-token'),
];
```

Обязательно определите эти переменные среды в вашем файле `.env`:

```env
NATIVEBPM_HOST=http://localhost:8080
NATIVEBPM_TOKEN=your-nativebpm-api-token
```

## 3. Регистрация сервис-провайдера

Необходимо зарегистрировать сервис-провайдер в вашем приложении Laravel (например, внутри массива `providers` в `config/app.php` или в `bootstrap/providers.php` в зависимости от версии Laravel):

```php
App\Providers\NativeBPMServiceProvider::class,
```

`NativeBPMServiceProvider` регистрирует `DefaultApi` в сервис-контейнере Laravel, настраивая хост и создавая кастомный HTTP-клиент Guzzle, который автоматически добавляет заголовок авторизации `Authorization: Bearer <API_TOKEN>` ко всем запросам:

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

## 4. Использование в контроллере

Внедряйте `NativeBPM\Client\Api\DefaultApi` напрямую в конструктор или методы ваших контроллеров:

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

Пример `WorkflowController` предоставляет роуты и методы для:
- **Получения списка определений процессов**: `GET /api/workflows`
- **Развертывания BPMN файла**: `POST /api/workflows/deploy` (загрузка файла через форму с ключом `bpmn_file`)
- **Запуска инстанса процесса**: `POST /api/workflows/{definitionId}/start` (с опциональным параметром JSON `custom_value`)
