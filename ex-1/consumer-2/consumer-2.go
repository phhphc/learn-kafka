package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

/* read message using high level Reader, print message then commit offset */
func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{":9092"},
		GroupID:  "group-1",
		Topic:    "2-part",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	ctx := context.Background()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			fmt.Print(err)
			break
		}

		fmt.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("fail to commit message:", err)
		}
	}
}
