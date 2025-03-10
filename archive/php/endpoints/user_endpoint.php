<?php

class UserEndpoint extends Endpoint{
    use CommandHandlingTrait;
    
    public function __construct(string $command, array $params, AuthenticationHandler &$authHandler) {
        parent::__construct($command, $params, $authHandler);
    }
}

?>