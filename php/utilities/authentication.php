<?php

use Firebase\JWT\JWT;
use Firebase\JWT\Key;
use Firebase\JWT\ExpiredException;

class AuthenticationHandler{

    public function createAuthToken(int $userId) : string{
        // TODO: probably remove userId later, maybe use valid endpoints?
        $this->m_UserId = $userId;

        $jwtPayload = [
            self::JWT_PAYLOAD_KEY_ISSUER => self::JWT_ISSUER,
            self::JWT_PAYLOAD_KEY_SUBJECT => "$userId",
            self::JWT_PAYLOAD_KEY_ISSUED_AT => time(),
            self::JWT_PAYLOAD_KEY_EXPIRATION => time() + self::JWT_AUTH_TOKEN_LIFESPAN_IN_MINS * 60,
        ];

        $jwt = JWT::encode($jwtPayload, self::JWT_KEY, self::JWT_ALGORITHM);   
        return $jwt;    
    }

    public function validateAuthToken() : bool{
        if( !AUTH_USER_VALIDATION_ENABLED )
        {
            // if authentication is disabled -> always true
            return true;
        }

        // Try to find a token and to decode it
        $decodedToken = $this->decodeTokenFromHeader();
        if( $decodedToken === false ){
            return false;
        }

        // If we got until here, token seems to be valid (in terms of format)
        // Check Issuer                     
        if( !isset($decodedToken[self::JWT_PAYLOAD_KEY_ISSUER]) ||$decodedToken[self::JWT_PAYLOAD_KEY_ISSUER] != self::JWT_ISSUER ){
            return false;
        }

        // Check Expiration Date
        if( !isset($decodedToken[self::JWT_PAYLOAD_KEY_EXPIRATION]) || $decodedToken[self::JWT_PAYLOAD_KEY_EXPIRATION] < time() ){
            return false;
        }
        
        // Check Issued_At Date
        if( !isset($decodedToken[self::JWT_PAYLOAD_KEY_ISSUED_AT]) || $decodedToken[self::JWT_PAYLOAD_KEY_ISSUED_AT] > time() ){
            return false;
        }

        if( !isset($decodedToken[self::JWT_PAYLOAD_KEY_SUBJECT])){
            return false;
        }

        // Valid token -> store user Id 
        $this->m_UserId = $decodedToken[self::JWT_PAYLOAD_KEY_SUBJECT];
        return true;
    }

    public function getUserIdFromAuthToken() : int{
        $this->validateAuthToken();
        return $this->m_UserId;
    }
    
    private function decodeTokenFromHeader() : array|bool{
        try{
            $headerValues = explode(' ', $this->getAuthorizationHeader());
            if( is_array($headerValues) && count($headerValues) >= 2 )
            {
                $jwt = $headerValues[1];               
                $decodedToken = JWT::decode($jwt, new Key(self::JWT_KEY, self::JWT_ALGORITHM));               
                return get_object_vars($decodedToken);
            }
        } 
        catch( DomainException | ExpiredException $e){
            return false;
        }

        return false;
    }

    private function getAuthorizationHeader() : string{
        // thanks SO!
        // https://stackoverflow.com/questions/40582161/how-to-properly-use-bearer-tokens
        $headers = '';
        if (isset($_SERVER['Authorization'])) {
            $headers = trim($_SERVER["Authorization"]);
        }
        else if (isset($_SERVER['HTTP_AUTHORIZATION'])) { //Nginx or fast CGI
            $headers = trim($_SERVER["HTTP_AUTHORIZATION"]);
        } elseif (function_exists('apache_request_headers')) {
            $requestHeaders = apache_request_headers();
            // Server-side fix for bug in old Android versions (a nice side-effect of this fix means we don't care about capitalization for Authorization)
            $requestHeaders = array_combine(array_map('ucwords', array_keys($requestHeaders)), array_values($requestHeaders));
            //print_r($requestHeaders);
            if (isset($requestHeaders['Authorization'])) {
                $headers = trim($requestHeaders['Authorization']);
            }
        }
        return $headers;
    }

    private int $m_UserId = 0;
    private const JWT_ISSUER = 'ecoland';
    private const JWT_KEY = 'd3iMudd4s31g51chT';
    private const JWT_ALGORITHM = 'HS256';

    private const JWT_PAYLOAD_KEY_ISSUER = 'iss';
    private const JWT_PAYLOAD_KEY_SUBJECT = 'sub';
    private const JWT_PAYLOAD_KEY_EXPIRATION = 'exp';
    private const JWT_PAYLOAD_KEY_ISSUED_AT = 'iat';

    private const JWT_AUTH_TOKEN_LIFESPAN_IN_MINS = 10;
    private const JWT_REFRESH_TOKEN_LIFESPAN_IN_MINS = 60;
}

?>