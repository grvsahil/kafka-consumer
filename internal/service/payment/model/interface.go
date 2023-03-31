package model

import (
	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/log"

	"context"
)

type Repository interface {
	Get(ctx context.Context, log log.Logger) ([]*Payment, error)
	Add(ctx context.Context, log log.Logger, records []*Payment) error
}

type Service interface {
	Get(ctx context.Context, log log.Logger) ([]*Payment, error)
	Add(ctx context.Context, log log.Logger, records []*Payment) error
}
