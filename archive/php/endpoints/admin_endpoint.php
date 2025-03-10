<?php

class AdministratorEndpoint extends Endpoint
{
    use CommandHandlingTrait;

    public function __construct(string $command, array $params, AuthenticationHandler &$authHandler)
    {
        parent::__construct($command, $params, $authHandler);

        $this->registerCommand('setupDatabase', 'setupDatabase');
        $this->registerCommand('initDefinitionData', 'initDefinitionData');
        $this->registerCommand('resetDatabase', 'resetDatabase');
    }

    private function setupDatabase()
    {
        if (!DEVELOPER_MODE_ENABLED) {
            return;
        }

        $setupService = new SetupService();
        return $setupService->setupDatabase();
    }

    private function initDefinitionData(): InternalStatus
    {
        if (!DEVELOPER_MODE_ENABLED) {
            return new InternalStatus(RequestStatus::DevModeRequired);
        }

        $setupService = new SetupService();
        return $setupService->initDefinitionData();
    }

    private function resetDatabase(): InternalStatus
    {
        if (!DEVELOPER_MODE_ENABLED) {
            return new InternalStatus(RequestStatus::DevModeRequired);
        }

        $setupService = new SetupService();
        return $setupService->resetDatabase();
    }
}
