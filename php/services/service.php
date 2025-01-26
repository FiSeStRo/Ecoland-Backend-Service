<?php

enum DefinitionType : string{
    case Building = 'Building';
    case Production = 'Production';
    case Product = 'Product';
}

trait HasDefinitionTrait{
    public function getDefinition(DefinitionType $type, int $definitionId) : BuildingDefinition|ProductionDefinition|bool{
        $definitionClassName = $type->value . 'Definition';
        try{
            $definition = new $definitionClassName($definitionId);
            return ($definition->isValid) ? $definition : false;
        }
        catch(Exception $e){}
        return false;
    }
}

trait DefinitionDBInitTrait{
    public readonly bool $isValid;

    private function initFromDatabaseResult(array $dbResult, array $propertyList){       
        $isValidResult = true;
        foreach($propertyList as $property){
            if(isset($dbResult[$property])) {
                $this->$property = $dbResult[$property];
            }
            else{
                $isValidResult = false;
                break;
            }
        }
        
        $this->isValid = $isValidResult;
    }
}

?>