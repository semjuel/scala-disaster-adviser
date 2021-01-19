package broker

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Name        string  `json:"name"`
	Date        int64   `json:"date"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Description string  `json:"description"`
}

func SendEvent(msg Message) {
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

	b, err := json.Marshal(msg)
	if err != nil {
		log.Printf("marshal error: %s", err)
	}

	_, err = conn.WriteMessages(
		kafka.Message{Value: b},
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
