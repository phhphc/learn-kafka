package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

var input = make(chan string)

func init() {
	go func() {
		var message string
		for {
			fmt.Print(">> ")
			fmt.Scanln(&message)
			input <- message
		}
	}()
}

/* write message using high level API */
func main() {
	chm := make(chan []kafka.Message)
	wg := sync.WaitGroup{}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go sendMessage(&wg, chm)
	go getMessage(ctx, &wg, chm)

	func() {
		<-ctx.Done()
		fmt.Println("\nClosing server...")
		wg.Wait()
	}()
}

func getMessage(ctx context.Context, wg *sync.WaitGroup, chm chan []kafka.Message) {
	wg.Add(1)
	defer wg.Done()

	messages := []kafka.Message{}
	for {
		select {
		case value := <-input:
			messages = append(messages, kafka.Message{Value: []byte(value)})
		case <-time.After(time.Second):
			chm <- messages
			messages = []kafka.Message{}
		case <-ctx.Done():
			chm <- messages
			close(chm)
			return
		}
	}
}

func sendMessage(wg *sync.WaitGroup, chm chan []kafka.Message) {
	wg.Add(1)
	defer wg.Done()

	w := kafka.Writer{
		Addr:        kafka.TCP("localhost:9092"),
		Topic:       "first-topic",
		Balancer:    &kafka.RoundRobin{},
		Compression: kafka.Snappy,
	}
	defer func() {
		if err := w.Close(); err != nil {
			log.Fatal("Fail to close writer:", err)
		}
	}()

	for messages := range chm {
		err := w.WriteMessages(context.Background(), messages...)
		if err != nil {
			log.Fatal("Fail to write message:", err)
		}
	}
}
