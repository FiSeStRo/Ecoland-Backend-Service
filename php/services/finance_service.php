<?php

class FinanceService{   
    use HasDatabaseAccessTrait;
    
    public function despositMoney(int $userId, float $sum) : InternalStatus{
        $userService = new UserService;
        if( !$userService->doesUserWithIdExist($userId)){
            return new InternalStatus(RequestStatus::UserDoesNotExist);
        }

        $sql = "UPDATE " . DbTables::UserResources->value ." SET money = money + ? WHERE user_id = ?";
        if( $this->m_Db->createStatement($sql)){
            $this->m_Db->bindStatementParamFloat($sum);
            $this->m_Db->bindStatementParamInt($userId);
            return $this->m_Db->executeStatement();            
        }
        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }

    public function withdrawMoney(int $userId, float $sum) : InternalStatus{
        $userService = new UserService;
        if( !$userService->doesUserWithIdExist($userId)){
            return new InternalStatus(RequestStatus::UserDoesNotExist);
        }

        $sql = "UPDATE " . DbTables::UserResources->value ." SET money = money - ? WHERE user_id = ?";
        if( $this->m_Db->createStatement($sql)){
            $this->m_Db->bindStatementParamFloat($sum);
            $this->m_Db->bindStatementParamInt($userId);
            return $this->m_Db->executeStatement();            
        }
        
        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }

    public function hasSufficientFunds(int $userId, float $sum) : InternalStatus{
        $fundStatus = $this->getCurrentFunds($userId, $sum);
        if( $fundStatus->isValidStatus()){
            $status = ($fundStatus->getData()['money'] >= $sum) ? RequestStatus::Valid : RequestStatus::UserInsufficientFunds;
            return new InternalStatus($status);
        }
        return $fundStatus;
    }

    public function getCurrentFunds(int $userId) : InternalStatus{
        $userService = new UserService();
        $userResStatus = $userService->getUserResources($userId);
        if( $userResStatus->isValidStatus()){
            $status = new InternalStatus(RequestStatus::ValidDatabaseRequest);
            $userResData = $userResStatus->getData();                      
            $currentFundsData['user_id'] = $userResData['user_id'];
            $currentFundsData['money'] = $userResData['money'];
            $status->setData($currentFundsData);
            return $status;
        }
        return $userResStatus;
    }
}

?>