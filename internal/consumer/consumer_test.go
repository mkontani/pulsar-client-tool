package consumer

import (
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"
)

func TestParseInitialPosition(t *testing.T) {
	tests := []struct {
		input   string
		want    pulsar.SubscriptionInitialPosition
		wantErr bool
	}{
		{"latest", pulsar.SubscriptionPositionLatest, false},
		{"Latest", pulsar.SubscriptionPositionLatest, false},
		{"earliest", pulsar.SubscriptionPositionEarliest, false},
		{"Earliest", pulsar.SubscriptionPositionEarliest, false},
		{"invalid", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseInitialPosition(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInitialPosition(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseInitialPosition(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseSubscriptionType(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"Exclusive", false},
		{"exclusive", false},
		{"Shared", false},
		{"shared", false},
		{"Failover", false},
		{"failover", false},
		{"KeyShared", false},
		{"key_shared", false},
		{"keyshared", false},
		{"invalid", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := ParseSubscriptionType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSubscriptionType(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
