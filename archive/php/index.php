<?php

// Libraries via Composer
include_once './vendor/autoload.php';

// Helper classes
include_once './utilities/constants.php';
include_once './utilities/config.php';
include_once './utilities/internalstatus.php';
include_once './utilities/helpers.php';
include_once './utilities/database.php';
include_once './utilities/response.php';
include_once './utilities/authentication.php';
include_once './utilities/gamestate.php';
include_once './utilities/request.php';

// Services 
include_once './services/service.php';
include_once './services/finance_service.php';
include_once './services/buildings_service.php';
include_once './services/storage_service.php';
include_once './services/production_service.php';
include_once './services/product_service.php';
include_once './services/user_service.php';
include_once './services/setup_service.php';

// Endpoints
include_once './endpoints/endpoint.php';
include_once './endpoints/admin_endpoint.php';
include_once './endpoints/authentication_endpoint.php';
include_once './endpoints/buildings_endpoint.php';
include_once './endpoints/production_endpoint.php';
include_once './endpoints/test_endpoint.php';
include_once './endpoints/user_endpoint.php';

// Set environment variables
$env = file_get_contents(__DIR__."/.env");
$lines = explode("\n",$env);
foreach($lines as $line){
    // Ignore comment lines ("#") or without a value after "="
    preg_match("/([^#]+)\=(.*)/",$line,$matches);
    if(isset($matches[2])){
        putenv(trim($line));
    }
}

$request = new RequestHandler();
$request->handleRequest();

?>