package broker

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type msg struct {
	Name        string  `json:"name"`
	Date        int64   `json:"date"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Description string  `json:"description"`
}

func SendEvent() {
	topic := "user-events"
	host := host()
	conn, err := kafka.DialLeader(context.Background(), "tcp", host, topic, 0)
	if err != nil {
		log.Printf("failed to dial leader: %s", err)
	}
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Printf("kafka error: %s", err)
	}

	msg1 := msg{
		"username",
		1610804729,
		45.234235,
		90.34234,
		"Event name #1",
	}
	b1, err := json.Marshal(msg1)
	if err != nil {
		log.Printf("marshal error: %s", err)
	}

	_, err = conn.WriteMessages(
		kafka.Message{Value: b1},
		// kafka.Message{Value: []byte("two!")},
		// kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Printf("failed to write messages: %s", err)
	}

	if err := conn.Close(); err != nil {
		log.Printf("failed to close writer: %s", err)
	}

}

func host() string {
	host := os.Getenv("KAFKA_HOST")

	if host != "" {
		return host
	}

	return "localhost:9092"
}
