# Kafka-consumer

This application demonstrates how kafka consumer works. We are using 
[Confluent Kafka](github.com/confluentinc/confluent-kafka-go/kafka) library to implement kafka consumer.

### Quickstart

You don't need to have MYSQL installed locally as this app is using docker to run MYSQL.
#### Steps:-
1. Just clone this repository and run "docker compose up -d" to pull the required images and start database in docker container in detached mode.
2. Run the application using "go run cmd/main.go"
3. That's it the application is up and ready!

### Description 

We have used [GIN](https://github.com/gin-gonic/gin) web framework to implement REST API.

This application has consumer group called "group1" and this group of consumer consumes messages from a remote kafka producer and perists it in MYSQL database.

Currently it supports only one GET route "/records" which fetches persisted data from database.

It also has a custom logger to log different errors and events.


