package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

/* read message using high level Reader, print message then commit offset */
func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"34.142.253.26:9092"},
		GroupID:  "usr1",
		Topic:    "topic-B",
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

		fmt.Printf("%v Message at offset %d: %s = %s\n", time.Now().UnixMilli(), m.Offset, string(m.Key), string(m.Value))

		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatal("fail to commit message:", err)
		}
	}
}

// /* Read message using high level Reader */
// func main() {
// 	r := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:   []string{"localhost:9092"},
// 		Partition: 0,
// 		Topic:     "first-topic",
// 		MinBytes:  10e3,
// 		MaxBytes:  10e6,
// 	})
// 	r.SetOffset(-1)

// 	for {
// 		m, err := r.ReadMessage(context.Background())
// 		if err != nil {
// 			break
// 		}
// 		fmt.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
// 	}

// 	if err := r.Close(); err != nil {
// 		log.Fatal("Fail to close reader:", err)
// 	}
// }

// /* read message using low level API */
// func main() {
// conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "first-topic", 0)
// 	if err != nil {
// 		log.Fatal("Fail to dial leader:", err)
// 	}
// 	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
// 	batch := conn.ReadBatch(1e3, 1e6)

// 	b := make([]byte, 1e3)
// 	for {
// 		n, err := batch.Read(b)
// 		if err != nil {
// 			break
// 		}
// 		fmt.Println(string(b[:n]))
// 	}

// 	if err = batch.Close(); err != nil {
// 		log.Fatal("Fail to close batch:", err)
// 	}

// 	if err = conn.Close(); err != nil {
// 		log.Fatal("Fail to close connection:", err)
// 	}
// }
