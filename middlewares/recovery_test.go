package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecoveryHandler_RecoversPanic(t *testing.T) {
	h := RecoveryHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	}))

	req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}

	if !strings.Contains(rr.Body.String(), "Internal Server Error") {
		t.Fatalf("body = %q, want to contain %q", rr.Body.String(), "Internal Server Error")
	}
}
