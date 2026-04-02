package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mkontani/pulsar-client-tool/internal/client"
	"github.com/mkontani/pulsar-client-tool/internal/producer"
	"github.com/spf13/cobra"
)

var produceCmd = &cobra.Command{
	Use:   "produce",
	Short: "Send messages to a Pulsar topic",
	Long: `Send messages to a Pulsar topic from a string, file, or stdin.

Examples:
  # Send a single message
  pulsar-client-tool produce -t my-topic -m "hello world"

  # Send from stdin
  echo "hello" | pulsar-client-tool produce -t my-topic

  # Send from file (one message per line)
  pulsar-client-tool produce -t my-topic -f messages.txt

  # Send with message key and properties
  pulsar-client-tool produce -t my-topic -m "hello" -k my-key -p env=prod -p version=1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := ConfigFromContext(cmd.Context())

		c, err := client.New(cfg)
		if err != nil {
			return fmt.Errorf("connect to Pulsar: %w", err)
		}
		defer c.Close()

		topic, _ := cmd.Flags().GetString("topic")
		message, _ := cmd.Flags().GetString("message")
		file, _ := cmd.Flags().GetString("file")
		key, _ := cmd.Flags().GetString("key")
		props, _ := cmd.Flags().GetStringSlice("property")
		separator, _ := cmd.Flags().GetString("separator")
		numMessages, _ := cmd.Flags().GetInt("num-messages")
		rate, _ := cmd.Flags().GetFloat64("rate")
		deliverAfter, _ := cmd.Flags().GetDuration("deliver-after")
		deliverAtStr, _ := cmd.Flags().GetString("deliver-at")

		var deliverAt time.Time
		if deliverAtStr != "" {
			var err error
			deliverAt, err = time.Parse(time.RFC3339, deliverAtStr)
			if err != nil {
				return fmt.Errorf("invalid --deliver-at format, expected RFC3339 (e.g. 2024-01-01T00:00:00Z): %w", err)
			}
		}

		if deliverAfter > 0 && !deliverAt.IsZero() {
			return fmt.Errorf("--deliver-after and --deliver-at are mutually exclusive")
		}

		properties := make(map[string]string)
		for _, p := range props {
			parts := strings.SplitN(p, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid property format %q, expected key=value", p)
			}
			properties[parts[0]] = parts[1]
		}

		opts := producer.Options{
			Topic:        topic,
			Key:          key,
			Properties:   properties,
			NumMessages:  numMessages,
			Rate:         rate,
			DeliverAfter: deliverAfter,
			DeliverAt:    deliverAt,
			Separator:    separator,
		}

		// Determine input source
		var input *os.File
		if file != "" {
			f, err := os.Open(file)
			if err != nil {
				return fmt.Errorf("open file: %w", err)
			}
			defer f.Close()
			input = f
		} else if message == "" {
			// Check if stdin is a pipe
			stat, _ := os.Stdin.Stat()
			if stat.Mode()&os.ModeCharDevice == 0 {
				input = os.Stdin
			}
		}

		return producer.Run(cmd.Context(), c, opts, message, input, cmd.OutOrStdout(), cfg.OutputFmt)
	},
}

func init() {
	f := produceCmd.Flags()
	f.StringP("topic", "t", "", "topic to produce to (required)")
	f.StringP("message", "m", "", "message content to send")
	f.StringP("file", "f", "", "file to read messages from (one per line)")
	f.StringP("key", "k", "", "message key")
	f.StringSliceP("property", "p", nil, "message property (key=value, repeatable)")
	f.IntP("num-messages", "n", 1, "number of times to send the message")
	f.StringP("separator", "d", "", `message delimiter for file/stdin input (default: newline). Use "---" or blank line separator with ""`)
	f.Float64("rate", 0, "messages per second rate limit (0=unlimited)")
	f.Duration("deliver-after", 0, "delay message delivery by duration (e.g. 10s, 5m)")
	f.String("deliver-at", "", "deliver message at specific time (RFC3339, e.g. 2024-01-01T00:00:00Z)")
	_ = produceCmd.MarkFlagRequired("topic")
	rootCmd.AddCommand(produceCmd)
}
