<?php

namespace App\Providers;

use Illuminate\Support\ServiceProvider;
use NativeBPM\Client\Configuration;
use NativeBPM\Client\Api\DefaultApi;
use GuzzleHttp\Client as GuzzleClient;

class NativeBPMServiceProvider extends ServiceProvider
{
    /**
     * Register any application services.
     */
    public function register(): void
    {
        $this->app->singleton(DefaultApi::class, function ($app) {
            $config = new Configuration();
            $config->setHost(config('nativebpm.host'));

            // Inject the Bearer Authorization header into Guzzle
            $token = config('nativebpm.token');
            $client = new GuzzleClient([
                'headers' => [
                    'Authorization' => "Bearer {$token}",
                    'Accept'        => 'application/json',
                ],
            ]);

            return new DefaultApi($client, $config);
        });
    }

    /**
     * Bootstrap any application services.
     */
    public function boot(): void
    {
        if ($this->app->runningInConsole()) {
            $this->publishes([
                __DIR__.'/../../config/nativebpm.php' => config_path('nativebpm.php'),
            ], 'config');
        }
    }
}
