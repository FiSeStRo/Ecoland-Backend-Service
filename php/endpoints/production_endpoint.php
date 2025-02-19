<?php

class ProductionEndpoint extends Endpoint{
    use CommandHandlingTrait;

    public function __construct(string $command, array $params, AuthenticationHandler &$authHandler) {
        parent::__construct($command, $params, $authHandler);

        $this->registerCommand('start', 'startProduction', CommandType::PostJson, true);       
        $this->registerCommand('cancel', 'cancelProduction', CommandType::GetWithId, true);
    }

    private function startProduction() : InternalStatus{
        $params = $this->getParams(); 
        
        $buildingId = intval($params["building_id"]) ?? 0;
        $productionDefId = intval($params["id"]) ?? 0;
        $numCycles = intval($params["cycles"]) ?? 0;
        $currentUserId = $this->getCurrentUserId();
        
        $productionService = new ProductionService;
        return $productionService->startProduction($buildingId, $currentUserId, $productionDefId, $numCycles);       
    }

    private function cancelProduction() : InternalStatus{
        $params = $this->getParams();
        $productionOrderId = intval($params['id']);

        if( $productionOrderId == 0){
            return new InternalStatus(RequestStatus::ProductionDoesNotExist);
        }
        $userId = $this->getCurrentUserId();
        $productionService = new ProductionService();

        // Check if the production is already done
        $productionStatus = $productionService->isProductionActive($productionOrderId);
        if( !$productionStatus->isValidStatus() ){
            return $productionStatus;
        }

        // Check if the production belongs to a building of the current user
        $status = $productionService->doesProductionBelongToUserId($productionOrderId, $userId);
        if(!$status->isValidStatus()){
            return $status;
        }

        // Everything should be valid -> cancel production.
        return $productionService->cancelProductionOrder( $productionOrderId );
    
    }
}

?>