<?php

class TestEndpoint extends Endpoint{
    use CommandHandlingTrait;

    public function __construct(string $command, array $params, AuthenticationHandler &$authHandler){
        parent::__construct($command, $params, $authHandler);

        $this->registerCommand('isItAlive', 'testIsItAlive');
        $this->registerCommand('isItAuthenticated', 'testIsItAuthenticated', CommandType::Get, true);
        $this->registerCommand('isGetAlive', 'testIsGetAlive', CommandType::GetWithId);
        $this->registerCommand('isPostAlive', 'testIsPostAlive', CommandType::PostFormData);
        $this->registerCommand('isDatabaseAlive', 'testIsDatabaseAlive', CommandType::Get);
    }

    // Dummy methods for get, post and authentication
    private function testIsItAlive(){}
    private function testIsItAuthenticated(){}
    private function testIsGetAlive(){}
    private function testIsPostAlive(){}

    private function testIsDatabaseAlive() : InternalStatus{
        $userService = new UserService();       
        return $userService->getUserList();
    }
}

?>