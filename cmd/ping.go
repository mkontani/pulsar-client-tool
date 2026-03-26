package cmd

import (
	"fmt"

	"github.com/mkontani/pulsar-client-tool/internal/client"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Test connectivity to Pulsar",
	Long: `Test connectivity to a Pulsar cluster by establishing and closing a connection.

Examples:
  # Test default connection
  pulsar-client-tool ping

  # Test specific service URL
  pulsar-client-tool ping --service-url pulsar://broker:6650`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := ConfigFromContext(cmd.Context())

		c, err := client.New(cfg)
		if err != nil {
			return fmt.Errorf("connect to Pulsar: %w", err)
		}
		c.Close()

		fmt.Fprintln(cmd.OutOrStdout(), "OK - connected to", cfg.ServiceURL)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
