package utils

import (
	"net/http/httptest"
	"testing"
)

func TestSetHeaderJson(t *testing.T) {
	w := httptest.NewRecorder()

	SetHeaderJson(w)
	expectedContentType := "application/vnd.api+json"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected Content-type %q, got %q", expectedContentType, contentType)
	}
}
