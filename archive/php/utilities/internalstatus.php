<?php

class InternalStatus{
    public function __construct(RequestStatus $status, int $newId = 0, int $numAffectedRows = 0) {     
        $this->m_RequestStatus = $status;  
        $this->m_NewId = $newId;
        $this->m_NumAffectedRows = $numAffectedRows;
    }

    public function setData(mixed $data){
        $this->m_Data = $data;
    }

    public function setAuthTokens(string $accessToken, string $refreshToken = ""){
        $this->m_AccessToken = $accessToken;
        $this->m_RefreshToken = $refreshToken;
    }

    public function getRequestStatus() : RequestStatus{
        return $this->m_RequestStatus;
    }

    public function getNewId() : int{
        return $this->m_NewId ?? 0;
    }

    public function getNumAffectedRows() : int{
        return $this->m_NumAffectedRows ?? 0;
    }

    public function getAuthAccessToken() : string{
        return $this->m_AccessToken ?? "";
    }

    public function getAuthRefreshToken() : string{
        return $this->m_RefreshToken ?? "";
    }

    public function getData(){
        return $this->m_Data;
    }

    public function isItTrue() : bool{
        return $this->getRequestStatus() == RequestStatus::Valid;
    }

    public function isValidStatus() : bool{
        return ($this->m_RequestStatus == RequestStatus::Valid 
                || $this->m_RequestStatus == RequestStatus::ValidCreation
                || $this->m_RequestStatus == RequestStatus::ValidDatabaseRequest
                || $this->m_RequestStatus == RequestStatus::ValidAuthentification);
    }

    private RequestStatus $m_RequestStatus;
    private readonly int $m_NewId;
    private readonly int $m_NumAffectedRows;
    private readonly string $m_RefreshToken;
    private readonly string $m_AccessToken;   
    private $m_Data;
}

?>