package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"odime-api/api/handlers" // Importing handlers package
	"odime-api/internal/rabbitmq/consumer"
	"odime-api/internal/rabbitmq/publisher"
	"odime-api/internal/repo"
	"odime-api/internal/service" // Importing service package
	"odime-api/pkg/config"       // Configuration package
	"strconv"
)

func main() {
	// Load configuration
	cfg, _ := config.LoadConfig()

	fileRepo, err := repo.NewRepository(cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer fileRepo.Close()

	rabbitMQPublisher, err := publisher.NewRabbitMQPublisher(cfg.MQHost, cfg.MQPort, cfg.MQUser, cfg.MQPassword, cfg.MQPendingQueue)
	if err != nil {
		log.Fatal(err)
	}
	println("publisher declared successfully")
	defer rabbitMQPublisher.Close()

	fileService := service.NewFileService(*fileRepo, rabbitMQPublisher)
	fileHandler := handlers.NewFileHandler(fileService)

	// Initialize RabbitMQ consumer
	rabbitMQConsumer, err := consumer.NewRabbitMQConsumer(cfg.MQHost, cfg.MQPort, cfg.MQUser, cfg.MQPassword, cfg.MQProcessedQueue, *fileService)
	if err != nil {
		log.Fatal(err)
	}
	println("consumer declared successfully")
	defer rabbitMQConsumer.Close()

	// Start the consumer in a separate goroutine
	go func() {
		if err := rabbitMQConsumer.Consume(); err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
	}()

	// Set up HTTP server and routes
	r := mux.NewRouter()
	r.Handle("/files", fileHandler.GetFiles()).Methods("GET")
	r.Handle("/upload", fileHandler.SaveAndPublish()).Methods("POST")

	// Start the server
	log.Println("Starting server on port: ", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.ServerPort), r))
}
