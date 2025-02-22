<?php

class TestEndpoint extends Endpoint{
    use CommandHandlingTrait;

    public function __construct(string $command, array $params, AuthenticationHandler &$authHandler){
        parent::__construct($command, $params, $authHandler);

        $this->registerCommand('isItAlive', 'testIsItAlive');
        $this->registerCommand('isItAuthenticated', 'testIsItAuthenticated', CommandType::Get, UserLevel::User);
        $this->registerCommand('isGetAlive', 'testIsGetAlive', CommandType::GetWithParams);
        $this->registerCommand('isPostAlive', 'testIsPostAlive', CommandType::PostFormData);
        $this->registerCommand('isDatabaseAlive', 'testIsDatabaseAlive', CommandType::Get);
        $this->registerCommand('isUserLevelCheckAlive', 'testUserLevel', CommandType::PostJson);
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

    private function testUserLevel() : InternalStatus{
        $params = $this->getParams();     
        $requestStatus = RequestStatus::Undefined;

        if( isset($params['user_id']) && isset($params['role'])){
            $userService = new UserService();            
            $requestStatus = RequestStatus::CommandInsufficientUserLevel;
            if( $userService->doesUserHaveUserLevel($params['user_id'], $params['role'])){
                $requestStatus = RequestStatus::Valid;
            }
        }

        return new InternalStatus($requestStatus);
    }
}

?>