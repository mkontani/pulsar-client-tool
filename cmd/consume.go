package cmd

import (
	"fmt"

	"github.com/mkontani/pulsar-client-tool/internal/client"
	"github.com/mkontani/pulsar-client-tool/internal/consumer"
	"github.com/spf13/cobra"
)

var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Consume messages from a Pulsar topic",
	Long: `Subscribe to a Pulsar topic and consume messages.

Examples:
  # Consume messages indefinitely
  pulsar-client-tool consume -t my-topic -s my-sub

  # Consume 10 messages then exit
  pulsar-client-tool consume -t my-topic -s my-sub -n 10

  # Consume with JSON output (pipe to jq)
  pulsar-client-tool consume -t my-topic -s my-sub -o json

  # Consume with shared subscription
  pulsar-client-tool consume -t my-topic -s my-sub -S Shared`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := ConfigFromContext(cmd.Context())

		c, err := client.New(cfg)
		if err != nil {
			return fmt.Errorf("connect to Pulsar: %w", err)
		}
		defer c.Close()

		topic, _ := cmd.Flags().GetString("topic")
		subscription, _ := cmd.Flags().GetString("subscription")
		subType, _ := cmd.Flags().GetString("subscription-type")
		numMessages, _ := cmd.Flags().GetInt("num-messages")

		st, err := consumer.ParseSubscriptionType(subType)
		if err != nil {
			return err
		}

		opts := consumer.Options{
			Topic:            topic,
			SubscriptionName: subscription,
			SubscriptionType: st,
			NumMessages:      numMessages,
		}

		return consumer.Run(cmd.Context(), c, opts, cmd.OutOrStdout(), cfg.OutputFmt)
	},
}

func init() {
	f := consumeCmd.Flags()
	f.StringP("topic", "t", "", "topic to consume from (required)")
	f.StringP("subscription", "s", "", "subscription name (required)")
	f.StringP("subscription-type", "S", "Exclusive", "subscription type (Exclusive, Shared, Failover, KeyShared)")
	f.IntP("num-messages", "n", 0, "number of messages to consume (0=unlimited)")
	_ = consumeCmd.MarkFlagRequired("topic")
	_ = consumeCmd.MarkFlagRequired("subscription")
	rootCmd.AddCommand(consumeCmd)
}
