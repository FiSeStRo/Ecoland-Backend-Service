package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsMethodPOST(t *testing.T) {
	w := httptest.NewRecorder()
	testReq := []struct {
		name     string
		req      *http.Request
		expected bool
	}{
		{name: "POST request",
			req:      httptest.NewRequest(http.MethodPost, "/root", strings.NewReader("testString")),
			expected: true},
		{name: "GET request",
			req:      httptest.NewRequest(http.MethodGet, "/root", strings.NewReader("testString")),
			expected: false},
	}
	for _, tr := range testReq {
		t.Run(tr.name, func(t *testing.T) {
			got := IsMethodPOST(w, tr.req)
			if got != tr.expected {
				t.Errorf("IsMethodPost(w, %q), got %v, expected %v", tr.name, got, tr.expected)
			}
		})

	}
}
