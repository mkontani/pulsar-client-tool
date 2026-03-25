package consumer

import "testing"

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
