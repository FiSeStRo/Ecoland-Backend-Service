<?php

class ProductionDefinition{
    use DefinitionDBInitTrait;

    public readonly int $id;
    public readonly string $token_name;
    public readonly float $cost;
    public readonly float $base_duration;
    private array $products = [];

    public function __construct(int $def_id) {
        $db = new DatabaseHandler;
        $sql = "SELECT * FROM " . DbTables::DefProduction->value . " WHERE id = ?";
        if( $db->createStatement($sql) ){
            $db->bindStatementParamInt($def_id);
            $status = $db->executeStatement();
            
            $this->initFromDatabaseResult($status->getData(), [
                "id",
                "token_name",
                "cost",
                "base_duration",
            ]);
        }

        $this->initializeProductionRecipe($db);
    }
    
    public function getProducts() : array{
        return $this->products;
    }

    private function initializeProductionRecipe(DatabaseHandler &$db){
        if( !$this->isValid ){
            return;
        }

        $sql = "SELECT ". DbTables::DefProductionRecipe->value .".*, " . DbTables::DefProducts->value. ".token_name 
                FROM " . DbTables::DefProductionRecipe->value . " 
                LEFT JOIN " . DbTables::DefProducts->value . " ON " . DbTables::DefProductionRecipe->value . ".product_id="
                . DbTables::DefProducts->value . ".id WHERE production_id = ?";
        if( $db->createStatement($sql)){
            $db->bindStatementParamInt($this->id);
            $status = $db->executeStatement(onSingleResultUseArray: true);
            if( $status->isValidStatus()){
                foreach($status->getData() as $product){
                    array_push($this->products, $product);
                }
            }
        }
    }
}

class ProductionService{
    use HasDefinitionTrait;
    use HasDatabaseAccessTrait;

    public function startProduction(int $buildingId, int $userId, int $productionDefId, int $numCycles) : InternalStatus{
        // Validate input data
        $buildingService = new BuildingService;

        if( $numCycles <= 0){
            return new InternalStatus(RequestStatus::InvalidInput);
        }

        // Is there currently an ongoing production at this building?       
        $buildingProductionStatus = $buildingService->isBuildingOnIdle($buildingId);
        if( !$buildingProductionStatus->isItTrue()){
            return $buildingProductionStatus;
        }

        // Does building belong to user?
        $status = $buildingService->doesBuildingBelongToUserId($buildingId, $userId);
        if( !$status->isItTrue() ){
            return $status;
        }

        // Does production definition exist?
        $productionDefinition = $this->getDefinition(DefinitionType::Production, $productionDefId);        
        if( !$productionDefinition->isValid){
            return new InternalStatus(RequestStatus::UnknownProductionDefinition);
        }

        // Can the building execute the production?
        $status = $buildingService->canBuildingDoProduction($buildingId, $productionDefId);
        if( !$status->isItTrue() ){
            return $status;
        }

        // Does the building has all the required wares?
        $storageService = new StorageService;
        $status = $storageService->checkStorageForProduction($buildingId, $productionDefinition, $numCycles);
        if( !$status->isValidStatus()){           
            return $status;
        }

        // Does the user have enough money?
        $costPerCycle = $productionDefinition->cost;
        $totalCost = $costPerCycle * $numCycles;
        $financeService = new FinanceService;
        $status = $financeService->hasSufficientFunds($userId, $totalCost);
        if( !$status->isValidStatus()){
            return $status;
        }
        
        // Every check passed -> save production order in database
        return $this->createProductionOrder($buildingId, $userId, $productionDefinition, $numCycles);
    }

    public function cancelAllActiveProductionOrders() : InternalStatus{
        // TODO: Implementation
        return new InternalStatus(RequestStatus::Undefined);
    }

    public function cancelProductionOrder(int $productionOrderId) : InternalStatus{
        
        return new InternalStatus(RequestStatus::Undefined);
    }

    public function finishProductionOrder(array &$order) : InternalStatus{

        $timeEnd = strtotime($order["time_end"]);
        if( $timeEnd === false || $timeEnd > time()){
            return new InternalStatus(RequestStatus::ProductionInvalidTimestamp);
        }
        
        if( $order["is_completed"] ){
            return new InternalStatus(RequestStatus::ProductionAlreadyCompleted);
        }
        
        // Get the production definition for the output items
        $productionDefinition = $this->getDefinition(DefinitionType::Production, $order["production_id"]);
        $buildingService = new BuildingService;
        $storageService = new StorageService;
        $financeService = new FinanceService;

        $buildingId = $order["building_id"];              
        $numCycles = $order["cycles"];
        $userIdStatus = $buildingService->getUserIdForBuilding($buildingId);
        if( !$userIdStatus->isValidStatus()){
            return $userIdStatus;
        }
        $userId = $userIdStatus->getData()["user_id"] ?? 0;
        if( $userId == 0 ){
            return new InternalStatus(RequestStatus::UserDoesNotExist);
        }       

        if( $productionDefinition->isValid ){
            // Add all the produced (= output) items to the storage
            $productStatus = new InternalStatus(RequestStatus::Valid);
            $productList = $productionDefinition->getProducts();
            foreach($productList as $product){               
                if( !$product["is_input"]){
                    $totalAmount = $product["amount"] * $numCycles;
                    if( $product["product_id"] == PRODUCT_ID_MONEY){   
                        // Selling Products -> deposit money                    
                        $financeStatus = $financeService->despositMoney($userId, floatval($totalAmount));
                        if( !$financeStatus->isValidStatus()){
                            $productStatus = $financeStatus;
                        }
                    }
                    else{
                        // Manufactoring Products -> store in storage
                        $storageStatus = $storageService->addProductToStorage($buildingId, $product["product_id"], $totalAmount);
                        if( !$storageStatus->isValidStatus()){
                            $productStatus = $storageStatus;
                        }
                    }
                }
            }

            if(!$productStatus->isValidStatus()){
                return $productStatus;
            }
        }

        // Mark the production order as complete
        $sql = "UPDATE " . DbTables::ManufacturingOrders->value . " SET is_completed = 1 WHERE id = ?";
        if( $this->m_Db->createStatement($sql) ){
            $this->m_Db->bindStatementParamInt($order["id"]);
            return $this->m_Db->executeStatement();
        }

        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }

    private function createProductionOrder(int $buildingId, int $userId, ProductionDefinition &$productionDefinition, int $numCycles) : InternalStatus{
        // Remove all products required from storage
        $storageService = new StorageService;
        foreach( $productionDefinition->getProducts() as $product){
            if( $product["is_input"]){
                $amountToRemove = $numCycles * $product["amount"];
                $storageStatus = $storageService->removeProductFromStorage($buildingId, $product["product_id"], $amountToRemove);
                if( !$storageStatus->isValidStatus()){
                    // TODO: products removed earlier are not restored
                    return $storageStatus;
                }
            }
        }

        $productionCost = $numCycles * $productionDefinition->cost;
        $financeService = new FinanceService;
        $status = $financeService->withdrawMoney($userId, $productionCost);
        if( !$status->isValidStatus()){
            return $status;
        }

        $sql = "INSERT INTO " . DbTables::ManufacturingOrders->value . " (building_id, production_id, time_end, cycles) VALUES (?, ?, ?, ?)";
        if( $this->m_Db->createStatement($sql)){
            $this->m_Db->bindStatementParamInt($buildingId);
            $this->m_Db->bindStatementParamInt($productionDefinition->id);
            $duration = $numCycles * $productionDefinition->base_duration;
            $time_end = time() + $duration;
            $this->m_Db->bindStatementParamTimestamp($time_end);
            $this->m_Db->bindStatementParamInt($numCycles);
            return $this->m_Db->executeStatement();
        }

        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }   
}

?>