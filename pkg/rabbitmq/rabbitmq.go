package rabbitmq

import (
	"github.com/streadway/amqp"
)

func NewRabbitMQConnection() (*amqp.Connection, error) {
	return amqp.Dial("amqp://guest:guest@localhost:5672/")
}
