<?php

enum StmtStatus
{
    case Inactive;
    case InPreparation;
    case PreparationComplete;
    case InvalidParamBinding; // mapping params to placeholders was invalid or incomplete.
    case ExecutionSuccess;
    case ExecutionError;
    case ExecutionDisabled; // just mocking data, no real sql executions
}

enum StmtDataType
{
    case Int;
    case Bool;
    case Float;
    case String;
}

enum DbTables: string
{
    case Users = "users";
    case UserResources = "user_resources";
    case Buildings = "buildings";
    case Storage = "rel_building_product";
    case ManufacturingOrders = "rel_building_def_production";
    case DefBuildings = "def_buildings";
    case DefProducts = "def_products";
    case DefProduction = "def_production";
    case DefProductionRecipe = "def_rel_production_product";
    case DefBuildingProduction = "def_rel_building_production";
}

trait HasDatabaseAccessTrait
{
    public function __construct()
    {
        $this->m_Db = new DatabaseHandler;
    }

    private DatabaseHandler $m_Db;
}

class DatabaseHandler
{
    public function __construct()
    {
        mysqli_report(MYSQLI_REPORT_ERROR | MYSQLI_REPORT_STRICT);
        $this->initializeDatabase();
    }

    public function createTables(SetupService &$setupService): bool
    {
        $createQueryList = $setupService->getCreateTableQueries();
        foreach ($createQueryList as $sql) {
            $this->executeQuery($sql);
        }
        return true;
    }

    public function deleteTables(SetupService &$setupService): bool
    {
        $sql = $setupService->getDeleteTableQuery();
        $this->executeQuery($sql);
        return true;
    }

    public function createStatement(string $sql): bool
    {
        $this->resetCurrentStatement();

        if (!empty($sql)) {
            try{
                $this->m_CurrentStatement = $this->m_DbHandle->prepare($sql);
                if ($this->m_CurrentStatement !== false) {
                    $this->m_StatementStatus = StmtStatus::InPreparation;
                }
            }
            catch(mysqli_sql_exception $e){}
        }

        return ($this->m_StatementStatus == StmtStatus::InPreparation);
    }

    public function bindStatementParamInt(int $value)
    {
        $this->storeStmtParam(StmtDataType::Int, $value);
    }

    public function bindStatementParamBool(bool $value)
    {
        $this->storeStmtParam(StmtDataType::Bool, $value);
    }

    public function bindStatementParamFloat(float $value)
    {
        $this->storeStmtParam(StmtDataType::Float, $value);
    }

    public function bindStatementParamString(string $value)
    {
        $this->storeStmtParam(StmtDataType::String, $value);
    }

    public function bindStatementParamTimestamp(int $timestamp)
    {
        $formattedTimestamp = date('Y-m-d H:i:s', $timestamp);
        $this->storeStmtParam(StmtDataType::String, $formattedTimestamp);
    }

    /**
     * Executes a prepared statement.
     * @param bool $onSelectNumRowsOnly if true: drops the data and returns just the amount of selected rows on a SELECT statement.
     * @param bool $onSingleResultUseArray if true: returns even a single result row as a array
     * @return InternalStatus OnSuccess with Status ValidDatabaseRequest or ValidCreation, RequestStatus error otherwise.
     */
    public function executeStatement(bool $onSelectNumRowsOnly = false, bool $onSingleResultUseArray = false): InternalStatus
    {
        $status = new InternalStatus(RequestStatus::Undefined);
        if ($this->m_StatementStatus == StmtStatus::InPreparation) {
            // Bind all collected params
            $this->bindAllStmtParams();
        }

        // Execute the statement
        if ($this->m_StatementStatus == StmtStatus::PreparationComplete) {
            if (DATABASE_EXECUTE_ENABLED) {
                try {
                    $this->m_StatementStatus = ($this->m_CurrentStatement->execute()) ? StmtStatus::ExecutionSuccess : StmtStatus::ExecutionError;
                } catch (mysqli_sql_exception $e) {
                    $this->m_StatementStatus = StmtStatus::ExecutionError;
                    dev_var_dump($e);
                }
            } else {
                dev_var_dump($this->m_CurrentStatement);
                $this->m_StatementStatus = StmtStatus::ExecutionDisabled;
            }
        }

        // Analyse results
        if ($this->m_StatementStatus == StmtStatus::ExecutionSuccess) {
            $requestStatus = RequestStatus::ValidDatabaseRequest;
            $numAffectedRows = $this->m_CurrentStatement->affected_rows; // rows affected by INSERT, DELETE or UPDATE
            $newEntryId = 0;
            $data = array();
            if ($numAffectedRows > 0) {
                // for INSERT the newly inserted id is returned
                if ($this->m_CurrentStatement->insert_id > 0) {
                    $newEntryId = $this->m_CurrentStatement->insert_id;
                    $requestStatus = RequestStatus::ValidCreation;
                }
            } else {
                $result = $this->m_CurrentStatement->get_result();
                $numAffectedRows = $result->num_rows;
                if (!$onSelectNumRowsOnly) {
                    // read data only if requested
                    if (!$onSingleResultUseArray && $result->num_rows == 1) {
                        // single result
                        $data = $result->fetch_assoc();
                    } else {
                        while ($row = $result->fetch_assoc()) {
                            array_push($data, $row);
                        }
                    }
                }
            }

            $status = new InternalStatus($requestStatus, $newEntryId, $numAffectedRows);
            $status->setData($data);
        }

        // Error Handling
        switch ($this->m_StatementStatus) {
            case StmtStatus::ExecutionDisabled;
                $status = new InternalStatus(RequestStatus::DatabaseExecutionDisabled);
                break;
            case StmtStatus::ExecutionError:
                $status = new InternalStatus(RequestStatus::DatabaseStmtExecutionError);
            case StmtStatus::InvalidParamBinding;
                $status = new InternalStatus(RequestStatus::DatabaseStmtParamsError);
                break;
        }

        // Whatever the current status of the statement is, after "executing" it resets.
        $this->resetCurrentStatement();

        return $status;
    }

    private function createDatabase()
    {
        if (!DEVELOPER_MODE_ENABLED) {
            return;
        }

        if (!$this->isDbHandleValid()) {
            return;
        }

        // manual query, since DatabaseHandler::executeQuery() does only work with an initialized
        // database handler - which requires a selected database (that we might be about to create).
        if (DATABASE_EXECUTE_ENABLED) {
            $sql = "CREATE DATABASE IF NOT EXISTS " . self::DB_DEFAULT_DATABASE . ";";
            try {
                $this->m_DbHandle->query($sql);
            } catch (mysqli_sql_exception $e) {
                dev_var_dump($e);
            }
        }
    }

    private function initializeDatabase()
    {
        try {
            $this->m_DbHandle = new mysqli(self::DB_HOST, self::DB_USER, self::DB_PASSWORD, null, self::DB_PORT);

            $this->createDatabase();

            if (DATABASE_EXECUTE_ENABLED) {
                $this->m_DbHandle->select_db(self::DB_DEFAULT_DATABASE);
                $this->m_IsDatabaseSelected = true;
            } else {
                trigger_error("Database Query Execution is disabled.", E_USER_WARNING);
            }
        } catch (mysqli_sql_exception $e) {
            dev_var_dump($e);
        }
    }

    private function isInitialized(): bool
    {
        return $this->isDbHandleValid() && $this->m_IsDatabaseSelected;
    }

    private function isDbHandleValid(): bool
    {
        return !is_null($this->m_DbHandle);
    }

    private function storeStmtParam(StmtDataType $dataType, mixed $value)
    {
        if ($this->m_StatementStatus != StmtStatus::InPreparation) {
            // Statement has to be in preparation mode!
            return;
        }

        array_push($this->m_StatementDataTypes, $dataType);
        array_push($this->m_StatementValues, $value);
    }

    private function bindAllStmtParams()
    {
        if (count($this->m_StatementValues) > 0) {
            $types = '';
            foreach ($this->m_StatementDataTypes as &$dataType) {
                // Build data type string for binding
                switch ($dataType) {
                    case StmtDataType::Int;
                    case StmtDataType::Bool;
                        $types .= 'i';
                        break;
                    case StmtDataType::Float;
                        $types .= 'd';
                        break;
                    case StmtDataType::String;
                    default;
                        $types .= 's';
                        break;
                }
            }

            try {
                $isBindingValid = $this->m_CurrentStatement->bind_param($types, ...$this->m_StatementValues);
                $this->m_StatementStatus = $isBindingValid ? StmtStatus::PreparationComplete : StmtStatus::InvalidParamBinding;
            } catch (ArgumentCountError $e) {
                $this->m_StatementStatus = StmtStatus::InvalidParamBinding;
            }
        } else {
            // No params -> should work aswell...?
            // TODO: that just means no params provided! but maybe params were required?
            $this->m_StatementStatus = StmtStatus::PreparationComplete;
        }
    }

    private function resetCurrentStatement()
    {
        unset($this->m_CurrentStatement);
        $this->m_StatementStatus = StmtStatus::Inactive;
        $this->m_StatementDataTypes = array();
        $this->m_StatementValues = array();
    }

    private function executeQuery(string $sql): bool
    {
        if (!DATABASE_EXECUTE_ENABLED) {
            dev_var_dump($sql);
            return false;
        }

        if ($this->isInitialized()) {
            try {
                $this->m_DbHandle->query($sql);
                // TODO: returns mysqli_result -> anaylize for return value
            } catch (mysqli_sql_exception $e) {
                dev_var_dump($e);
            }
        }

        return false;
    }

    private const DB_HOST = 'mariadb_php';
    private const DB_USER = 'maria';
    private const DB_PASSWORD = 'maria123';
    private const DB_DEFAULT_DATABASE = 'mariadb';
    private const DB_PORT = 3306;

    private $m_DbHandle = null;
    private $m_IsDatabaseSelected = false;

    private StmtStatus $m_StatementStatus = StmtStatus::Inactive;
    private mysqli_stmt $m_CurrentStatement;
    private array $m_StatementDataTypes = [];
    private array $m_StatementValues = [];
}
