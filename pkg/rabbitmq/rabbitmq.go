package rabbitmq

import (
	"github.com/streadway/amqp"
)

func NewRabbitMQConnection() (*amqp.Connection, error) {
	return amqp.Dial("amqp://guest:guest@rabbitmq:5672/") // Change localhost to rabbitmq
}
