package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mkontani/pulsar-client-tool/internal/config"
	"github.com/spf13/cobra"
)

type contextKey string

const configKey contextKey = "config"

var (
	serviceURL string
	authToken  string
	tlsCert    string
	tlsKey     string
	tlsCA      string
	timeout    time.Duration
	outputFmt  string
)

var rootCmd = &cobra.Command{
	Use:   "pulsar-client-tool",
	Short: "A user-friendly CLI for Apache Pulsar",
	Long:  "pulsar-client-tool is a Go-based CLI for producing and consuming messages on Apache Pulsar.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg := &config.Config{
			ServiceURL: serviceURL,
			AuthToken:  authToken,
			TLSCert:    tlsCert,
			TLSKey:     tlsKey,
			TLSCA:      tlsCA,
			Timeout:    timeout,
			OutputFmt:  outputFmt,
		}
		if err := cfg.Validate(); err != nil {
			return err
		}

		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		_ = cancel // cancel is called when signal is received
		cmd.SetContext(context.WithValue(ctx, configKey, cfg))
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func ConfigFromContext(ctx context.Context) *config.Config {
	return ctx.Value(configKey).(*config.Config)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&serviceURL, "service-url", config.EnvOrDefault("PULSAR_SERVICE_URL", "pulsar://localhost:6650"), "Pulsar service URL")
	pf.StringVar(&authToken, "auth-token", config.EnvOrDefault("PULSAR_AUTH_TOKEN", ""), "authentication token")
	pf.StringVar(&tlsCert, "tls-cert", "", "TLS client certificate path")
	pf.StringVar(&tlsKey, "tls-key", "", "TLS client key path")
	pf.StringVar(&tlsCA, "tls-ca", "", "TLS trusted CA certificate path")
	pf.DurationVar(&timeout, "timeout", 30*time.Second, "connection and operation timeout")
	pf.StringVarP(&outputFmt, "output", "o", "text", "output format (text, json)")
}
