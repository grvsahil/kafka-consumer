package kafkaconsumer

import (
	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/consumer"
)

func (a *Application) BuildConsumer() (*consumer.Consumer, error) {
	consumer, err := consumer.NewConsumer(a.log, a.services.paymentSvc, bootstrapServer)
	return consumer, err
}
