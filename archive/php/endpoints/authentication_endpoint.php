<?php

class AuthenticationEndpoint extends Endpoint{
    use CommandHandlingTrait;

    public function __construct(string $command, array $params, AuthenticationHandler &$authHandler)
    {
        parent::__construct($command, $params, $authHandler);

        $this->m_UserService = new UserService();

        $this->registerCommand('sign-up', 'signUp', CommandType::PostJson);
        $this->registerCommand('sign-in', 'signIn', CommandType::PostJson);
        $this->registerCommand('refresh-token', 'refreshToken',CommandType::PostJson, UserLevel::User);
    }

    public function signUp() : InternalStatus{
        $signUpParams = $this->getParams();

        if( !isset($signUpParams['username']) || !isset($signUpParams['email']) || !isset($signUpParams['password'])
        || empty($signUpParams['username']) || empty($signUpParams['email']) || empty($signUpParams['password']) ){
            $status = new InternalStatus(RequestStatus::InvalidInput);           
            return $status;
        }

        $username = $signUpParams['username'];
        $email = $signUpParams['email'];
        $passwordRaw = $signUpParams['password'];

        // Params seem to be valid -> try to create a new account.
        $statusNewUser = $this->m_UserService->createNewUser($username, $email, $passwordRaw);

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
            $authStatus->setAuthTokens($this->createAuthToken($userId, true), $this->createAuthToken($userId, false));
            return $authStatus;
        }
    }

    public function refreshToken() : InternalStatus{
        $userId = $this->getCurrentUserId();
        
        if( $userId > 0){
            $authStatus = new InternalStatus(RequestStatus::ValidAuthentification);
            $authStatus->setAuthTokens($this->createAuthToken($userId, true), $this->createAuthToken($userId, false));
            return $authStatus;
        }
        return new InternalStatus(RequestStatus::Undefined);
    }
    
    private UserService $m_UserService;
}
?>