<?php

class StorageService
{
    use HasDatabaseAccessTrait;

    public function checkStorageForProduction(int $buildingId, ProductionDefinition &$productionDef, int $numCycles): InternalStatus{
        // Get storage information for building
        $storageStatus = $this->getStorageForBuilding($buildingId);
        if ($storageStatus->isValidStatus()) {
            $storage = $storageStatus->getData();
            $productList = $productionDef->getProducts();
            foreach($productList as $product){               
                $totalAmountRequired = $numCycles * $product['amount'];
                $status = $this->checkStorageForProduct($storage, ( $product['is_input'] == false), $product['product_id'], $totalAmountRequired );
                if( !$status->isValidStatus()){
                    return $status;                   
                }
            }

        }
        return $storageStatus;
    }
    
    public function addProductToStorage(int $buildingId, int $productId, int $amountToAdd ) : InternalStatus{
        if( $amountToAdd <= 0){
            return new InternalStatus(RequestStatus::StorageInvalidAmount);
        }
        $storageStatus = $this->getStorageForBuilding($buildingId);
        if( $storageStatus->isValidStatus()){
            $storage = $storageStatus->getData();
            $status = $this->checkStorageForProduct($storage, true, $productId, $amountToAdd);
            if( $status->isValidStatus()){
                if($this->hasStorageEntryForProduct($storage, $productId)){
                    // Product was stored before -> edit entry
                    $sql = "UPDATE " . DbTables::Storage->value . " SET amount = amount + ? WHERE building_id = ? AND product_id = ?";                   
                    if( $this->m_Db->createStatement($sql)){
                        $this->m_Db->bindStatementParamInt($amountToAdd);
                        $this->m_Db->bindStatementParamInt($buildingId);
                        $this->m_Db->bindStatementParamInt($productId);
                        return $this->m_Db->executeStatement();
                    }
                }
                else{
                    // Product was not stored before -> new entry
                    $sql = "INSERT INTO " . DbTables::Storage->value . " (building_id, product_id, amount) VALUES (?,?,?)";
                    if( $this->m_Db->createStatement($sql)){
                        $this->m_Db->bindStatementParamInt($buildingId);
                        $this->m_Db->bindStatementParamInt($productId);
                        $this->m_Db->bindStatementParamInt($amountToAdd);
                        return $this->m_Db->executeStatement();
                    }
                    return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
                }
            }
            return $status;
        }
        return $storageStatus;
    }

    public function removeProductFromStorage(int $buildingId, int $productId, int $amountToRemove) : InternalStatus{
        if( $amountToRemove <= 0){
            return new InternalStatus(RequestStatus::StorageInvalidAmount);
        }
        $storageStatus = $this->getStorageForBuilding($buildingId);
        if( $storageStatus->isValidStatus()){
            $storage = $storageStatus->getData();
            $status = $this->checkStorageForProduct($storage, false, $productId, $amountToRemove);
            if( $status->isValidStatus()){
                if($this->hasStorageEntryForProduct($storage, $productId)){
                    // Product was stored before -> edit entry
                    $sql = "UPDATE " . DbTables::Storage->value . " SET amount = amount - ? WHERE building_id = ? AND product_id = ?";                   
                    if( $this->m_Db->createStatement($sql)){
                        $this->m_Db->bindStatementParamInt($amountToRemove);
                        $this->m_Db->bindStatementParamInt($buildingId);
                        $this->m_Db->bindStatementParamInt($productId);
                        return $this->m_Db->executeStatement();
                    }
                }
                return new InternalStatus(RequestStatus::StorageMissingProducts);
            }
            return $status;
        }
        return $storageStatus;
    }

    private function getStorageForBuilding(int $buildingId): InternalStatus{
        $sql = "SELECT * FROM " . DbTables::Storage->value . " WHERE building_id = ?";
        if ($this->m_Db->createStatement($sql)) {
            $this->m_Db->bindStatementParamInt($buildingId);
            return $this->m_Db->executeStatement();
        }
        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }

    private function checkStorageForProduct( mixed &$storage, bool $isProductAddedToStorage, int $productId, int $amount) : InternalStatus{
        $amountFound = 0;
        $storageCapacity = STORAGE_DEFAULT_CAPACITY;

        $storedProduct = $this->getStorageEntryForProduct($storage, $productId);
        if( is_array($storedProduct)){
            $amountFound = $storedProduct['amount'];
            $storageCapacity = $storedProduct['capacity'];
        }

        $status = RequestStatus::Valid;
        // Check if enough storage or products is available
        if( $isProductAddedToStorage ){
            // product to store -> check Storage capacity
            if( $amount + $amountFound > $storageCapacity ){
                $status = RequestStatus::StorageMissingCapacity;
            }
        }
        else{
            // product to remove -> check available product amount
            if( $amount > $amountFound){
                $status = RequestStatus::StorageMissingProducts;
            }
        }

        return new InternalStatus($status);
    }

    private function hasStorageEntryForProduct(mixed &$storage, int $productId) : bool{
        foreach($storage as $storedProduct){
            if( $storedProduct['product_id'] == $productId){
                return true;
            }
        }
        return false;
    }

    private function getStorageEntryForProduct(mixed &$storage, int $productId) : array|bool{
        foreach($storage as $storedProduct){
            if( $storedProduct['product_id'] == $productId){
                return $storedProduct;
            }
        }
        return false;
    }
}
