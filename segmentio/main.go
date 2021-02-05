package main

import (
	"context"
	"fmt"

	"github.com/segmentio/nsq-go"
	"golang.org/x/sync/errgroup"
)

type (
	HelloNsq interface {
		Publisher() *nsq.Producer
		Consumer() *nsq.Consumer
		Publish([]byte) error
		Consume(func(msg *nsq.Message) error) error
	}

	HelloNsqImpl struct {
		NsqConsumer *nsq.Consumer
		NsqProducer *nsq.Producer
	}
)

func NewHelloNsq() HelloNsq {
	// Starts a new producer that publishes to the TCP endpoint of a nsqd node.
	// The producer automatically handles connections in the background.
	producer, err := nsq.NewProducer(nsq.ProducerConfig{
		Topic:   "hello",
		Address: "localhost:4150",
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	// Create a new consumer, looking up nsqd nodes from the listed nsqlookup
	// addresses, pulling messages from the 'world' channel of the 'hello' topic
	// with a maximum of 250 in-flight messages.
	consumer, err := nsq.NewConsumer(nsq.ConsumerConfig{
		Topic:       "hello",
		Channel:     "world",
		Address:     "localhost:4150",
		MaxInFlight: 250,
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	return &HelloNsqImpl{
		NsqProducer: producer,
		NsqConsumer: consumer,
	}
}

func (h *HelloNsqImpl) Publisher() *nsq.Producer {
	return h.NsqProducer
}

func (h *HelloNsqImpl) Consumer() *nsq.Consumer {
	return h.NsqConsumer
}

func (h *HelloNsqImpl) Publish(message []byte) error {
	h.NsqProducer.Start()
	// Publishes a message to the topic that this producer is configured for,
	// the method returns when the operation completes, potentially returning an
	// error if something went wrong.
	err := h.NsqProducer.Publish(message)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Stops the producer, all in-flight requests will be canceled and no more
	// messages can be published through this producer.
	h.NsqProducer.Stop()
	return nil
}

func (h *HelloNsqImpl) Consume(handler func(msg *nsq.Message) error) error {
	h.NsqConsumer.Start()
	// Consume messages, the consumer automatically connects to the nsqd nodes
	// it discovers and handles reconnections if something goes wrong.
	for msg := range h.NsqConsumer.Messages() {
		err := handler(&msg)
		if err != nil {
			return err
		}
		msg.Finish()
		if len(h.NsqConsumer.Messages()) == 0 {
			break
		}
	}
	h.NsqConsumer.Stop()
	return nil
}

func (h *HelloNsqImpl) WaitConsume(handle func(msg *nsq.Message) error) error {
	h.NsqConsumer.Start()
	// Consume messages, the consumer automatically connects to the nsqd nodes
	// it discovers and handles reconnections if something goes wrong.
	for msg := range h.NsqConsumer.Messages() {
		err := handle(&msg)
		if err != nil {
			fmt.Println(err.Error())
		}
		msg.Finish()
	}
	return nil
}

func main() {
	// Run produer with errgroup goroutine, to add group use .Go()
	// After group, you can start run with g.Wait() which return error
	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)
	helloNsq := NewHelloNsq()

	// Producer
	g.Go(func() error {
		return helloNsq.Publish([]byte("Hello World!"))
	})
	if err := g.Wait(); err != nil {
		fmt.Println(err.Error())
	}

	// Consumer
	g.Go(func() error {
		return helloNsq.Consume(func(msg *nsq.Message) error {
			fmt.Println(string(msg.Body))
			return nil
		})
	})
	if err := g.Wait(); err != nil {
		fmt.Println(err.Error())
	}
}
