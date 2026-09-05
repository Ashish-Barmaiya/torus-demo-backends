package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		service  string
		instance string
		port     string
		wantErr  bool
		wantPort int
	}{
		{
			name:     "valid configuration",
			service:  "users",
			instance: "users-a",
			port:     "3000",
			wantPort: 3000,
		},
		{
			name:     "default port",
			service:  "users",
			instance: "users-a",
			port:     "",
			wantPort: 3000,
		},
		{
			name:     "missing service",
			service:  "",
			instance: "users-a",
			port:     "3000",
			wantErr:  true,
		},
		{
			name:     "missing instance",
			service:  "users",
			instance: "",
			port:     "3000",
			wantErr:  true,
		},
		{
			name:     "invalid port",
			service:  "users",
			instance: "users-a",
			port:     "invalid",
			wantErr:  true,
		},
		{
			name:     "port too low",
			service:  "users",
			instance: "users-a",
			port:     "0",
			wantErr:  true,
		},
		{
			name:     "port too high",
			service:  "users",
			instance: "users-a",
			port:     "65536",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SERVICE", tt.service)
			t.Setenv("INSTANCE", tt.instance)
			t.Setenv("PORT", tt.port)

			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if cfg.Service != tt.service {
				t.Fatalf("expected service %q, got %q", tt.service, cfg.Service)
			}

			if cfg.Instance != tt.instance {
				t.Fatalf(
					"expected instance %q, got %q",
					tt.instance,
					cfg.Instance,
				)
			}

			if cfg.Port != tt.wantPort {
				t.Fatalf("expected port %d, got %d", tt.wantPort, cfg.Port)
			}
		})
	}
}

func TestLoadDoesNotDependOnAmbientEnvironment(t *testing.T) {
	for _, key := range []string{"SERVICE", "INSTANCE", "PORT"} {
		_ = os.Unsetenv(key)
	}

	_, err := Load()

	if err == nil {
		t.Fatal("expected error when required environment is missing")
	}
}
