package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestFormatMessage_Text(t *testing.T) {
	var buf bytes.Buffer
	msg := MessageInfo{
		Topic:     "persistent://public/default/test",
		Key:       "key-1",
		Payload:   "hello world",
		MessageID: "1:0:0:0",
		Properties: map[string]string{
			"env": "prod",
		},
		Timestamp: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	err := FormatMessage(&buf, msg, "text")
	if err != nil {
		t.Fatalf("FormatMessage() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "1:0:0:0") {
		t.Errorf("output should contain message ID")
	}
	if !strings.Contains(output, "hello world") {
		t.Errorf("output should contain payload")
	}
	if !strings.Contains(output, "key-1") {
		t.Errorf("output should contain key")
	}
	if !strings.Contains(output, "env: prod") {
		t.Errorf("output should contain properties")
	}
}

func TestFormatMessage_JSON(t *testing.T) {
	var buf bytes.Buffer
	msg := MessageInfo{
		Topic:     "test-topic",
		Payload:   "hello",
		MessageID: "1:0:0:0",
		Timestamp: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	err := FormatMessage(&buf, msg, "json")
	if err != nil {
		t.Fatalf("FormatMessage() error = %v", err)
	}

	var decoded MessageInfo
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if decoded.Payload != "hello" {
		t.Errorf("payload = %q, want %q", decoded.Payload, "hello")
	}
}

func TestFormatProduceConfirmation_Text(t *testing.T) {
	var buf bytes.Buffer
	result := ProduceResult{
		Topic:     "my-topic",
		MessageID: "1:0:0:0",
	}

	err := FormatProduceConfirmation(&buf, result, "text")
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "my-topic") {
		t.Errorf("output should contain topic")
	}
	if !strings.Contains(output, "1:0:0:0") {
		t.Errorf("output should contain message ID")
	}
}

func TestFormatProduceConfirmation_JSON(t *testing.T) {
	var buf bytes.Buffer
	result := ProduceResult{
		Topic:     "my-topic",
		MessageID: "1:0:0:0",
	}

	err := FormatProduceConfirmation(&buf, result, "json")
	if err != nil {
		t.Fatalf("error = %v", err)
	}

	var decoded ProduceResult
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if decoded.Topic != "my-topic" {
		t.Errorf("topic = %q, want %q", decoded.Topic, "my-topic")
	}
}
