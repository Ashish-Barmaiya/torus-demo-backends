package simulation

import "testing"

func TestParseMode(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    Mode
		wantErr bool
	}{
		{
			name:  "normal",
			value: "normal",
			want:  ModeNormal,
		},
		{
			name:  "slow",
			value: "slow",
			want:  ModeSlow,
		},
		{
			name:  "error",
			value: "error",
			want:  ModeError,
		},
		{
			name:    "offline is unsupported",
			value:   "offline",
			wantErr: true,
		},
		{
			name:    "invalid",
			value:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMode(tt.value)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("mode = %q, want %q", got, tt.want)
			}
		})
	}
}
