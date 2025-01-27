package utils

import "net/http"

// IsMethodPOST is a utils function that checks if the request was a POST if not it writes an http.Error into the ResponseWriter
func IsMethodPOST(w http.ResponseWriter, req *http.Request) bool {
	if req.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowd", http.StatusMethodNotAllowed)
		return false
	}

	return true
}

// IsMethodGET is a utils function that checks if the request was a GET if not it writes an http.Error into the ResponseWriter
func IsMethodGET(w http.ResponseWriter, req *http.Request) bool {
	if req.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not allowd", http.StatusMethodNotAllowed)
		return false
	}

	return true
}
