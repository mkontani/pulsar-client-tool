package config

import (
	"os"
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				ServiceURL: "pulsar://localhost:6650",
				OutputFmt:  "text",
				Timeout:    30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "valid json output",
			cfg: Config{
				ServiceURL: "pulsar://localhost:6650",
				OutputFmt:  "json",
			},
			wantErr: false,
		},
		{
			name: "missing service URL",
			cfg: Config{
				OutputFmt: "text",
			},
			wantErr: true,
		},
		{
			name: "invalid output format",
			cfg: Config{
				ServiceURL: "pulsar://localhost:6650",
				OutputFmt:  "xml",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_ClientOptions(t *testing.T) {
	cfg := Config{
		ServiceURL: "pulsar://example.com:6650",
		AuthToken:  "my-token",
		Timeout:    10 * time.Second,
		OutputFmt:  "text",
	}

	opts := cfg.ClientOptions()
	if opts.URL != "pulsar://example.com:6650" {
		t.Errorf("URL = %q, want %q", opts.URL, "pulsar://example.com:6650")
	}
	if opts.OperationTimeout != 10*time.Second {
		t.Errorf("OperationTimeout = %v, want %v", opts.OperationTimeout, 10*time.Second)
	}
	if opts.Authentication == nil {
		t.Error("Authentication should not be nil when token is set")
	}
}

func TestEnvOrDefault(t *testing.T) {
	t.Run("returns env value when set", func(t *testing.T) {
		os.Setenv("TEST_PULSAR_VAR", "from-env")
		defer os.Unsetenv("TEST_PULSAR_VAR")

		got := EnvOrDefault("TEST_PULSAR_VAR", "default")
		if got != "from-env" {
			t.Errorf("EnvOrDefault() = %q, want %q", got, "from-env")
		}
	})

	t.Run("returns default when unset", func(t *testing.T) {
		got := EnvOrDefault("TEST_PULSAR_UNSET_VAR", "default")
		if got != "default" {
			t.Errorf("EnvOrDefault() = %q, want %q", got, "default")
		}
	})
}
