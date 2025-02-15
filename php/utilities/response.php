<?php

enum ResponseType{
    case None; // respond with nothing
    case NewEntry; // respond with new id
    case Authentication; // respond with tokens
    case Error; // respond with error (error code and message token)
}

enum HttpResponseCode : int{
    case Ok = 200;
    case Created = 201;
    case BadRequest = 400;
    case Unauthorized = 401;
    case NotFound = 404;
    case NotImplemented = 501;
}

/**
 * Response class for handling data that gets send back to the user.
 */
class Response{
    public function __construct() {
        $this->m_Content = new stdClass;
    }

    public function handleInternalStatus(InternalStatus &$status){
        // TODOs: general checking -> error for db query execution disabled

        $this->handleResponseCodes($status); // http response codes
        $this->handleResponseType($status); // internal response types
        $this->handleResponseContent($status); // actual content that gets returned to the user
    }

    public function send(){
        http_response_code($this->m_HttpResponseCode->value);
        header('Content-type: ' . CONTENT_TYPE_APPLICATION_JSON);
        $responseData = json_encode($this->m_Content);       
        echo $responseData;
    }

    public function getRequestStatus() : RequestStatus{
        return $this->m_RequestStatus;
    }

    public function isRequestStatusValid() : bool{
        return $this->getRequestStatus() == RequestStatus::Valid;
    }

    private function handleResponseCodes(InternalStatus &$status){
        $this->m_RequestStatus = $status->getRequestStatus();

        switch($this->m_RequestStatus){
            case RequestStatus::AuthenticationInvalid:
                $this->m_HttpResponseCode = HttpResponseCode::Unauthorized; // 401
                break;
            case RequestStatus::UnknownEndpoint;
            case RequestStatus::UnknownCommand;
            case RequestStatus::MissingCommandImplementation;
                $this->m_HttpResponseCode = HttpResponseCode::NotFound; // 404
                break;
            case RequestStatus::InvalidInput;
            case RequestStatus::MissingCommand;
            case RequestStatus::MissingCommandInternalStatus;
            case RequestStatus::MissingRequiredParamsGet;
            case RequestStatus::MissingRequiredParamsPost;
            case RequestStatus::UserDoesAlreadyExist;
            case RequestStatus::UserDoesNotExist;
            case RequestStatus::UserInsufficientFunds;
            case RequestStatus::StorageMissingProducts;
            case RequestStatus::StorageMissingCapacity;
            case RequestStatus::StorageInvalidAmount;
                $this->m_HttpResponseCode = HttpResponseCode::BadRequest; // 400
                break;
            case RequestStatus::ValidCreation;
                $this->m_HttpResponseCode = HttpResponseCode::Created; // 201
                break;
            case RequestStatus::Valid;
            case RequestStatus::ValidAuthentification;
                $this->m_HttpResponseCode = HttpResponseCode::Ok; // 200
                break;
            default;
                $this->m_HttpResponseCode = HttpResponseCode::NotImplemented; // 501
                break;
        }
    }

    private function handleResponseType(InternalStatus &$status){
        $this->m_ResponseType = ResponseType::None;
        switch($this->m_RequestStatus){
            case RequestStatus::Undefined;
            case RequestStatus::AuthenticationInvalid:               
            case RequestStatus::UnknownEndpoint;
            case RequestStatus::UnknownCommand;
            case RequestStatus::InvalidInput;
            case RequestStatus::MissingCommand;
            case RequestStatus::MissingCommandInternalStatus;
            case RequestStatus::MissingCommandImplementation;
            case RequestStatus::MissingRequiredParamsGet;
            case RequestStatus::MissingRequiredParamsPost;     
            case RequestStatus::UserDoesAlreadyExist;
            case RequestStatus::UserDoesNotExist;  
            case RequestStatus::UserInsufficientFunds; 
            case RequestStatus::StorageInvalidAmount;
            case RequestStatus::StorageMissingCapacity;
            case RequestStatus::StorageMissingProducts;
            default; 
                $this->m_ResponseType = ResponseType::Error;
                break;
            case RequestStatus::ValidAuthentification;
                $this->m_ResponseType = ResponseType::Authentication;
                break;
            case RequestStatus::ValidCreation;
                $this->m_ResponseType = ResponseType::NewEntry;
                break;
        }
    }

    private function handleResponseContent(InternalStatus &$status){
        switch($this->m_ResponseType){
            case ResponseType::Authentication;
                $this->handleResponseContentAuthentication($status);
                break;
            case ResponseType::NewEntry;
                $this->handleResponseContentNewEntry($status);
                break;
            case ResponseType::Error;
                $this->handleResponseContentError($status);
                break;
            case ResponseType::None;
            default;
                break;
        }
    }

    private function handleResponseContentAuthentication(InternalStatus &$status){
        $this->m_Content = (object)[
            'accessToken' => $status->getAuthAccessToken(),
            'refreshToken' => $status->getAuthRefreshToken(),
        ];
    }

    private function handleResponseContentNewEntry(InternalStatus &$status){
        $this->m_Content = (object)[
            'id' => $status->getNewId(),
        ];
    }

    private function handleResponseContentError(InternalStatus &$status){       
        $token = $this->m_RequestStatus;
        $this->m_Content = (object)[
            'code' => 0,
            'message_token' => $token,
        ];
    }

    private RequestStatus $m_RequestStatus = RequestStatus::Uninitialized;
    private HttpResponseCode $m_HttpResponseCode = HttpResponseCode::Ok;   
    private ResponseType $m_ResponseType = ResponseType::None;   
    private $m_Content;    
}

?>