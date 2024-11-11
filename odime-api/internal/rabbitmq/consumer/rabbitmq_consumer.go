package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"odime-api/internal/service"
	"odime-api/pkg/models"
	"time"
)

type MessageConsumer interface {
	Consume() error
}

type RabbitConsumer struct {
	channel     *amqp.Channel
	queue       amqp.Queue
	fileService service.FileService
}

func NewRabbitMQConsumer(host string, port int, user string, password string, queueName string, service service.FileService) (*RabbitConsumer, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", user, password, host, port)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}
	// Declare the rabbitmq
	queue, err := channel.QueueDeclare(
		queueName, // rabbitmq name
		true,      // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to declare a rabbitmq: %v", err)
	}

	return &RabbitConsumer{
		channel:     channel,
		queue:       queue,
		fileService: service,
	}, nil
}

func (c *RabbitConsumer) Consume() error {
	msgs, err := c.channel.Consume(
		c.queue.Name, // rabbitmq name
		"",           // consumer tag
		true,         // auto-acknowledge
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %v", err)
	}

	for msg := range msgs {
		if err := processMessage(c, msg.Body); err != nil {
			log.Printf("Error processing message: %v", err)
		}
	}

	return nil
}

func processMessage(c *RabbitConsumer, body []byte) error {
	var msg models.FileMessage

	err := json.Unmarshal(body, &msg)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return err
	}
	file := models.File{
		ReceiptID:          msg.ReceiptID,
		Path:               msg.Path,
		Status:             msg.Status,
		ProcessedTimestamp: time.Now().Unix(),
	}

	expense := models.Expense{
		ReceiptID: msg.ReceiptID,
		Category:  msg.Data.Category,
		Amount:    float32(msg.Data.Amount),
		Timestamp: time.Now().Unix(),
	}

	if err := c.fileService.ProcessConsumedFile(file, &expense); err != nil {
		return err
	}
	return nil
}

func (c *RabbitConsumer) Close() {
	if err := c.channel.Close(); err != nil {
		log.Printf("Error closing the channel: %v", err)
	}
}
