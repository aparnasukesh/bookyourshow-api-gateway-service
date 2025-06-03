package rabbitmq

import "github.com/streadway/amqp"

// func NewRabbitMQConnection() (*amqp.Connection, error) {
// 	return amqp.Dial("amqp://guest:guest@rabbitmq:5672/") // Change localhost to rabbitmq
// }

//	func NewRabbitMQConnection() (*amqp.Connection, error) {
//		return amqp.Dial("amqp://user:password@rabbitmq:5672/") // Use updated credentials
//	}

func NewRabbitMQConnection() (*amqp.Connection, error) {
	return amqp.Dial("amqp://admin:admin123@rabbitmq:5672/")
}
