package simulation

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestMiddlewareNormal(t *testing.T) {
	handler := Middleware(http.HandlerFunc(func(
		w http.ResponseWriter,
		_ *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/users/1",
		nil,
	)

	rec := httptest.NewRecorder()

	start := time.Now()
	handler.ServeHTTP(rec, req)
	elapsed := time.Since(start)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if elapsed >= 100*time.Millisecond {
		t.Fatalf("normal request took too long: %s", elapsed)
	}
}

func TestMiddlewareSlow(t *testing.T) {
	handler := Middleware(http.HandlerFunc(func(
		w http.ResponseWriter,
		_ *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/users/1",
		nil,
	)

	req.Header.Set(ModeHeader, "slow")
	req.Header.Set(DelayHeader, "250")

	rec := httptest.NewRecorder()

	start := time.Now()
	handler.ServeHTTP(rec, req)
	elapsed := time.Since(start)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if elapsed < 250*time.Millisecond {
		t.Fatalf("request completed too quickly: %s", elapsed)
	}
}

func TestMiddlewareError(t *testing.T) {
	handlerCalled := false

	handler := Middleware(http.HandlerFunc(func(
		w http.ResponseWriter,
		_ *http.Request,
	) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/users/1",
		nil,
	)

	req.Header.Set(ModeHeader, "error")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf(
			"status = %d, want %d",
			rec.Code,
			http.StatusServiceUnavailable,
		)
	}

	if handlerCalled {
		t.Fatal("next handler should not be called for error simulation")
	}
}

func TestMiddlewareHealthBypassesSimulation(t *testing.T) {
	handler := Middleware(http.HandlerFunc(func(
		w http.ResponseWriter,
		_ *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("healthy"))
	}))

	req := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	req.Header.Set(ModeHeader, "slow")
	req.Header.Set(DelayHeader, "2000")

	rec := httptest.NewRecorder()

	start := time.Now()
	handler.ServeHTTP(rec, req)
	elapsed := time.Since(start)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if elapsed >= 100*time.Millisecond {
		t.Fatalf("health check was delayed: %s", elapsed)
	}

	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if string(body) != "healthy" {
		t.Fatalf("body = %q, want %q", string(body), "healthy")
	}
}

func TestMiddlewareConcurrentRequests(t *testing.T) {
	handler := Middleware(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("normal"))
	}))

	const (
		normalRequests = 10
		slowRequests   = 10
		errorRequests  = 10
	)

	var wg sync.WaitGroup

	results := make(chan int, normalRequests+slowRequests+errorRequests)

	send := func(mode Mode, delay string) {
		defer wg.Done()

		req := httptest.NewRequest(
			http.MethodGet,
			"/api/v1/users/1",
			nil,
		)

		req.Header.Set(ModeHeader, string(mode))

		if delay != "" {
			req.Header.Set(DelayHeader, delay)
		}

		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		results <- rec.Code
	}

	wg.Add(normalRequests + slowRequests + errorRequests)

	for i := 0; i < normalRequests; i++ {
		go send(ModeNormal, "")
	}

	for i := 0; i < slowRequests; i++ {
		go send(ModeSlow, "250")
	}

	for i := 0; i < errorRequests; i++ {
		go send(ModeError, "")
	}

	wg.Wait()
	close(results)

	var statusOK int
	var statusUnavailable int

	for status := range results {
		switch status {
		case http.StatusOK:
			statusOK++

		case http.StatusServiceUnavailable:
			statusUnavailable++

		default:
			t.Fatalf("unexpected status: %d", status)
		}
	}

	if statusOK != normalRequests+slowRequests {
		t.Fatalf(
			"successful requests = %d, want %d",
			statusOK,
			normalRequests+slowRequests,
		)
	}

	if statusUnavailable != errorRequests {
		t.Fatalf(
			"error requests = %d, want %d",
			statusUnavailable,
			errorRequests,
		)
	}
}
