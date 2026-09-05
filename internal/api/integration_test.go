package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/config"
	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/payload"
	"github.com/Ashish-Barmaiya/torus-demo-backends/internal/simulation"
)

func newIntegrationServer(t *testing.T, service string) *httptest.Server {
	t.Helper()

	server := NewServer(config.Config{
		Service:  service,
		Instance: service + "-a",
		Port:     3000,
	})

	return httptest.NewServer(server.Handler())
}

func TestIntegrationHealth(t *testing.T) {
	ts := newIntegrationServer(t, "users")
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	if len(body) == 0 {
		t.Fatal("expected non-empty health response")
	}
}

func TestIntegrationUsersGet(t *testing.T) {
	ts := newIntegrationServer(t, "users")
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + "/api/v1/users/1")
	if err != nil {
		t.Fatalf("GET /api/v1/users/1: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	if len(body) == 0 {
		t.Fatal("expected non-empty user response")
	}
}

func TestIntegrationOrdersGet(t *testing.T) {
	ts := newIntegrationServer(t, "orders")
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + "/api/v1/orders/1")
	if err != nil {
		t.Fatalf("GET /api/v1/orders/1: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	if len(body) == 0 {
		t.Fatal("expected non-empty order response")
	}
}

func TestIntegrationResponsePayloadSize(t *testing.T) {
	ts := newIntegrationServer(t, "users")
	defer ts.Close()

	tests := []struct {
		name       string
		size       string
		wantStatus int
		wantBytes  int
	}{
		{
			name:       "1kb",
			size:       "1kb",
			wantStatus: http.StatusOK,
			wantBytes:  1 << 10,
		},
		{
			name:       "64kb",
			size:       "64kb",
			wantStatus: http.StatusOK,
			wantBytes:  64 << 10,
		},
		{
			name:       "4mb",
			size:       "4mb",
			wantStatus: http.StatusOK,
			wantBytes:  4 << 20,
		},
		{
			name:       "0b",
			size:       "0b",
			wantStatus: http.StatusOK,
			wantBytes:  0,
		},
		{
			name:       "unsupported",
			size:       "2mb",
			wantStatus: http.StatusBadRequest,
			wantBytes:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(
				http.MethodGet,
				ts.URL+"/api/v1/users/1",
				nil,
			)
			if err != nil {
				t.Fatalf("create request: %v", err)
			}

			req.Header.Set(payload.ResponseSizeHeader, tt.size)

			resp, err := ts.Client().Do(req)
			if err != nil {
				t.Fatalf("request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf(
					"status = %d, want %d",
					resp.StatusCode,
					tt.wantStatus,
				)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("read response body: %v", err)
			}

			if tt.wantBytes >= 0 && len(body) != tt.wantBytes {
				t.Fatalf(
					"body size = %d, want %d",
					len(body),
					tt.wantBytes,
				)
			}

			if tt.wantBytes == -1 && len(body) == 0 {
				t.Fatal("expected error response body")
			}
		})
	}
}

func TestIntegrationSlowSimulation(t *testing.T) {
	ts := newIntegrationServer(t, "users")
	defer ts.Close()

	req, err := http.NewRequest(
		http.MethodGet,
		ts.URL+"/api/v1/users/1",
		nil,
	)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	req.Header.Set(simulation.ModeHeader, string(simulation.ModeSlow))
	req.Header.Set(
		simulation.DelayHeader,
		strconv.Itoa(int(250)),
	)

	start := time.Now()

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if elapsed < 250*time.Millisecond {
		t.Fatalf(
			"request completed too quickly: %s",
			elapsed,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	if len(body) == 0 {
		t.Fatal("expected normal response body")
	}
}

func TestIntegrationErrorSimulation(t *testing.T) {
	ts := newIntegrationServer(t, "users")
	defer ts.Close()

	req, err := http.NewRequest(
		http.MethodGet,
		ts.URL+"/api/v1/users/1",
		nil,
	)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	req.Header.Set(simulation.ModeHeader, string(simulation.ModeError))

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf(
			"status = %d, want %d",
			resp.StatusCode,
			http.StatusServiceUnavailable,
		)
	}
}

func TestIntegrationHealthIgnoresSimulation(t *testing.T) {
	ts := newIntegrationServer(t, "users")
	defer ts.Close()

	req, err := http.NewRequest(
		http.MethodGet,
		ts.URL+"/health",
		nil,
	)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	req.Header.Set(simulation.ModeHeader, string(simulation.ModeSlow))
	req.Header.Set(
		simulation.DelayHeader,
		strconv.Itoa(int(2_000)),
	)

	start := time.Now()

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if elapsed >= 100*time.Millisecond {
		t.Fatalf(
			"health request was delayed by simulation: %s",
			elapsed,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	if len(body) == 0 {
		t.Fatal("expected health response body")
	}
}

func TestIntegrationSimulationAndPayloadSize(t *testing.T) {
	ts := newIntegrationServer(t, "orders")
	defer ts.Close()

	req, err := http.NewRequest(
		http.MethodGet,
		ts.URL+"/api/v1/orders/1",
		nil,
	)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	req.Header.Set(simulation.ModeHeader, string(simulation.ModeSlow))
	req.Header.Set(
		simulation.DelayHeader,
		strconv.Itoa(int(250)),
	)
	req.Header.Set(payload.ResponseSizeHeader, "1kb")

	start := time.Now()

	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	if elapsed < 250*time.Millisecond {
		t.Fatalf(
			"request completed too quickly: %s",
			elapsed,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	if len(body) != 1<<10 {
		t.Fatalf(
			"body size = %d, want %d",
			len(body),
			1<<10,
		)
	}
}
