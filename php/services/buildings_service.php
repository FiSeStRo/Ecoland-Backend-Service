<?php

class BuildingDefinition
{
    use DefinitionDBInitTrait;

    public readonly int $id;
    public readonly string $token_name;
    public readonly float $base_construction_cost;
    public readonly float $base_construction_time;
    private array $production_list = [];

    public function __construct(int $def_id)
    {
        $db = new DatabaseHandler;
        $sql = "SELECT * FROM " . DbTables::DefBuildings->value . " WHERE id = ?";
        if ($db->createStatement($sql)) {
            $db->bindStatementParamInt($def_id);
            $status = $db->executeStatement();

            $this->initFromDatabaseResult($status->getData(), [
                "id",
                "token_name",
                "base_construction_cost",
                "base_construction_time",
            ]);
        }

        $this->initProductionList($db);
    }

    public function getPossibleProductionList(): array
    {
        return $this->production_list;
    }

    private function initProductionList(DatabaseHandler &$db)
    {
        $sql = "SELECT production_id FROM " . DbTables::DefBuildingProduction->value . " WHERE building_id = ?";
        if ($db->createStatement($sql)) {
            $db->bindStatementParamInt($this->id);
            $status = $db->executeStatement();
            if ($status->isValidStatus()) {
                foreach ($status->getData() as $production) {
                    array_push($this->production_list, $production["production_id"]);
                }
            }
        }
    }
}

class BuildingService
{
    use HasDefinitionTrait;
    use HasDatabaseAccessTrait;

    public function constructBuilding(int $defId, int $userId, string $displayName): InternalStatus
    {
        $buildingDefinition = $this->getDefinition(DefinitionType::Building, $defId);
        if ($buildingDefinition === false) {
            return new InternalStatus(RequestStatus::BuildingDefinitionMissing);
        }

        // Check for enough user money
        $financeService = new FinanceService();
        $financeStatus = $financeService->hasSufficientFunds($userId, $buildingDefinition->base_construction_cost);
        if (!$financeStatus->isValidStatus()) {
            return $financeStatus;
        }

        // Create new building
        $sql = "INSERT INTO " . DbTables::Buildings->value . " (user_id, def_id, name) VALUES (?, ?, ?)";
        if ($this->m_Db->createStatement($sql)) {
            $this->m_Db->bindStatementParamInt($userId);
            $this->m_Db->bindStatementParamInt($defId);
            $this->m_Db->bindStatementParamString($displayName);
            $status = $this->m_Db->executeStatement();
            if ($status->isValidStatus()) {
                // Pay the money for construction               
                $statusFinance = $financeService->withdrawMoney($userId, $buildingDefinition->base_construction_cost);
                if (!$statusFinance->isValidStatus()) {
                    // if it didn't work -> remove building
                    $this->demolishBuilding($status->getNewId(), $userId);
                    return $statusFinance;
                }
                return $status;
            }
        }
        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }

    public function demolishBuilding(int $buildingId, int $userId): InternalStatus
    {
        $sql = "DELETE FROM " . DbTables::Buildings->value . " WHERE id = ? AND user_id = ?";
        if ($this->m_Db->createStatement($sql)) {
            $this->m_Db->bindStatementParamInt($buildingId);
            $this->m_Db->bindStatementParamInt($userId);
            return $this->m_Db->executeStatement();
        }
        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }

    public function getUserIdForBuilding(int $buildingId): InternalStatus
    {
        $sql = "SELECT user_id FROM " . DbTables::Buildings->value . " WHERE id = ?";
        if ($this->m_Db->createStatement($sql)) {
            $this->m_Db->bindStatementParamInt($buildingId);
            $status = $this->m_Db->executeStatement();
            if (!$status->isValidStatus()) {
                return $status;
            }
            $data = $status->getData();
            return (empty($data)) ? new InternalStatus(RequestStatus::BuildingDoesNotExist) : $status;
        }
        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }

    public function doesBuildingBelongToUserId(int $buildingId, int $userId): InternalStatus
    {
        $status = $this->getUserIdForBuilding($buildingId);
        $data = $status->getData();
        if (empty($data)) {
            return new InternalStatus(RequestStatus::BuildingDoesNotExist);
        }

        if ($status->isValidStatus()) {
            return new InternalStatus(($userId == $data['user_id']) ? RequestStatus::Valid : RequestStatus::BuildingDoesNotBelongToUser);
        }
        return $status;
    }

    public function isBuildingOnIdle(int $buildingId): InternalStatus
    {
        $sql = "SELECT * FROM " . DbTables::ManufacturingOrders->value . " WHERE building_id = ? AND is_completed = ?";
        if ($this->m_Db->createStatement($sql)) {
            $this->m_Db->bindStatementParamInt($buildingId);
            $this->m_Db->bindStatementParamBool(false);
            $status = $this->m_Db->executeStatement(true);
            if ($status->isValidStatus()) {
                return new InternalStatus(($status->getNumAffectedRows() == 0) ? RequestStatus::Valid : RequestStatus::BuildingHasUnfinishedProduction);
            }
            return $status;
        }
        return new InternalStatus(RequestStatus::DatabaseStmtCreationError);
    }

    public function canBuildingDoProduction(int $buildingId, int $productionDefinitionId): InternalStatus
    {
        $buildingDefinition = $this->getDefinitionForBuilding($buildingId);
        if ($buildingDefinition->isValid) {
            $canBuildingDoProduction = in_array($productionDefinitionId, $buildingDefinition->getPossibleProductionList());
            return new InternalStatus(($canBuildingDoProduction) ? RequestStatus::Valid : RequestStatus::BuildingCannotDoProduction);
        }
        return new InternalStatus(RequestStatus::UnknownBuildingDefinition);
    }

    private function getDefinitionForBuilding(int $buildingId): BuildingDefinition
    {
        $defId = 0;
        $sql = "SELECT def_id FROM " . DbTables::Buildings->value . " WHERE id = ?";
        if ($this->m_Db->createStatement($sql)) {
            $this->m_Db->bindStatementParamInt($buildingId);
            $status = $this->m_Db->executeStatement();
            if ($status->isValidStatus()) {
                $defId = $status->getData()["def_id"];
            }
        }

        return $this->getDefinition(DefinitionType::Building, $defId);
    }
}
