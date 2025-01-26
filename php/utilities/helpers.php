<?php

function dev_var_dump($data){
    if( !defined('DEVELOPER_MODE_ENABLED') || !DEVELOPER_MODE_ENABLED ){
        return;
    }
    
    echo "<pre>";
        var_dump($data);
        echo "</pre>";
}
?>