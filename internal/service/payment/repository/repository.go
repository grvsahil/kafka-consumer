package repository

import (
	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/log"
	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/service/payment/model"

	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Get(ctx context.Context, log log.Logger) ([]*model.Payment, error) {
	rows, err := r.db.Query("SELECT name, amount FROM payments")
	if err != nil {
		log.Errorf("failed to get records from db: ", err)
		return nil, err
	}
	defer rows.Close()

	payments := make([]*model.Payment, 0)
	for rows.Next() {
		var payment model.Payment
		if err := rows.Scan(&payment.Name, &payment.Amount); err != nil {
			log.Errorf("failed to get records from db: ", err)
			return nil, err
		}
		payments = append(payments, &payment)
	}
	if err := rows.Err(); err != nil {
		log.Errorf("failed to get records from db: ", err)
		return nil, err
	}
	return payments, nil
}

func (r *Repository) Add(ctx context.Context, log log.Logger, records []*model.Payment) error {
	stmt, err := r.db.Prepare("INSERT INTO payments (name, amount) VALUES (?, ?)")
	if err != nil {
		log.Errorf("failed to add records into db %v", err)
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		_, err = stmt.Exec(record.Name, record.Amount)
		if err != nil {
			log.Errorf("failed to add records into db %v", err)
			return err
		}
	}
	return nil
}
