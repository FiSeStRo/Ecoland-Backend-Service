<?php

enum DefinitionDataType{
    case OriginTableId;
    case RelationSingleValue;
    case RelationObjectProperty;
}

class SetupService
{
    use HasDatabaseAccessTrait;

    public function setupDatabase() : InternalStatus
    {
        if (!DEVELOPER_MODE_ENABLED) {
            return new InternalStatus(RequestStatus::DevModeRequired);
        }

        if($this->m_Db->createTables($this)){
            return new InternalStatus(RequestStatus::Valid);
        }

        return new InternalStatus(RequestStatus::GeneralError);
    }

    public function initDefinitionData() : InternalStatus
    {
        if (!DEVELOPER_MODE_ENABLED) {
            return new InternalStatus(RequestStatus::DevModeRequired);
        }

        // check definition data list for files to parse and include
        foreach (array_keys(self::SQL_DEFINITION_DATA_LIST) as $definitionName) {
            $definition = self::SQL_DEFINITION_DATA_LIST[$definitionName];
            if ($definition[self::SQL_DEFINITION_FROM_FILE]) {
                // search for a matching json file
                $pathToJsonFile = self::CONFIG_STATIC_DATA_PATH . $definitionName . ".json";
                $fileContent = file_get_contents($pathToJsonFile);
                if ($fileContent !== false) {
                    $data = json_decode($fileContent);
                    $this->analyseDataForBaseDefinition($definitionName, $data);
                }
            }
        }

        return new InternalStatus(RequestStatus::Valid);
    }

    public function resetDatabase() : InternalStatus
    {
        if (!DEVELOPER_MODE_ENABLED) {
            return new InternalStatus(RequestStatus::DevModeRequired);
        }
        
        $isDeleteSuccess = $this->m_Db->deleteTables($this);
        $isCreateSuccess = $this->m_Db->createTables($this);
        if( $isDeleteSuccess && $isCreateSuccess){
            return new InternalStatus(RequestStatus::Valid);
        }
        return new InternalStatus(RequestStatus::GeneralError);
    }

    public function getCreateTableQueries(): array
    {
        return $this->buildCreateTableQueries();
    }

    public function getDeleteTableQuery(): string
    {
        $sql = "DROP TABLE IF EXISTS ";
        $tableList = implode(', ', array_keys(self::SQL_TABLE_DEFINITION));
        $sql .= $tableList;
        return $sql;
    }

    private function analyseDataForBaseDefinition(string $definitionName, $data)
    {
        // find definition
        $definition = self::SQL_DEFINITION_DATA_LIST[$definitionName] ?? false;
        if ($definition !== false) {
            // check each object (= table row) from the json data
            foreach ($data as $row) {
                $valueList = [];

                // try to find a property in the object for each column defined
                foreach ($definition[self::SQL_DEFINITION_COLUMNS] as $column) {
                    if (isset($row->$column)) {
                        $valueList[$column] = $row->$column;
                    }
                    // TODO: Error handling when missing a column
                }

                $this->insertDefinitionData($definitionName, $valueList);

                $numColumns = count($valueList);
                if ($numColumns > 0) {
                    // build sql state from collected columns and values
                    $sql = "INSERT INTO $definitionName (" . implode(", ", array_keys($valueList)) . ") VALUES (" . implode(", ", array_fill(0, $numColumns, "?")) . ")";
                }

                // check for relationship-table definitions (there might be multiple)
                if (isset($definition[self::SQL_DEFINITION_RELATION_TARGET])) {
                    foreach ($definition[self::SQL_DEFINITION_RELATION_TARGET] as $objectProperty => $relTableName) {
                        $relData = $row->$objectProperty ?? null;
                        if (is_array($relData)) {
                            // extract the data object from the definition and evaluate it row by row
                            foreach ($relData as $relRow) {
                                $originTableId = $row->id ?? 0;
                                if ($originTableId > 0) {
                                    $this->analyseDataForRelationshipDefinition($relTableName, $originTableId, $relRow);
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    private function analyseDataForRelationshipDefinition(string $relTableName, int $originTableId, object|int $relData)
    {
        $tableDefinition = self::SQL_DEFINITION_DATA_LIST[$relTableName] ?? null;
        if (isset($tableDefinition) && $tableDefinition[self::SQL_DEFINITION_FROM_FILE] == false) {           
            $valueList = [];
            // try to find a property in the object for each column defined
            foreach ($tableDefinition[self::SQL_DEFINITION_RELATION_COLUMNS] as $columnName => $valueType) {
                if( $valueType == DefinitionDataType::OriginTableId ){
                    $valueList[$columnName] = $originTableId;
                }
                else if( $valueType == DefinitionDataType::RelationSingleValue){                   
                    if( is_int($relData) ){
                        $valueList[$columnName] = $relData;
                    }
                    // TODO errorHandling!
                }
                else if( isset($relData->$columnName) ){
                    // relationship data has to be an object => search for property
                    $valueList[$columnName] = $relData->$columnName;
                }
            }
            $this->insertDefinitionData( $relTableName, $valueList );
        }
    }

    private function insertDefinitionData( string $tableName, array $valueList ){       
        $numColumns = count($valueList);
        if ($numColumns > 0) {
            // build sql state from collected columns and values
            $sql = "INSERT INTO $tableName (" . implode(", ", array_keys($valueList)) . ") VALUES (" . implode(", ", array_fill(0, $numColumns, "?")) . ")";
            if( $this->m_Db->createStatement($sql) ){
                foreach($valueList as $value){
                    if( is_bool($value) ){
                        $this->m_Db->bindStatementParamBool($value);
                    }
                    else if( is_int($value)){
                        $this->m_Db->bindStatementParamInt($value);
                    }
                    else if( is_float($value)){
                        $this->m_Db->bindStatementParamFloat($value);
                    }
                    else{
                        $this->m_Db->bindStatementParamString($value);
                    }                    
                }
                $this->m_Db->executeStatement();
            }
        }
    }

    private function buildCreateTableQueries(): array
    {
        $createQueries = [];
        foreach (array_keys(self::SQL_TABLE_DEFINITION) as $tableName) {
            $sql = $this->buildCreateTableQuery($tableName);
            if (!empty($sql)) {
                array_push($createQueries, $sql);
            }
        }

        return $createQueries;
    }

    private function buildCreateTableQuery(string $tableName): string
    {
        $sql = "";
        $sqlDefinition = $this->getSQLDefinitionForTable($tableName);
        if ($sqlDefinition !== false && $sqlDefinition !== "") {
            $sql = "CREATE TABLE IF NOT EXISTS $tableName(" . $sqlDefinition . ")";
        }

        return $sql;
    }

    private function getSQLDefinitionForTable(string $tableName): mixed
    {
        if (array_key_exists($tableName, self::SQL_TABLE_DEFINITION)) {
            $tableData = self::SQL_TABLE_DEFINITION[$tableName];
            $columnDefinitionString = "";
            // search for column definitions
            if (array_key_exists(self::SQL_DEFINITION_COLUMNS, $tableData)) {
                $columnDefinitionString = implode(', ', $tableData[self::SQL_DEFINITION_COLUMNS]);
            }

            // search for primary key definition
            if (array_key_exists(self::SQL_DEFINITION_PRIMARY, $tableData)) {
                $primaryKeyColumnNames = $tableData[self::SQL_DEFINITION_PRIMARY];
                $primaryKeyDefinitionString = "";
                // definitions required as array (for single- or multi-column primary keys)
                if (is_array($primaryKeyColumnNames)) {
                    $primaryKeyDefinitionString = implode(', ', $primaryKeyColumnNames);
                    $columnDefinitionString .= ", CONSTRAINT " . $tableName . "_pk PRIMARY KEY ($primaryKeyDefinitionString)";
                }
            }
            return $columnDefinitionString;
        }
        return false;
    }   

    private const SQL_DEFINITION_COLUMNS = 'columns'; // name of columns in table
    private const SQL_DEFINITION_PRIMARY = 'primary'; // name of primary key column
    private const SQL_DEFINITION_FROM_FILE = 'from_file'; // from json file to read definition data?
    private const SQL_DEFINITION_RELATION_TARGET = 'relation_target'; // Which properties describe another rel-table   
    private const SQL_DEFINITION_RELATION_COLUMNS = 'relation_columns'; // defines columns in relationship table that get the original values assigned

    private const CONFIG_STATIC_DATA_PATH = '/var/www/html/config/';

    private const SQL_TABLE_DEFINITION = [
        // user accounts table
        DbTables::Users->value => [
            self::SQL_DEFINITION_COLUMNS => [
                'id INT NOT NULL AUTO_INCREMENT',
                'username VARCHAR(255)',
                'password VARCHAR(255)',
                'role INT',
                'time_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP()',
                'time_last_activity TIMESTAMP DEFAULT CURRENT_TIMESTAMP()',
            ],
            self::SQL_DEFINITION_PRIMARY => ['id'],
        ],
        // user resources table
        DbTables::UserResources->value => [
            self::SQL_DEFINITION_COLUMNS => [
                'user_id INT NOT NULL',
                'money DECIMAL',
                'prestige INT',
            ],
            self::SQL_DEFINITION_PRIMARY => ['user_id'],
        ],
        // buildings table
        DbTables::Buildings->value => [
            self::SQL_DEFINITION_COLUMNS => [
                'id INT NOT NULL AUTO_INCREMENT',
                'user_id INT',
                'def_id INT',
                'name VARCHAR(255)',
                'time_build TIMESTAMP DEFAULT CURRENT_TIMESTAMP()',
            ],
            self::SQL_DEFINITION_PRIMARY => ['id'],
        ],
        // buildings definitions table (configuration)
        DbTables::DefBuildings->value => [
            self::SQL_DEFINITION_COLUMNS => [
                'id INT NOT NULL',
                'token_name VARCHAR(255)',
                'base_construction_cost DECIMAL',
                'base_construction_time INT',
            ],
            self::SQL_DEFINITION_PRIMARY => ['id'],
        ],
        // products definitions table (configuration)
        DbTables::DefProducts->value => [
            self::SQL_DEFINITION_COLUMNS => [
                'id INT NOT NULL',
                'token_name VARCHAR(255)',
            ],
            self::SQL_DEFINITION_PRIMARY => ['id'],
        ],
        // production definitions table (configuration, basics)
        DbTables::DefProduction->value => [
            self::SQL_DEFINITION_COLUMNS => [
                'id INT NOT NULL',
                'token_name VARCHAR(255)',
                'cost DECIMAL',
                'base_duration INT',
            ],
            self::SQL_DEFINITION_PRIMARY => ['id'],
        ],
        // production recipe table (configuration, n-to-n)
        DbTables::DefProductionRecipe->value => [
            self::SQL_DEFINITION_COLUMNS => [
                'production_id INT NOT NULL',
                'product_id INT NOT NULL',
                'is_input BOOL',
                'amount INT NOT NULL',
            ],
            self::SQL_DEFINITION_PRIMARY => [
                'production_id',
                'product_id',
                'is_input',
            ],
        ],
        DbTables::DefBuildingProduction->value => [
            self::SQL_DEFINITION_COLUMNS => [
                'building_id INT NOT NULL',
                'production_id INT NOT NULL',
            ],
            self::SQL_DEFINITION_PRIMARY => [
                'building_id',
                'production_id',
            ]
        ],
        // "Storage" table
        DbTables::Storage->value => [
            self::SQL_DEFINITION_COLUMNS => [
                'building_id INT NOT NULL',
                'product_id INT NOT NULL',
                'amount INT',
                'capacity INT DEFAULT 500' //TODO: define default capacity via config
            ],
            self::SQL_DEFINITION_PRIMARY => [
                'building_id',
                'product_id',
            ],
        ],
        // manufacturing orders
        DbTables::ManufacturingOrders->value => [
            self::SQL_DEFINITION_COLUMNS => [
                'id INT NOT NULL AUTO_INCREMENT',
                'building_id INT NOT NULL',
                'production_id INT NOT NULL',
                'time_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP()',
                'time_end TIMESTAMP DEFAULT CURRENT_TIMESTAMP()',
                'cycles INT DEFAULT 1',
                'is_completed BOOLEAN DEFAULT FALSE',
            ],
            self::SQL_DEFINITION_PRIMARY => [
                'id',
            ],
        ],
    ];

    // Defines what to load from json-definitions into database
    // table => [columns to fill]
    private const SQL_DEFINITION_DATA_LIST = [
        // Definition Products
        DbTables::DefProducts->value => [
            self::SQL_DEFINITION_FROM_FILE => true,
            self::SQL_DEFINITION_COLUMNS => ['id', 'token_name',],
        ],
        // Definition Buildings
        DbTables::DefBuildings->value => [
            self::SQL_DEFINITION_FROM_FILE => true,
            self::SQL_DEFINITION_COLUMNS => ['id', 'token_name', 'base_construction_cost', 'base_construction_time'],
            self::SQL_DEFINITION_RELATION_TARGET => [
                'productions' => DbTables::DefBuildingProduction->value,
            ],
        ],
        // Definition Production
        DbTables::DefProduction->value => [
            self::SQL_DEFINITION_FROM_FILE => true,
            self::SQL_DEFINITION_COLUMNS => ['id', 'token_name', 'cost', 'base_duration'],
            self::SQL_DEFINITION_RELATION_TARGET => [
                'products' => DbTables::DefProductionRecipe->value,
            ],
        ],
        // Definition relation buildings <-> productions
        DbTables::DefBuildingProduction->value => [
            self::SQL_DEFINITION_FROM_FILE => false,
            self::SQL_DEFINITION_RELATION_COLUMNS => [
                'building_id' => DefinitionDataType::OriginTableId,
                'production_id' => DefinitionDataType::RelationSingleValue,
            ],
        ],
        // Definition relation production <-> products (recipe)
        DbTables::DefProductionRecipe->value => [
            self::SQL_DEFINITION_FROM_FILE => false,
            self::SQL_DEFINITION_RELATION_COLUMNS => [
                'production_id' => DefinitionDataType::OriginTableId,
                'product_id' => DefinitionDataType::RelationObjectProperty,
                'is_input' => DefinitionDataType::RelationObjectProperty,
                'amount' => DefinitionDataType::RelationObjectProperty,
            ],
        ],
    ];
}
