<?php

class RequestHandler{
    public function handleRequest() {
        $this->m_Response = new Response();
        $this->m_AuthHandler = new AuthenticationHandler();
        $this->m_GameState = new GameStateHandler();

        $this->parseRequest();       
    }

    private function parseRequest(){
        // Separate the path into smaller pieces.       
        $pathSections = $this->extractPathSections($_GET['route']);       

        // Extract the endpoint, command and params for the command
        $requestEndpoint = $this->extractEndpointFromPath($pathSections);
        $requestCommand = $this->extractCommandFromPath($pathSections);
        $requestParams = $this->extractParams($pathSections);

        $isEndpointValid = array_key_exists($requestEndpoint, self::ENDPOINT_LIST);
        $endpoint = null;
        $status = new InternalStatus(RequestStatus::Undefined);
        if( !$isEndpointValid || self::ENDPOINT_LIST[$requestEndpoint] == ''){
            $status = new InternalStatus(RequestStatus::UnknownEndpoint);            
        }
        else if( $requestCommand == '' ){
            $status = new InternalStatus(RequestStatus::MissingCommand);           
        }
        else{
            // Check current game state before handling the request
            $this->m_GameState->checkGameState();

            // Handle the request
            $endpointClassName = self::ENDPOINT_LIST[$requestEndpoint];
            if( $endpoint = new $endpointClassName($requestCommand, $requestParams, $this->m_AuthHandler) ){
                $status = $endpoint->handleCommand();
            }
        }

        if( isset($status) ){
            $this->m_Response->handleInternalStatus($status);
        }
        else{
            // Use default status?
            $status = new InternalStatus(RequestStatus::Undefined);
            $this->m_Response->handleInternalStatus($status);
        }
        $this->sendResult();
    }

    private function extractPathSections(string $route) : array{
        // With redirect from .htaccess this gets quite simple.
        return (isset($_GET[self::ROUTE_PARAMS_KEY])) ? explode('/', $_GET[self::ROUTE_PARAMS_KEY]) : array();        
    }

    private function extractEndpointFromPath(array $pathSections ) : string {      
        return ( $pathSections ) ? strtolower($pathSections[0]) : '';       
    }

    private function extractCommandFromPath(array $pathSections ) : string {
        $numPathSections = count($pathSections);
        return ($pathSections && $numPathSections >= 2) ? strtolower($pathSections[1]) : '';
    }

    private function extractParams(array $pathSections) : array {
        $getParams = $_GET;
        $postParams = $_POST; 
        if( isset($_SERVER['CONTENT_TYPE'])){
            if( $_SERVER['CONTENT_TYPE'] == CONTENT_TYPE_APPLICATION_JSON){
                $postJsonString = file_get_contents('php://input');
                $postJsonParams = json_decode($postJsonString);
                if( !empty($postJsonParams)){
                    $postParams = get_object_vars($postJsonParams);
                }
            }
        }

        // Store potential GET id
        if(isset($pathSections[2]) ){
            $getParams['id'] = intval($pathSections[2]);
        }

        // Remove route from params
        unset($getParams[self::ROUTE_PARAMS_KEY]);

        // Sanitize remaining parameters
        foreach( $getParams as &$param ){
            $param = preg_replace('/[^-a-zA-Z0-9_.@]/', '', $param);
        }
        foreach( $postParams as $key => &$param){           
            if( $key != 'display_name'){
                $param = preg_replace('/[^-a-zA-Z0-9_.@]/', '', $param);
            }
        }

        $params = [
            'GET' => $getParams,
            'POST' => $postParams,
        ];
        return $params;
    }

    private function sendResult(){
        $this->m_Response->send();
    }

    private const ROUTE_PARAMS_KEY = 'route';

    // valid content-type-headers (for POST requests)

    // list of endpoint strings <-> endpoint class names    
    private const ENDPOINT_LIST = array(
        'admin' => 'AdministratorEndpoint',
        'authentication' => 'AuthenticationEndpoint',
        'test' => 'TestEndpoint',
        'buildings' => 'BuildingEndpoint',
        'production' => 'ProductionEndpoint',
        'user' => 'UserEndpoint',
    );
    
    private Response $m_Response;
    private AuthenticationHandler $m_AuthHandler;
    private GameStateHandler $m_GameState;
}

?>