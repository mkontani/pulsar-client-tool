package config

import (
	"fmt"
	"os"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

type Config struct {
	ServiceURL string
	AuthToken  string
	TLSCert    string
	TLSKey     string
	TLSCA      string
	Timeout    time.Duration
	OutputFmt  string
}

func (c *Config) Validate() error {
	if c.ServiceURL == "" {
		return fmt.Errorf("service URL is required (--service-url or PULSAR_SERVICE_URL)")
	}
	if c.OutputFmt != "text" && c.OutputFmt != "json" {
		return fmt.Errorf("output format must be 'text' or 'json', got %q", c.OutputFmt)
	}
	return nil
}

func (c *Config) ClientOptions() pulsar.ClientOptions {
	opts := pulsar.ClientOptions{
		URL:               c.ServiceURL,
		OperationTimeout:  c.Timeout,
		ConnectionTimeout: c.Timeout,
	}

	if c.AuthToken != "" {
		opts.Authentication = pulsar.NewAuthenticationToken(c.AuthToken)
	}

	if c.TLSCA != "" {
		opts.TLSTrustCertsFilePath = c.TLSCA
	}
	if c.TLSCert != "" && c.TLSKey != "" {
		opts.Authentication = pulsar.NewAuthenticationTLS(c.TLSCert, c.TLSKey)
	}

	return opts
}

// EnvOrDefault returns the environment variable value if set, otherwise the default.
func EnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
