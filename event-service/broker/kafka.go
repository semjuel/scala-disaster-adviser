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
	Name    string
	Date    int64
	Lat     float64
	Lon     float64
	Summary string
}

func SendEvent() {
	topic := "user-events"
	host := host()
	conn, err := kafka.DialLeader(context.Background(), "tcp", host, topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal("kafka error:", err)
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
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}

func host() string {
	host := os.Getenv("KAFKA_HOST")

	if host != "" {
		return host
	}

	return "localhost:9092"
}
