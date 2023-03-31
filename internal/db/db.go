package db

import (
	"fmt"

	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/config"
	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/log"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	database = "mysql"
)

func NewDB(config *config.Config, log log.Logger) (*sql.DB, error) {
	db, err := sql.Open(database, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name))
	if err != nil {
		log.Errorf("error connecting db %v", err)
		return nil, fmt.Errorf("error connecting db %v", err)
	}

	return db, nil
}
