<?php

class BuildingEndpoint extends Endpoint {
    use CommandHandlingTrait;

    public function __construct(string $command, array $params, AuthenticationHandler &$authHandler) {
        parent::__construct($command, $params, $authHandler);

        $this->registerCommand('construct', 'constructBuilding', CommandType::PostJson, UserLevel::User);
        $this->registerCommand( 'details', 'getBuildingDetails', CommandType::GetWithId, UserLevel::User);
    }   

    public function constructBuilding() : InternalStatus{
        $userId = $this->getCurrentUserId();
        $userService = new UserService;
        if( !$userService->doesUserWithIdExist($userId)){           
            return new InternalStatus(RequestStatus::UserDoesNotExist);
        }

        $defId = intval($this->getParams()["def_id"] ?? 0);
        $displayName = $this->getParams()["display_name"] ?? "";
        $buildingService = new BuildingService;
        $status = $buildingService->constructBuilding($defId, $userId, $displayName);
        return $status;
    }

    public function getBuildingDetails() : InternalStatus{
        $buildingId = $this->getParams()['id'];
        $userId = $this->getCurrentUserId();

        $buildingService = new BuildingService;
        $validatationStatus = $buildingService->doesBuildingBelongToUserId($buildingId, $userId);
        if( !$validatationStatus->isValidStatus()){
            return $validatationStatus;
        }
       
        // Request all the required data from the different service and build the data for the response.
        $buildingStatus = $buildingService->getBuildingDetails($buildingId);
        if(!$buildingStatus->isValidStatus()){
            return $buildingStatus;
        }   
        $buildingData = $buildingStatus->getData();
        $buildingDefinition = null;        
        $buildingDefinitionId = $buildingData['def_id'];
        if( $buildingDefinitionId > 0){
            $buildingDefinition = $buildingService->getDefinition(DefinitionType::Building, $buildingDefinitionId);
        }

        $productionService = new ProductionService;
        $productionStatus = $productionService->getAllProductionsForBuildingId($buildingId);
        if(!$productionStatus->isValidStatus()){
            return $productionStatus;
        }

        $storageService = new StorageService;
        $storageStatus = $storageService->getStorageForBuilding($buildingId);
        if(!$storageStatus->isValidStatus()){
            return $storageStatus;
        }

        $detailStatus = new InternalStatus(RequestStatus::ValidDatabaseRequest);
        $detailData = (object)[];
        $detailData->id = $buildingData['id'];
        $detailData->name = $buildingData['name'];
        $detailData->type = (object)[];
        $detailData->type->def_id = $buildingData['def_id'];
        $detailData->type->token_name = $buildingDefinition->token_name ?? '';

        // #TODO: Implementation weiter

        //dev_var_dump($detailData);
        return new InternalStatus(RequestStatus::Undefined);
    }
}

?>