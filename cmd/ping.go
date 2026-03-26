package cmd

import (
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/mkontani/pulsar-client-tool/internal/client"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Test connectivity to Pulsar",
	Long: `Test connectivity to a Pulsar cluster by verifying broker responsiveness.

Creates a temporary producer on a non-persistent topic to ensure a full
protocol handshake with the broker, then reports the connection latency.

Examples:
  # Test default connection
  pulsar-client-tool ping

  # Test specific service URL
  pulsar-client-tool ping --service-url pulsar://broker:6650`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := ConfigFromContext(cmd.Context())

		start := time.Now()

		c, err := client.New(cfg)
		if err != nil {
			return fmt.Errorf("connect to Pulsar: %w", err)
		}
		defer c.Close()

		// Verify broker responsiveness by creating and immediately closing
		// a producer on a non-persistent topic (no data stored)
		p, err := c.CreateProducer(pulsar.ProducerOptions{
			Topic: "non-persistent://public/default/__pulsar_client_tool_ping",
		})
		if err != nil {
			return fmt.Errorf("broker handshake: %w", err)
		}
		p.Close()

		elapsed := time.Since(start).Round(time.Millisecond)
		fmt.Fprintf(cmd.OutOrStdout(), "OK - connected to %s (%s)\n", cfg.ServiceURL, elapsed)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
