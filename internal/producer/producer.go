package producer

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/mkontani/pulsar-client-tool/internal/client"
	"github.com/mkontani/pulsar-client-tool/internal/output"
)

type Options struct {
	Topic        string
	Key          string
	Properties   map[string]string
	NumMessages  int
	Rate         float64
	DeliverAfter time.Duration
	DeliverAt    time.Time
	Separator    string
	Raw          bool
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
		if opts.DeliverAfter > 0 {
			msg.DeliverAfter = opts.DeliverAfter
		}
		if !opts.DeliverAt.IsZero() {
			msg.DeliverAt = opts.DeliverAt
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

	// Send from reader (stdin or file)
	if input != nil {
		if opts.Raw {
			// Read entire input as a single message
			data, err := io.ReadAll(input)
			if err != nil {
				return fmt.Errorf("read input: %w", err)
			}
			payload := bytes.TrimRight(data, "\n")
			if len(payload) == 0 {
				return fmt.Errorf("no message provided: input is empty")
			}
			return send(payload)
		}

		scanner := bufio.NewScanner(input)
		first := true

		if opts.Separator != "" {
			// Multi-line mode: accumulate lines until separator line is found
			var lines []string
			for scanner.Scan() {
				line := scanner.Text()
				if line == opts.Separator {
					if len(lines) > 0 {
						select {
						case <-ctx.Done():
							return ctx.Err()
						default:
						}
						if !first {
							throttle()
						}
						first = false
						msg := strings.Join(lines, "\n")
						if err := send([]byte(msg)); err != nil {
							return err
						}
						lines = lines[:0]
					}
				} else {
					lines = append(lines, line)
				}
			}
			// Send remaining lines as final message
			if len(lines) > 0 {
				if !first {
					throttle()
				}
				msg := strings.Join(lines, "\n")
				if err := send([]byte(msg)); err != nil {
					return err
				}
			}
		} else {
			// Default: one message per line
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
		}
		return scanner.Err()
	}

	return fmt.Errorf("no message provided: use --message, --file, or pipe via stdin")
}
