<?php

include_once './services/service.php';

class GameStateHandler{
    use HasDatabaseAccessTrait;
    use HasDefinitionTrait;   

    public function checkGameState() : InternalStatus{
        return $this->checkFinishedManufactoringOrders();
    }

    private function checkFinishedManufactoringOrders() : InternalStatus{
        $productionService = new ProductionService;
        $sql = "SELECT * FROM " . DbTables::ManufacturingOrders->value . " WHERE time_end < ? AND is_completed = 0";
        if( $this->m_Db->createStatement($sql) ){
            $currentTime = time();
            $this->m_Db->bindStatementParamTimestamp($currentTime);
            $status = $this->m_Db->executeStatement(onSingleResultUseArray: true);
            if( $status->isValidStatus()){
                $finishedManufactoringOrders = $status->getData();                
                foreach($finishedManufactoringOrders as $order){
                    $productionService->finishProductionOrder($order);
                }
            }
        }
        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }
}

?>