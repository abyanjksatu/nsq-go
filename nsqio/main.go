package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/nsqio/go-nsq"
	"golang.org/x/sync/errgroup"
)

type (
	NsqProducer interface {
		Publish(message []byte) error
	}

	NsqProducerImpl struct {
		Producer *nsq.Producer
	}

	NsqConsumer interface {
		Consume(handle func(msg *nsq.Message) error) error
	}

	NsqConsumerImpl struct {
		Consumer *nsq.Consumer
	}
)

/*
	Producer
*/

func NewNsqProducer() NsqProducer {
	// Create new config and producer with address nsq
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("127.0.0.1:4150", config)
	if err != nil {
		log.Panic(err)
	}

	return &NsqProducerImpl{
		Producer: producer,
	}
}

func (n *NsqProducerImpl) Publish(message []byte) error {
	// Publish with topic
	err := n.Producer.Publish("hello-topic", message)
	if err != nil {
		return err
	}
	return nil
}

/*
	Consumer
*/

func NewNsqConsumer() NsqConsumer {
	// Initiate Consumer with new config, topic, and channel
	decodeConfig := nsq.NewConfig()
	consumer, err := nsq.NewConsumer("hello-topic", "hello-channel", decodeConfig)
	if err != nil {
		log.Panic("Could not create consumer")
	}

	return &NsqConsumerImpl{
		Consumer: consumer,
	}
}

func (n *NsqConsumerImpl) Consume(handle func(msg *nsq.Message) error) error {
	// Consume with handle func(msg *nsq.Message) error
	// Connect to address to listen message
	n.Consumer.AddHandler(nsq.HandlerFunc(handle))
	err := n.Consumer.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Panic("Could not connect")
	}
	return err
}

/*
	Main
*/
func main() {
	// Run produer with errgroup goroutine, to add group use .Go()
	// After group, you can start run with g.Wait() which return error
	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)
	producer := NewNsqProducer()
	consumer := NewNsqConsumer()

	// Producer
	g.Go(func() error {
		return producer.Publish([]byte("Hello World"))
	})
	if err := g.Wait(); err != nil {
		fmt.Println(err.Error())
	}

	// Listener wait for signal use sync.WaitGroup{}
	// After get & process signal message, we state to w.Done()
	// Which mean the process is final and 1 way situation
	w := sync.WaitGroup{}
	w.Add(1)

	// Consumer
	err := consumer.Consume(func(message *nsq.Message) error {
		log.Printf("NSQ message received, msg: %s", string(message.Body))
		w.Done()
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	w.Wait()
}
