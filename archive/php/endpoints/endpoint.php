<?php

enum CommandType
{
    case Get;
    case GetWithId;
    case PostFormData; // Content-Type multipart/form-data
    case PostJson; // Content-Type application/json
    case Patch; // Updating part of data
    case Delete;
}

trait CommandHandlingTrait{    
    public function handleCommand() : InternalStatus{
        $status = $this->getRequestedCommandPath();
        $commandPath = $status->getData();
        if( !empty($commandPath)){  
            try{
                return $this->$commandPath();
            }         
            catch(TypeError $e){
                return new InternalStatus(RequestStatus::MissingCommandInternalStatus);
            }
        }
        return $status;
    }
}

class Endpoint
{
    public function __construct(string $command, array $params, AuthenticationHandler &$authHandler)
    {
        $this->m_Command = $command;
        $this->m_RequestMethod = $_SERVER['REQUEST_METHOD'];
        $this->m_Params = $params;
        $this->m_Db = new DatabaseHandler();
        $this->m_AuthHandler =& $authHandler;
    }

    protected function getDatabaseHandler(): ?DatabaseHandler
    {
        return $this->m_Db;
    }

    protected function getCommand(): string
    {
        return $this->m_Command;
    }

    protected function getParams(): array
    {
        $currenctCommandType = $this->getCurrentCommandType();
        if( $currenctCommandType === false)
        {
            // Command probably invalid -> no params
            return array();
        }

        switch($currenctCommandType){

            case CommandType::GetWithId;
            case CommandType::Delete;
                return $this->m_Params['GET'];
            case CommandType::PostFormData;
            case CommandType::PostJson;
                return $this->m_Params['POST'];
            case CommandType::Get;
            default;
                return array();
                break;
        }       
    }

    protected function getCurrentUserId(){
        return $this->m_AuthHandler->getUserIdFromAuthToken();
    }

    protected function getCurrentCommandType() : CommandType | bool{   
        if( !$this->isCommandRegistered($this->m_Command, $this->m_RequestMethod) ){
            return false;
        }
        
        $command = $this->m_RegisteredCommands[$this->m_RequestMethod][$this->m_Command];
        return $command[self::KEY_CMD_TYPE];
    }

    protected function registerCommand(string $commandName, string $commandPath, CommandType $commandType = CommandType::Get, UserLevel $requiredUserLevel = UserLevel::Unregistered): bool
    {
        $commandName = trim($commandName);
        $commandName = strtolower($commandName);
        $commandPath = trim($commandPath);

        if ($commandName == '') {
            return false;
        }

        if ($commandPath == '') {
            return false;
        }

        if (array_key_exists($commandName, $this->m_RegisteredCommands) === true) {
            return false;
        }

        $userLevelValue = $requiredUserLevel->value;

        // Add new commands in two levels: 
        // 1) check request method
        $requestMethod = $this->getRequestMethodFromCommandType($commandType);
        if( !isset($this->m_RegisteredCommands[$requestMethod])){
            $this->m_RegisteredCommands[$requestMethod] = [];
        }

        // 2) check if command already exists for that request method
        if( isset($this->m_RegisteredCommands[$requestMethod][$commandName])){
            return false;
        }

        $this->m_RegisteredCommands[$requestMethod][$commandName] = [
            self::KEY_CMD_PATH => $commandPath,
            self::KEY_CMD_TYPE => $commandType,
            self::KEY_CMD_USER_LVL => $userLevelValue,
        ];
        return true;
    }

    protected function getRequestedCommandPath(): InternalStatus
    {
        // Retreives the command path of the requested command if it is registered and valid.
        // Empty string otherwise. 
        // Multiple validation steps are required.
        $foundCommand = '';
        $commandValidationStatus = RequestStatus::Valid;

        // - Is the command registered?
        if( !$this->isCommandRegistered( $this->m_Command, $this->m_RequestMethod) ){
            $commandValidationStatus = RequestStatus::UnknownCommand;
        }
        else{     
            $command = $this->m_RegisteredCommands[$this->m_RequestMethod][$this->m_Command];

            if ($command[self::KEY_CMD_USER_LVL] > UserLevel::Unregistered && !$this->isAuthenticationValid()) {
                // - Is Authentication required?
                $commandValidationStatus = RequestStatus::AuthenticationInvalid;
            } else if($command[self::KEY_CMD_USER_LVL] > UserLevel::User) {
                // - Does the user have the required user level to execute the command?

            } else if (!method_exists($this, $command[self::KEY_CMD_PATH])) {
                // - Does the method for the command exist?
                $commandValidationStatus = RequestStatus::MissingCommandImplementation;                
            } else {
                // - Is the required data available? (GET, POST, ...)
                $typeValidation = $this->validateCommandType($command[self::KEY_CMD_TYPE]);
                if ($typeValidation != RequestStatus::Valid) {
                    $commandValidationStatus = $typeValidation;
                }
            }
            $foundCommand = $command[self::KEY_CMD_PATH];
        }

        $status = new InternalStatus($commandValidationStatus);
        $status->setData(($commandValidationStatus == RequestStatus::Valid) ? $foundCommand : '');
        return $status;
    }

    private function getRequestMethodFromCommandType(CommandType $commandType) : string{
        switch ($commandType){
            case CommandType::Get:
            case CommandType::GetWithId:
                break;
            case CommandType::PostFormData:
            case CommandType::PostJson:
                return 'POST';
                break;
            case CommandType::Patch:
                return 'PATCH';
            case CommandType::Delete;
                return 'DELETE';
                break;
        }

        return 'GET';
    }

    private function isCommandRegistered(string $commandName, string $requestMethod) : bool{        
        return array_key_exists($requestMethod, $this->m_RegisteredCommands) && array_key_exists($commandName, $this->m_RegisteredCommands[$requestMethod]);
    }

    private function isAuthenticationValid(): bool
    {       
        return $this->m_AuthHandler->validateAuthToken();
    }
    
    protected function createAuthToken(int $userId, bool $isAuthToken) : string{
        return ($userId > 0) ?$this->m_AuthHandler->createAuthToken($userId, $isAuthToken) : '';
    }

    private function validateCommandType(CommandType $type): RequestStatus
    {
        if ($type == CommandType::GetWithId || $type == CommandType::Delete){
            if(empty($this->m_Params['GET']) || intval($this->m_Params['GET']['id']) == 0) {                
                return ($type == CommandType::GetWithId) ? RequestStatus::MissingRequiredParamsGet : RequestStatus::MissingRequiredParamsDelete;
            }
        }  else if ($type == CommandType::PostFormData && empty($this->m_Params['POST'])) {
            return RequestStatus::MissingRequiredParamsPost;
        } else if($type == CommandType::Patch && empty($this->m_Params['PATCH'])){
            return RequestStatus::MissingRequiredParamsPatch;
        }
        return RequestStatus::Valid;
    }

    protected function debugDumpRegisteredCommands()
    {
        dev_var_dump($this->m_RegisteredCommands);
    }

    private const KEY_CMD_PATH = 0; // command path
    private const KEY_CMD_TYPE = 1; // type of command
    private const KEY_CMD_USER_LVL = 2; // required user level to be able to execute command

    private string $m_Command = ''; // Currently requested command
    private string $m_RequestMethod = '';
    private array $m_Params = array();

    private DatabaseHandler $m_Db;   
    private AuthenticationHandler $m_AuthHandler;

    private array $m_RegisteredCommands = array();
}
