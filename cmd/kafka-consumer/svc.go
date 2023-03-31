package kafkaconsumer

import (
	"database/sql"

	paymentService "github.com/grvsahil/golang-kafka/kafka-consumer/internal/service/payment"
	paymentModel "github.com/grvsahil/golang-kafka/kafka-consumer/internal/service/payment/model"
	paymentRepo "github.com/grvsahil/golang-kafka/kafka-consumer/internal/service/payment/repository"

	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/config"
)

type services struct {
	paymentSvc paymentModel.Service
}

type repos struct {
	paymentRepo paymentModel.Repository
}

func buildServices(cfg *config.Config, db *sql.DB) *services {
	svc := &services{}
	repos := &repos{}
	repos.buildRepos(db)
	svc.buildPaymentService(repos)

	return svc
}

func (r *repos) buildRepos(db *sql.DB) {
	r.paymentRepo = paymentRepo.NewRepository(db)
}

func (s *services) buildPaymentService(repo *repos) {
	s.paymentSvc = paymentService.NewService(repo.paymentRepo)
}
