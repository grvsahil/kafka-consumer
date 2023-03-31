package kafkaconsumer

import (
	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/consumer"
)

func (a *Application) BuildConsumer() *consumer.Consumer  {
	consumer := consumer.NewConsumer(a.services.paymentSvc)
	return consumer
}
