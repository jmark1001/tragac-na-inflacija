package queue

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
	println(connStr)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		} // Close the connection properly if the channel creation fails
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare a queue (ensure it exists)
	queue, err := channel.QueueDeclare(
		queueName, // queue name
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
		} // Close the connection properly if queue declaration fails
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Return the publisher with the channel and queue
	return &RabbitPublisher{
		channel: channel,
		queue:   queue,
	}, nil
}

func (p *RabbitPublisher) Publish(file models.File) error {
	message, err := json.Marshal(file)
	if err != nil {
		log.Printf("Failed to serialize file data: %v", err)
		return fmt.Errorf("failed to serialize file data: %w", err)
	}

	// Publish the message to RabbitMQ
	err = p.channel.Publish(
		"",           // default exchange
		p.queue.Name, // routing key (queue name)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message to RabbitMQ: %v", err)
		return fmt.Errorf("failed to publish message to RabbitMQ: %w", err)
	}

	log.Println("Message published to RabbitMQ successfully:", string(message))
	return nil
}

func (p *RabbitPublisher) Close() {
	// Close the channel properly
	if err := p.channel.Close(); err != nil {
		log.Printf("Error closing the channel: %v", err)
	}
}
