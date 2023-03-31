package kafkaconsumer

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/config"
	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/consumer"
	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/db"
	"github.com/grvsahil/golang-kafka/kafka-consumer/internal/log"
)

const (
	bootstrapServer = "localhost:9092"
	topic           = "payment"
	server          = "localhost:8089"
)

type Application struct {
	db       *sql.DB
	log      log.Logger
	cfg      *config.Config
	services *services
	router   *gin.Engine
	consumer *consumer.Consumer
}

func (a *Application) Init(ctx context.Context, configFile string) {
	log := log.New().With(ctx)
	a.log = log

	config, err := config.Load(log, configFile)
	if err != nil {
		log.Fatalf("failed to read config: %s ", err)
		return
	}
	a.cfg = config

	db, err := db.NewDB(config, log)
	if err != nil {
		log.Fatalf("error connecting db: %s ", err)
		return
	}
	a.db = db

	services := buildServices(config, db)
	a.services = services

	router := gin.Default()
	a.router = router
	a.SetupHandlers()
	go func(ctx context.Context) {
		err := a.router.Run(server)
		if err != nil {
			a.log.Fatalf("failed to start server: %s ", err)
			return
		}
	}(ctx)

	consumer := a.BuildConsumer()
	a.consumer = consumer

	err = a.consumer.Consume(ctx, a.log, bootstrapServer, []string{topic})
	if err != nil {
		log.Fatalf("error starting consumer: %s ", err)
		return
	}
}
