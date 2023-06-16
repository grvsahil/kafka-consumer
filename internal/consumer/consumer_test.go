package consumer_test

import (
	"context"
	"errors"
	"testing"

	consumer "github.com/grvsahil/golang-kafka/kafka-consumer/internal/consumer"
	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/log"
	logMock "github.com/grvsahil/golang-kafka/kafka-consumer/internal/log/mocks"
	paymentModel "github.com/grvsahil/golang-kafka/kafka-consumer/internal/service/payment/model"
	paymentMocks "github.com/grvsahil/golang-kafka/kafka-consumer/internal/service/payment/model/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	bootstrapServer = "localhost:9092"
	errService      = errors.New("payment service error")
)

func TestConsumer_Consume(t *testing.T) {
	t.Parallel()
	type fields struct {
		log             log.Logger
		paymentService  paymentModel.Service
		bootstrapServer string
	}
	type args struct {
		ctx    context.Context
		topics []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "consumer stopped",
			fields: fields{
				log: func() log.Logger {
					logger := new(logMock.Logger)
					logger.On("Infof", mock.Anything, mock.Anything)
					logger.On("Info", mock.Anything, mock.Anything)
					logger.On("Errorf", mock.Anything, mock.Anything)
					return logger
				}(),
				paymentService: func() paymentModel.Service {
					service := new(paymentMocks.Service)
					service.On("Add", mock.Anything, mock.Anything, mock.Anything).Return(nil)
					return service
				}(),
				bootstrapServer: bootstrapServer,
			},
			args: args{
				ctx: context.Background(),
				topics: func() []string {
					topics := []string{"testTopics"}
					return topics
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, err := consumer.NewConsumer(tt.fields.log, tt.fields.paymentService, tt.fields.bootstrapServer)
			if err != nil {
				t.Errorf("Failed to initialize consumer: %v", err)
			}
			if err := c.Consume(tt.args.ctx, tt.args.topics); (err != nil) != tt.wantErr {
				t.Errorf("Consume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
