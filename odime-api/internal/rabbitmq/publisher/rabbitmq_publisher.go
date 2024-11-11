package publisher

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"odime-api/pkg/models"
)

type MessageQueue interface {
	Publish(file models.File) error
}

type RabbitPublisher struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitMQPublisher(host string, port int, user string, password string, queueName string) (*RabbitPublisher, error) {
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
		} // Close the connection properly if the channel creation fails
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	// Declare a rabbitmq (ensure it exists)
	queue, err := channel.QueueDeclare(
		queueName, // rabbitmq name
		true,      // durable (survives broker restart)
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		} // Close the connection properly if rabbitmq declaration fails
		return nil, fmt.Errorf("failed to declare a rabbitmq: %v", err)
	}

	// Return the publisher with the channel and rabbitmq
	return &RabbitPublisher{
		channel: channel,
		queue:   queue,
	}, nil
}

func (p *RabbitPublisher) Publish(file models.File) error {
	message, err := json.Marshal(file)
	if err != nil {
		return fmt.Errorf("failed to serialize file data: %v", err)
	}

	// Publish the message to RabbitMQ
	err = p.channel.Publish(
		"",           // default exchange
		p.queue.Name, // routing key (rabbitmq name)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to RabbitMQ: %v", err)
	}

	return nil
}

func (p *RabbitPublisher) Close() {
	// Close the channel properly
	if err := p.channel.Close(); err != nil {
		log.Printf("Error closing the channel: %v", err)
	}
}
