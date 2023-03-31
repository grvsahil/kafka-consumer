package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/log"
	paymentModel "github.com/grvsahil/golang-kafka/kafka-consumer/internal/service/payment/model"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	paymentService paymentModel.Service
}

func NewConsumer(paymentService paymentModel.Service) *Consumer {
	return &Consumer{
		paymentService: paymentService,
	}
}

func (con Consumer) Consume(ctx context.Context, log log.Logger, bootstrapServer string, topics []string) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          "group-1",
		"auto.offset.reset": "smallest"})
	if err != nil {
		log.Errorf("failed to consume %v", err)
		return err
	}
	defer consumer.Close()

	err = consumer.SubscribeTopics(topics, nil)
	if err != nil {
		log.Errorf("failed to subscribe topic %v", err)
		return err
	}

	run := true
	for run == true {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			record, err := bytesToPayment(e.Value)
			if err != nil {
				log.Errorf("failed to parse bytes into payment %v", err)
			}
			err = con.paymentService.Add(ctx, log, []*paymentModel.Payment{&record})
			if err != nil {
				log.Errorf("failed to add record %v", err)
			}
		case kafka.Error:
			log.Errorf("Error: %v\n", e)
			run = false
		}
	}

	return errors.New("consumer stopped!")
}

func bytesToPayment(b []byte) (paymentModel.Payment, error) {
	var p paymentModel.Payment
	err := json.Unmarshal(b, &p)
	return p, err
}
