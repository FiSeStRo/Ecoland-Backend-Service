<?php
// Settings
declare(strict_types=1);

// Constants
define('CONTENT_TYPE_APPLICATION_JSON', 'application/json'); // Content-header

define ('STORAGE_DEFAULT_CAPACITY', 500); // Default capacity of a new building
define ('PRODUCT_ID_MONEY', 1); // Special id for products that represents money 

enum RequestStatus : string{
    case Undefined = 'Undefined';
    // Success 
    case Valid = 'Valid'; // gets used for "true on bool-response methods (that use InternalStatus)
    case ValidCreation = 'ValidCreation'; // new building, production order, ...
    case ValidDatabaseRequest = 'ValidDatabaseRequest';
    case ValidAuthentification = 'ValidAuthentification';
    // General Errors    
    case DevModeRequired = 'DevModeRequired';
    case GeneralError = 'GeneralError';
    case Uninitialized = 'Unitialized';
    case InvalidInput = 'InvalidInput';
    case InvalidInputEmail = 'InvalidInputEmail'; // Email address input is not of form email
    case AuthenticationInvalid = 'AuthenticationInvalid';
    case CommandInsufficientUserLevel = 'CommandInsufficientUserLevel'; // User does not have required user level to be able to execute command
    case UnknownEndpoint = 'UnknownEndpoint';
    case MissingCommand = 'MissingCommand';
    case UnknownCommand = 'UnknownCommand';
    case MissingCommandImplementation = 'MissingCommandImplementation';
    case MissingCommandInternalStatus = 'MissingCommandInternalStatus'; // command valid, but doesn't return InternalStatus
    case MissingRequiredParamsGet = 'MissingRequiredParamsGet';   
    case MissingRequiredParamsPost = 'MissingRequiredParamsPost'; 
    case MissingRequiredParamsPatch = 'MissingRequiredParamsPatch';
    case MissingRequiredParamsDelete = 'MissingRequiredParamsDelete';
    // Definition Data Error
    case UnknownBuildingDefinition = 'UnknownBuildingDefinition';
    case UnknownProductDefinition = 'UnknownProductDefinition';
    case UnknownProductionDefinition = 'UnknownProuctionDefinition';
    // Database Errors
    case DatabaseExecutionDisabled = 'DatabaseExecutionDisabled';
    case DatabaseStmtCreationError = 'DatabaseStmtCreationError';
    case DatabaseStmtParamsError = 'DatabaseStmtParamsError';
    case DatabaseStmtExecutionError = 'DatabaseStmtExecutionError';
    // Buildings Errors
    case BuildingDefinitionMissing = 'BuildingDefinitionMissing';
    case BuildingDoesNotExist = 'BuildingDoesNotExist';
    case BuildingDoesNotBelongToUser = 'BuildingDoesNotBelongToUser';
    case BuildingHasUnfinishedProduction = 'BuildingHasUnfinishedProduction';
    case BuildingCannotDoProduction = 'BuildingCannotDoProduction';   
    // Storage
    case StorageMissingProducts = 'StorageMissingProducts';
    case StorageMissingCapacity = 'StorageMissingCapacity';
    case StorageInvalidAmount = 'StorageInvalidAmount';
    // User (Resources) Errors
    case UserDoesAlreadyExist = 'UserDoesAlreadyExist';
    case UserEmailDoesAlreadyExist = 'UserEmailDoesAlreadyExist';
    case UserDoesNotExist = 'UserDoesNotExist';
    case UserInsufficientFunds = 'UserInsufficientFunds';
    // Financial Errors
    case FinanceWithdrawError = 'FinanceWithdrawError';
    case FinanceDepositError = 'FinanceDepositError';
    // Product Errors
    // Production Errors
    case ProductionInvalidTimestamp = 'ProductionInvalidTimestamp';
    case ProductionAlreadyCompleted = 'ProductionAlreadyCompleted';
    case ProductionDoesNotExist = 'ProductionDoesNotExist';
    case ProductionDoesNotBelongToUser = 'ProductionDoesNotBelongToUser';
}

enum UserLevel : int{
    case Unregistered = -1;
    case User = 0;
    case Administrator = 50;
}
?>
