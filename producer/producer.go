package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

var input = make(chan string)
var username []byte

func init() {
	// push stdin's input to input channel
	fmt.Print("Enter your username: ")
	fmt.Scanln(&username)

	r := bufio.NewReader(os.Stdin)
	go func() {
		fmt.Print(">> ")
		for {
			text, err := r.ReadString('\n')
			if err == io.EOF {
				continue
			} else if err != nil {
				log.Fatal("Fail to read user input:", err)
			}

			input <- text[:len(text)-1]
			fmt.Print(">> ")
		}
	}()
}

func main() {
	chm := make(chan []kafka.Message, 3)
	wg := sync.WaitGroup{}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	w := kafka.Writer{
		Addr:     kafka.TCP("34.142.253.26:9092"),
		Topic:    "topic-B",
		Balancer: &kafka.RoundRobin{},

		Compression: kafka.Snappy,
	}
	defer func() {
		if err := w.Close(); err != nil {
			log.Fatal("Fail to close writer:", err)
		}
	}()

	go sendMessage(&wg, w, chm)
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
			messages = append(messages, kafka.Message{Value: []byte(value), Key: username})
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

func sendMessage(wg *sync.WaitGroup, w kafka.Writer, chm chan []kafka.Message) {
	wg.Add(1)
	defer wg.Done()

	for messages := range chm {
		err := w.WriteMessages(context.Background(), messages...)
		if err != nil {
			log.Fatal("Fail to write message:", err)
		}
	}
}
