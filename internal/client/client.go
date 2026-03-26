package client

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/mkontani/pulsar-client-tool/internal/config"
)

// PulsarClient abstracts the pulsar.Client interface for testability.
type PulsarClient interface {
	CreateProducer(pulsar.ProducerOptions) (pulsar.Producer, error)
	Subscribe(pulsar.ConsumerOptions) (pulsar.Consumer, error)
	Close()
}

// New creates a new Pulsar client from the given config.
func New(cfg *config.Config) (PulsarClient, error) {
	return pulsar.NewClient(cfg.ClientOptions())
}
