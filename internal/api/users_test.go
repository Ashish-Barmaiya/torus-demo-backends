package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/config"
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
