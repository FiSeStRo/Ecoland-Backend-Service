<?php

class BuildingEndpoint extends Endpoint {
    use CommandHandlingTrait;

    public function __construct(string $command, array $params, AuthenticationHandler &$authHandler) {
        parent::__construct($command, $params, $authHandler);

        $this->registerCommand('construct', 'constructBuilding', CommandType::PostJson, true);
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
}

?>