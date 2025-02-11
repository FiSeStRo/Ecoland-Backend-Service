package service

import (
	"io"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Service is up and helthy")

}
