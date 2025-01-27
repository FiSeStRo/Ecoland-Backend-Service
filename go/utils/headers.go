package utils

import "net/http"

// SetHeaderJson is a utils funcitonality to set they Header of a response as JsonFormat
func SetHeaderJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
}
