<?php

class AuthenticationEndpoint extends Endpoint{
    use CommandHandlingTrait;

    public function __construct(string $command, array $params, AuthenticationHandler &$authHandler)
    {
        parent::__construct($command, $params, $authHandler);

        $this->m_UserService = new UserService();

        $this->registerCommand('sign-up', 'signUp', CommandType::PostJson);
        $this->registerCommand('sign-in', 'signIn', CommandType::PostJson);
    }

    public function signUp() : InternalStatus{
        $signUpParams = $this->getParams();

        if( !isset($signUpParams['username']) || !isset($signUpParams['password'])
        || empty($signUpParams['username']) || empty($signUpParams['password']) ){
            $status = new InternalStatus(RequestStatus::InvalidInput);           
            return $status;
        }

        $username = $signUpParams['username'];
        $passwordRaw = $signUpParams['password'];

        // Params seem to be valid -> try to create a new account.
        $statusNewUser = $this->m_UserService->createNewUser($username, $passwordRaw);

        return $statusNewUser;
    }

    public function signIn() : InternalStatus{
        $signUpParams = $this->getParams();
        if( !isset($signUpParams['username']) || !isset($signUpParams['password'])
        || empty($signUpParams['username']) || empty($signUpParams['password']) ){
            return new InternalStatus(RequestStatus::InvalidInput);           
        }

        $username = $signUpParams['username'];
        $passwordRaw = $signUpParams['password'];
        $userId = $this->m_UserService->getUserIdByLogin($username, $passwordRaw);
        if( $userId == 0 ){
            return new InternalStatus(RequestStatus::AuthenticationInvalid);
        }
        else{
            // Create Authentication Token and save in Internal Status
            $authStatus = new InternalStatus(RequestStatus::ValidAuthentification);
            $authStatus->setAuthTokens($this->createAuthToken($userId));
            return $authStatus;
        }
    }
    
    private UserService $m_UserService;
}
?>