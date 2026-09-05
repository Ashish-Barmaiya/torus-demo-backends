package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/config"
	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/payload"
)

func newOrdersTestServer() *Server {
	return NewServer(config.Config{
		Service:  "orders",
		Instance: "orders-test",
		Port:     3000,
	})
}

func TestOrdersGET(t *testing.T) {
	server := newOrdersTestServer()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/orders",
		nil,
	)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()

	if !strings.Contains(body, `"service":"orders"`) {
		t.Fatalf("expected orders service: %s", body)
	}

	if !strings.Contains(body, `"instance":"orders-test"`) {
		t.Fatalf("expected orders-test instance: %s", body)
	}
}

func TestOrderGET(t *testing.T) {
	server := newOrdersTestServer()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/orders/7",
		nil,
	)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), `"id":"ord_000007"`) {
		t.Fatalf("expected order 7 in response: %s", rec.Body.String())
	}
}

func TestOrderGETNotFound(t *testing.T) {
	server := newOrdersTestServer()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/orders/999",
		nil,
	)
	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestOrderPOST(t *testing.T) {
	server := newOrdersTestServer()

	body := `{
		"customer_id": "usr_000005",
		"currency": "USD",
		"total": 129900
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/orders",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	response := rec.Body.String()

	if !strings.Contains(response, `"customer_id":"usr_000005"`) {
		t.Fatalf("expected customer ID: %s", response)
	}
}

func TestOrderPOSTInvalidInput(t *testing.T) {
	server := newOrdersTestServer()

	body := `{
		"customer_id": "",
		"currency": "USD",
		"total": 129900
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/orders",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestOrderPATCH(t *testing.T) {
	server := newOrdersTestServer()

	body := `{
		"status": "completed",
		"total": 149900
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/orders/7",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	response := rec.Body.String()

	if !strings.Contains(response, `"status":"completed"`) {
		t.Fatalf("expected completed status: %s", response)
	}

	if !strings.Contains(response, `"total":149900`) {
		t.Fatalf("expected updated total: %s", response)
	}
}

func TestOrderPATCHEmptyRequest(t *testing.T) {
	server := newOrdersTestServer()

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/orders/7",
		strings.NewReader(`{}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestOrderDELETE(t *testing.T) {
	server := newOrdersTestServer()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/orders/7",
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

func TestOrderResponsePayloadSize(t *testing.T) {
	cfg := config.Config{
		Service:  "orders",
		Instance: "orders-a",
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
				"/api/v1/orders/1",
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
