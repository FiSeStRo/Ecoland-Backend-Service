<?php

class UserService{
    use HasDatabaseAccessTrait;

    public function getUsernameById(int $userId) : string{
        $username = '';
        $sql = "SELECT username FROM ". DbTables::Users->value ." WHERE id=?";       
        if($this->m_Db->createStatement($sql) ){
            $this->m_Db->bindStatementParamInt($userId);
            $status = $this->m_Db->executeStatement();
            if( $status->getRequestStatus() == RequestStatus::ValidDatabaseRequest){
                $data = $status->getData();
                $username = $data[0]['username'];
            }
        }
        return $username;
    }

    public function getUserIdByLogin(string $username, string $password) : int{        
        $userId = 0;
        $sql = "SELECT * FROM ". DbTables::Users->value ." WHERE username=?";       
        if($this->m_Db->createStatement($sql) ){
            $this->m_Db->bindStatementParamString($username);
            $status = $this->m_Db->executeStatement();
            if( $status->getRequestStatus() == RequestStatus::ValidDatabaseRequest){
                $user = $status->getData();
                if( !empty($user)){
                    // Password check
                    if( $password == $user['password']){
                        $userId = $user['id'];
                    }
                }
            }
        }
        return $userId;
    }

    public function doesUserWithIdExist( int $userId ) : bool{
        $sql = "SELECT id FROM ". DbTables::Users->value ." WHERE id=?";
        $status = new InternalStatus(RequestStatus::Undefined);
        if($this->m_Db->createStatement($sql) ){
            $this->m_Db->bindStatementParamInt($userId);
            $status = $this->m_Db->executeStatement(true);           
        }
        return $status->getNumAffectedRows() > 0;
    }

    public function doesUserWithUsernameExist(string $username) : bool{
        $sql = "SELECT id FROM ". DbTables::Users->value ." WHERE username=?";
        $status = new InternalStatus(RequestStatus::Undefined);
        if($this->m_Db->createStatement($sql) ){
            $this->m_Db->bindStatementParamString($username);
            $status = $this->m_Db->executeStatement(true);
        }
        return $status->getNumAffectedRows() > 0;
    }

    public function doesUserHaveUserLevel(int $userId, int $requiredUserLevel) : bool{
        $userStatus = $this->getUserData($userId);
        if( !$userStatus->isValidStatus() ){
            return false;
        }

        $userLevel = intval($userStatus->getData()['role']);
        return $userLevel >= $requiredUserLevel;
    }
    
    private function doesUserWithEmailExist(string $email) : InternalStatus{
        if( !$this->isValidEmail($email)){
            return new InternalStatus(RequestStatus::InvalidInputEmail);
        }

        $sql = "SELECT id FROM " . DbTables::Users->value ." WHERE email=?";
        $status = new InternalStatus(RequestStatus::Undefined);
        if( $this->m_Db->createStatement($sql)){
            $this->m_Db->bindStatementParamString($email);
            $status = $this->m_Db->executeStatement(true);
        }

        return new InternalStatus(($status->getNumAffectedRows() > 0) 
            ? RequestStatus::UserEmailDoesAlreadyExist
            : RequestStatus::Valid
        );
    }

    public function getUserList() : InternalStatus{
        $sql = "SELECT * FROM ". DbTables::Users->value ."";       
        if( $this->m_Db->createStatement($sql)){
            return $this->m_Db->executeStatement();           
        }

        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }   

    public function getUserResources(int $userId) : InternalStatus{
        $sql = "SELECT * FROM ". DbTables::UserResources->value ." WHERE user_id = ?";
        if( $this->m_Db->createStatement($sql)){
            $this->m_Db->bindStatementParamInt($userId);
            return $this->m_Db->executeStatement();
        }

        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }

    public function createNewUser(string $username, string $email, string $password) : InternalStatus{
        if( $this->doesUserWithUsernameExist($username) ){
            return new InternalStatus(RequestStatus::UserDoesAlreadyExist);
        }       

        if( !$this->isValidEmail($email)){
            return new InternalStatus(RequestStatus::InvalidInputEmail);
        }

        $emailStatus = $this->doesUserWithEmailExist($email);
        if( !$emailStatus->isValidStatus() ){
            return $emailStatus;
        }

        $sql = "INSERT INTO ". DbTables::Users->value ." (username, email, password, role) VALUES (?, ?, ?, ?)";
        $status = new InternalStatus(RequestStatus::Undefined);
        if( $this->m_Db->createStatement($sql)){
            $this->m_Db->bindStatementParamString($username);
            $this->m_Db->bindStatementParamString($email);
            $this->m_Db->bindStatementParamString($password);
            $this->m_Db->bindStatementParamInt(self::NEW_USER_DEFAULT_ROLE);
            $status = $this->m_Db->executeStatement();        
        }
        else{
            return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
        }

        // Initialize entry in user_resource table aswell
        if( $status->isValidStatus() ){
            $statusUserResources = new InternalStatus(RequestStatus::Undefined);
            $sql = "INSERT INTO ". DbTables::UserResources->value ." (user_id, money, prestige) VALUES (?, ?, ?)";
            if( $this->m_Db->createStatement($sql)){
                $this->m_Db->bindStatementParamInt($status->getNewId());
                $this->m_Db->bindStatementParamFloat(self::NEW_USER_START_MONEY);
                $this->m_Db->bindStatementParamInt(self::NEW_USER_START_PRESTIGE);
                $statusUserResources = $this->m_Db->executeStatement();               
            }     
            if( !$statusUserResources->isValidStatus() ){
                // User Resources entry error -> delete user entry aswell
                $this->deleteUser($status->getNewId());
                return $statusUserResources;               
            }
        }

        return $status;
    }

    private function deleteUser(int $userId){
        // TODO: check authorization for delete
        $sql = "DELETE FROM ". DbTables::Users->value ." WHERE id = ?";
        if( $this->m_Db->createStatement($sql)){
            $this->m_Db->bindStatementParamInt($userId);
            $this->m_Db->executeStatement();
        }

        $sql = "DELETE FROM ". DbTables::UserResources->value ." WHERE user_id = ?";
        if( $this->m_Db->createStatement($sql)){
            $this->m_Db->bindStatementParamInt($userId);
            $this->m_Db->executeStatement();
        }
    }

    private function isValidEmail(string $email) : bool{       
        return ( filter_var($email, FILTER_VALIDATE_EMAIL) !== false);
    }

    private function getUserData(int $userId) : InternalStatus{
        $sql = "SELECT * FROM ". DbTables::Users->value ." WHERE id=?";
        $status = new InternalStatus(RequestStatus::Undefined);
        if($this->m_Db->createStatement($sql) ){
            $this->m_Db->bindStatementParamInt($userId);
            $status = $this->m_Db->executeStatement();           
        }

        return $status;
    }

    private const NEW_USER_DEFAULT_ROLE = 0;
    private const NEW_USER_START_MONEY = 100000.0;
    private const NEW_USER_START_PRESTIGE = 0;
}

?>