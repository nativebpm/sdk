<?php

return [
    /*
    |--------------------------------------------------------------------------
    | NativeBPM Engine Host URL
    |--------------------------------------------------------------------------
    |
    | The base URL of the running NativeBPM engine API.
    |
    */
    'host' => env('NATIVEBPM_HOST', 'http://localhost:8080'),

    /*
    |--------------------------------------------------------------------------
    | NativeBPM API Token
    |--------------------------------------------------------------------------
    |
    | The Bearer token required for authentication against the API.
    |
    */
    'token' => env('NATIVEBPM_TOKEN', 'test-bearer-token'),
];
