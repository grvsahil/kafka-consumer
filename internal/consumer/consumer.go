package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/log"
	paymentModel "github.com/grvsahil/golang-kafka/kafka-consumer/internal/service/payment/model"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	log            log.Logger
	paymentService paymentModel.Service
	Consumers      []*kafka.Consumer
}

func NewConsumer(log log.Logger, paymentService paymentModel.Service, bootstrapServer string) (*Consumer, error) {
	consumer1, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          "group-1",
		"auto.offset.reset": "smallest"})
	if err != nil {
		log.Errorf("failed to initialize consumer1 %v", err)
		return nil, err
	}

	consumer2, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          "group-2",
		"auto.offset.reset": "smallest"})
	if err != nil {
		log.Errorf("failed to initialize consumer2 %v", err)
		return nil, err
	}

	consumer3, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          "group-3",
		"auto.offset.reset": "smallest"})
	if err != nil {
		log.Errorf("failed to initialize consumer3 %v", err)
		return nil, err
	}

	consumers := []*kafka.Consumer{consumer1, consumer2, consumer3}

	return &Consumer{
		log:            log,
		paymentService: paymentService,
		Consumers:      consumers,
	}, nil
}

func (con Consumer) Consume(ctx context.Context, topics []string) error {
	done := make(chan bool)

	for i, consumer := range con.Consumers {
		consumer.SubscribeTopics(topics, nil)
		go con.consumeMessages(ctx, consumer, i+1, done)
	}

	// Handle termination signals
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Wait for termination signal or error
	select {
	case <-done:
		con.log.Info("Consumer error occurred. Exiting...")
		return errors.New("consumer stopped!")
	case <-sigterm:
		con.log.Info("Termination signal received. Exiting...")
	}

	return nil
}

func bytesToPayment(b []byte) ([]*paymentModel.Payment, error) {
	var p []*paymentModel.Payment
	err := json.Unmarshal(b, &p)
	return p, err
}

func (con Consumer) consumeMessages(ctx context.Context, consumer *kafka.Consumer, consumerIndex int, done chan bool) {
	con.log.Infof("Consumer %d initialized", consumerIndex)
	defer close(done)

	for {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			con.log.Infof("Consumed message from consumer %d", consumerIndex)
			record, err := bytesToPayment(e.Value)
			if err != nil {
				con.log.Errorf("failed to parse bytes into payment %v", err)
			}
			err = con.paymentService.Add(ctx, con.log, record)
			if err != nil {
				con.log.Errorf("failed to add record %v", err)
			}
		case kafka.Error:
			con.log.Errorf("Error: %v\n", e)
			done <- true
		}
	}
}
