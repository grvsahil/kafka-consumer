package kafkaconsumer

import (
	"github.com/grvsahil/golang-kafka/kafka-consumer/api/v1/payment"
)

func (a *Application) SetupHandlers() {
	payment.RegisterHandlers(
		a.router,
		a.services.paymentSvc,
		a.log,
	)
}
