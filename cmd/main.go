package main

import (
	"context"

	kafkaConsumer "github.com/grvsahil/golang-kafka/kafka-consumer/cmd/kafka-consumer"
)

const defaultConfPath = "./local.yaml"

func main() {
	application := &kafkaConsumer.Application{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application.Init(ctx, defaultConfPath)
}
