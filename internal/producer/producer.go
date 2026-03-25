package producer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/mkontani/pulsar-client-tool/internal/client"
	"github.com/mkontani/pulsar-client-tool/internal/output"
)

type Options struct {
	Topic       string
	Key         string
	Properties  map[string]string
	NumMessages int
	Rate        float64
}

// Run sends messages from the given input to the topic.
// If input is nil and message is empty, it returns an error.
func Run(ctx context.Context, c client.PulsarClient, opts Options, message string, input io.Reader, w io.Writer, outputFmt string) error {
	p, err := c.CreateProducer(pulsar.ProducerOptions{
		Topic: opts.Topic,
	})
	if err != nil {
		return fmt.Errorf("create producer: %w", err)
	}
	defer p.Close()

	send := func(payload []byte) error {
		msg := &pulsar.ProducerMessage{
			Payload:    payload,
			Properties: opts.Properties,
		}
		if opts.Key != "" {
			msg.Key = opts.Key
		}

		id, err := p.Send(ctx, msg)
		if err != nil {
			return fmt.Errorf("send message: %w", err)
		}

		return output.FormatProduceConfirmation(w, output.ProduceResult{
			Topic:     opts.Topic,
			MessageID: fmt.Sprintf("%v", id),
		}, outputFmt)
	}

	// Rate limiter interval
	var interval time.Duration
	if opts.Rate > 0 {
		interval = time.Duration(float64(time.Second) / opts.Rate)
	}

	throttle := func() {
		if interval > 0 {
			time.Sleep(interval)
		}
	}

	// Send literal message
	if message != "" {
		n := opts.NumMessages
		if n <= 0 {
			n = 1
		}
		for i := 0; i < n; i++ {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			if err := send([]byte(message)); err != nil {
				return err
			}
			if i < n-1 {
				throttle()
			}
		}
		return nil
	}

	// Send from reader (stdin or file), one message per line
	if input != nil {
		scanner := bufio.NewScanner(input)
		first := true
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			if !first {
				throttle()
			}
			first = false
			if err := send(scanner.Bytes()); err != nil {
				return err
			}
		}
		return scanner.Err()
	}

	return fmt.Errorf("no message provided: use --message, --file, or pipe via stdin")
}
