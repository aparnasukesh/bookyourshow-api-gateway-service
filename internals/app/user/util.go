package user

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"

	"math/rand"

	"github.com/aparnasukesh/inter-communication/user_admin"
)

func BuildGetUserProfile(res *user_admin.GetProfileResponse) (*UserProfileDetails, error) {
	return &UserProfileDetails{
		Username:    res.ProfileDetails.Username,
		FirstName:   res.ProfileDetails.FirstName,
		LastName:    res.ProfileDetails.LastName,
		PhoneNumber: res.ProfileDetails.Phone,
		DateOfBirth: res.ProfileDetails.DateOfBirth,
		Gender:      res.ProfileDetails.Gender,
	}, nil
}

func ExtractErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	errMsg := err.Error()
	if index := strings.Index(errMsg, "desc = "); index != -1 {
		return errMsg[index+len("desc = "):]
	}
	return errMsg
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewUpgrader(ctx *gin.Context) (conn *websocket.Conn, err error) {
	conn, err = upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	return
}

type RabbitMQQueueChannel struct {
	ch *amqp.Channel
	*amqp.Queue
}

func RabbitMQQueue(conn *amqp.Connection, queueName string) (*RabbitMQQueueChannel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &RabbitMQQueueChannel{
		ch:    ch,
		Queue: &q,
	}, nil
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func setupReplyQueue(queue *RabbitMQQueueChannel) (amqp.Queue, error) {
	return queue.ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
}

func sendMessage(correlationID string, queue *RabbitMQQueueChannel, body []byte, replyQueue amqp.Queue) error {
	return queue.ch.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			Body:          body,
			ReplyTo:       replyQueue.Name,
			CorrelationId: correlationID,
		},
	)
}

func waitForResponse(queue *RabbitMQQueueChannel, replyQueue amqp.Queue, correlationID string) (*amqp.Delivery, error) {
	messages, err := queue.ch.Consume(
		replyQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for msg := range messages {
		if msg.CorrelationId == correlationID {
			return &msg, nil
		}
	}
	return nil, fmt.Errorf("response with correlation ID %s not found", correlationID)
}
