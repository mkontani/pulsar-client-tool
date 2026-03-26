package consumer

import (
	"context"
	"fmt"
	"io"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/mkontani/pulsar-client-tool/internal/client"
	"github.com/mkontani/pulsar-client-tool/internal/output"
)

type Options struct {
	Topic            string
	SubscriptionName string
	SubscriptionType pulsar.SubscriptionType
	NumMessages      int
}

// ParseSubscriptionType parses a string into a pulsar.SubscriptionType.
func ParseSubscriptionType(s string) (pulsar.SubscriptionType, error) {
	switch s {
	case "Exclusive", "exclusive":
		return pulsar.Exclusive, nil
	case "Shared", "shared":
		return pulsar.Shared, nil
	case "Failover", "failover":
		return pulsar.Failover, nil
	case "KeyShared", "key_shared", "keyshared":
		return pulsar.KeyShared, nil
	default:
		return 0, fmt.Errorf("unknown subscription type %q (use: Exclusive, Shared, Failover, KeyShared)", s)
	}
}

// Run subscribes to the topic and consumes messages, writing each to w.
func Run(ctx context.Context, c client.PulsarClient, opts Options, w io.Writer, outputFmt string) error {
	consumer, err := c.Subscribe(pulsar.ConsumerOptions{
		Topic:            opts.Topic,
		SubscriptionName: opts.SubscriptionName,
		Type:             opts.SubscriptionType,
	})
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}
	defer consumer.Close()

	count := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		msg, err := consumer.Receive(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("receive: %w", err)
		}

		consumer.Ack(msg)

		info := output.MessageInfo{
			Topic:      msg.Topic(),
			Key:        msg.Key(),
			Payload:    string(msg.Payload()),
			MessageID:  fmt.Sprintf("%v", msg.ID()),
			Properties: msg.Properties(),
			Timestamp:  msg.PublishTime(),
		}
		if err := output.FormatMessage(w, info, outputFmt); err != nil {
			return fmt.Errorf("format output: %w", err)
		}

		count++
		if opts.NumMessages > 0 && count >= opts.NumMessages {
			return nil
		}
	}
}
