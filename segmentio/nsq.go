package main

import (
	"context"
	"fmt"

	"github.com/segmentio/nsq-go"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	// Producer
	g.Go(func() error {
		return Producer()
	})
	if err := g.Wait(); err != nil {
		fmt.Println(err.Error())
	}

	// Consumer
	g.Go(func() error {
		return Consumer()
	})
	if err := g.Wait(); err != nil {
		fmt.Println(err.Error())
	}
}

func Producer() error {
	// Starts a new producer that publishes to the TCP endpoint of a nsqd node.
	// The producer automatically handles connections in the background.
	producer, _ := nsq.StartProducer(nsq.ProducerConfig{
		Topic:   "hello",
		Address: "localhost:4150",
	})

	// Publishes a message to the topic that this producer is configured for,
	// the method returns when the operation completes, potentially returning an
	// error if something went wrong.
	err := producer.Publish([]byte("Hello World!"))
	if err != nil {
		fmt.Println(err.Error())
	}

	// Stops the producer, all in-flight requests will be canceled and no more
	// messages can be published through this producer.
	producer.Stop()

	return nil
}

func Consumer() error {
	// Create a new consumer, looking up nsqd nodes from the listed nsqlookup
	// addresses, pulling messages from the 'world' channel of the 'hello' topic
	// with a maximum of 250 in-flight messages.
	consumer, err := nsq.StartConsumer(nsq.ConsumerConfig{
		Topic:       "hello",
		Channel:     "world",
		Address:     "localhost:4150",
		MaxInFlight: 250,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	// Consume messages, the consumer automatically connects to the nsqd nodes
	// it discovers and handles reconnections if something goes wrong.
	for msg := range consumer.Messages() {
		fmt.Println(string(msg.Body))
		msg.Finish()
		if len(consumer.Messages()) == 0 {
			break
		}
	}
	consumer.Stop()
	return err
}
