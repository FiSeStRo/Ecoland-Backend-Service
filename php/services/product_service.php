<?php

class ProductDefinition{
    use DefinitionDBInitTrait;

    public readonly int $id;
    public readonly string $token_name;

    public function __construct(int $def_id) {
        $db = new DatabaseHandler;
        $sql = "SELECT * FROM " . DbTables::DefProducts->value . " WHERE id = ?";
        if( $db->createStatement($sql) ){
            $db->bindStatementParamInt($def_id);
            $status = $db->executeStatement();
            
            $this->initFromDatabaseResult($status->getData(), [
                "id",
                "token_name",
            ]);
        }
    }
}

?>