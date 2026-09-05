package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/config"
	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/payload"
)

func newUsersTestServer() *Server {
	return NewServer(config.Config{
		Service:  "users",
		Instance: "users-test",
		Port:     3000,
	})
}

func TestUsersGET(t *testing.T) {
	server := newUsersTestServer()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/users",
		nil,
	)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"service":"users"`) {
		t.Fatalf("expected users service in response: %s", rec.Body.String())
	}

	if !strings.Contains(rec.Body.String(), `"instance":"users-test"`) {
		t.Fatalf("expected users-test instance in response: %s", rec.Body.String())
	}
}

func TestUserGET(t *testing.T) {
	server := newUsersTestServer()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/users/5",
		nil,
	)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()

	if !strings.Contains(body, `"id":"usr_000005"`) {
		t.Fatalf("expected user 5 in response: %s", body)
	}
}

func TestUserGETInvalidID(t *testing.T) {
	server := newUsersTestServer()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/users/not-a-number",
		nil,
	)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUserGETNotFound(t *testing.T) {
	server := newUsersTestServer()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/users/999",
		nil,
	)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUserPOST(t *testing.T) {
	server := newUsersTestServer()

	body := `{
		"name": "Alice Johnson",
		"email": "alice@example.com",
		"plan": "pro"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/users",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"name":"Alice Johnson"`) {
		t.Fatalf("expected created user in response: %s", rec.Body.String())
	}
}

func TestUserPOSTInvalidJSON(t *testing.T) {
	server := newUsersTestServer()

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/users",
		strings.NewReader(`{broken`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUserPOSTInvalidInput(t *testing.T) {
	server := newUsersTestServer()

	body := `{
		"email": "alice@example.com",
		"plan": "pro"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/users",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUserPATCH(t *testing.T) {
	server := newUsersTestServer()

	body := `{
		"plan": "enterprise",
		"status": "suspended"
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/users/5",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	response := rec.Body.String()

	if !strings.Contains(response, `"plan":"enterprise"`) {
		t.Fatalf("expected updated plan: %s", response)
	}

	if !strings.Contains(response, `"status":"suspended"`) {
		t.Fatalf("expected updated status: %s", response)
	}
}

func TestUserPATCHEmptyRequest(t *testing.T) {
	server := newUsersTestServer()

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/users/5",
		strings.NewReader(`{}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUserDELETE(t *testing.T) {
	server := newUsersTestServer()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/users/5",
		nil,
	)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"deleted":true`) {
		t.Fatalf("expected deleted=true: %s", rec.Body.String())
	}
}

func TestUserDELETEUnsupportedMethod(t *testing.T) {
	server := newUsersTestServer()

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/users/5",
		nil,
	)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusMethodNotAllowed,
			rec.Code,
		)
	}
}

func TestUserResponsePayloadSize(t *testing.T) {
	cfg := config.Config{
		Service:  "users",
		Instance: "users-a",
		Port:     3000,
	}

	server := NewServer(cfg)

	tests := []struct {
		name       string
		header     string
		wantStatus int
		wantLength int
	}{
		{
			name:       "normal response",
			header:     "",
			wantStatus: http.StatusOK,
			wantLength: -1,
		},
		{
			name:       "1kb response",
			header:     "1kb",
			wantStatus: http.StatusOK,
			wantLength: 1 << 10,
		},
		{
			name:       "64kb response",
			header:     "64kb",
			wantStatus: http.StatusOK,
			wantLength: 64 << 10,
		},
		{
			name:       "empty response",
			header:     "0b",
			wantStatus: http.StatusOK,
			wantLength: 0,
		},
		{
			name:       "unsupported response size",
			header:     "2mb",
			wantStatus: http.StatusBadRequest,
			wantLength: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				"/api/v1/users/1",
				nil,
			)

			if tt.header != "" {
				req.Header.Set(payload.ResponseSizeHeader, tt.header)
			}

			rec := httptest.NewRecorder()

			server.Handler().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"status = %d, want %d",
					rec.Code,
					tt.wantStatus,
				)
			}

			if tt.wantLength >= 0 && rec.Body.Len() != tt.wantLength {
				t.Fatalf(
					"body length = %d, want %d",
					rec.Body.Len(),
					tt.wantLength,
				)
			}

			if tt.wantLength == -1 && rec.Body.Len() == 0 {
				t.Fatal("expected normal JSON response body")
			}
		})
	}
}
