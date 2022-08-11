package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/segmentio/kafka-go"
)

/* write message using high level API */
func main() {
	w := kafka.Writer{
		Addr:        kafka.TCP("localhost:9092"),
		Topic:       "first-topic",
		Balancer:    &kafka.RoundRobin{},
		Compression: kafka.Snappy,
	}

	wg := sync.WaitGroup{}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	func() {
		var message string

		for {
			select {
			case s := <-ch:
				fmt.Println("Got signal:", s)
				fmt.Println("Quitting...")
				return
			default:
				fmt.Print(">> ")
				fmt.Scanln(&message)

				wg.Add(1)
				go func(message string) {
					defer wg.Done()
					err := w.WriteMessages(context.Background(), kafka.Message{
						Value: []byte(message),
					})
					if err != nil {
						log.Fatal("Fail to write message:", err)
					}
				}(message)
			}
		}
	}()

	wg.Wait()
	if err := w.Close(); err != nil {
		log.Fatal("Fail to close writer:", err)
	}
}

// /* write message using low level API*/
// func main() {
// 	topic := "first-topic"
// 	partition := 0

// 	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
// 	if err != nil {
// 		log.Fatal("Fail to dial leader:", err)
// 	}
// 	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

// 	_, err = conn.WriteMessages(kafka.Message{Value: []byte("hello world")})
// 	if err != nil {
// 		log.Fatal("Fail to write message:", err)
// 	}

// 	if err := conn.Close(); err != nil {
// 		log.Fatal("Fail to close writer:", err)
// 	}
// }
