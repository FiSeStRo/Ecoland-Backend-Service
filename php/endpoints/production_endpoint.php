<?php

class ProductionEndpoint extends Endpoint{
    use CommandHandlingTrait;

    public function __construct(string $command, array $params, AuthenticationHandler &$authHandler) {
        parent::__construct($command, $params, $authHandler);

        $this->registerCommand('start', 'startProduction', CommandType::PostJson, true);       
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
}

?>