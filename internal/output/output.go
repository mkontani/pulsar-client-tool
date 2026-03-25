package output

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type MessageInfo struct {
	Topic      string            `json:"topic"`
	Key        string            `json:"key,omitempty"`
	Payload    string            `json:"payload"`
	MessageID  string            `json:"messageId"`
	Properties map[string]string `json:"properties,omitempty"`
	Timestamp  time.Time         `json:"publishTime"`
}

type ProduceResult struct {
	Topic     string `json:"topic"`
	MessageID string `json:"messageId"`
}

func FormatMessage(w io.Writer, msg MessageInfo, format string) error {
	if format == "json" {
		return json.NewEncoder(w).Encode(msg)
	}

	fmt.Fprintf(w, "--- Message [%s] ---\n", msg.MessageID)
	fmt.Fprintf(w, "Topic: %s\n", msg.Topic)
	if msg.Key != "" {
		fmt.Fprintf(w, "Key:   %s\n", msg.Key)
	}
	fmt.Fprintf(w, "Time:  %s\n", msg.Timestamp.Format(time.RFC3339))
	if len(msg.Properties) > 0 {
		fmt.Fprintf(w, "Properties:\n")
		for k, v := range msg.Properties {
			fmt.Fprintf(w, "  %s: %s\n", k, v)
		}
	}
	fmt.Fprintf(w, "Payload:\n%s\n\n", msg.Payload)
	return nil
}

func FormatProduceConfirmation(w io.Writer, result ProduceResult, format string) error {
	if format == "json" {
		return json.NewEncoder(w).Encode(result)
	}
	fmt.Fprintf(w, "Message sent to %s (ID: %s)\n", result.Topic, result.MessageID)
	return nil
}
