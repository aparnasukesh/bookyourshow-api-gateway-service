package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// func NewRabbitMQConnection() (*amqp.Connection, error) {
// 	return amqp.Dial("amqp://guest:guest@rabbitmq:5672/") // Change localhost to rabbitmq
// }

//	func NewRabbitMQConnection() (*amqp.Connection, error) {
//		return amqp.Dial("amqp://user:password@rabbitmq:5672/") // Use updated credentials
//	}

//	func NewRabbitMQConnection() (*amqp.Connection, error) {
//		return amqp.Dial("amqp://admin:admin123@rabbitmq:5672/") // Updated creadentials
//	}
func NewRabbitMQConnection() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial("amqp://admin:admin123@rabbitmq:5672/")
		if err == nil {
			log.Println("✅ Connected to RabbitMQ")
			return conn, nil
		}
		log.Printf("❌ RabbitMQ not ready, retrying in 5s... (%d/10) - %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after retries: %v", err)
}
