package rabbitmq

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

// func NewRabbitMQConnection() (*amqp.Connection, error) {
// 	return amqp.Dial("amqp://guest:guest@rabbitmq:5672/") // Change localhost to rabbitmq
// }

// func NewRabbitMQConnection() (*amqp.Connection, error) {
// 	return amqp.Dial("amqp://user:password@rabbitmq:5672/") // Use updated credentials
// }

// func NewRabbitMQConnection() (*amqp.Connection, error) {
// 	return amqp.Dial("amqp://myuser:mypassword@rabbitmq:5672/")
// }

func NewRabbitMQConnection() (*amqp.Connection, error) {
	username := os.Getenv("RABBITMQ_USERNAME")
	password := os.Getenv("RABBITMQ_PASSWORD")

	rabbitmqURL := fmt.Sprintf("amqp://%s:%s@rabbitmq:5672/", username, password)
	return amqp.Dial(rabbitmqURL)
}
