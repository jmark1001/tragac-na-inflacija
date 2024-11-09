package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"odime-api/api/handlers" // Importing handlers package
	"odime-api/internal/queue"
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

	rabbitMQ, err := queue.NewRabbitMQPublisher(cfg.MQHost, cfg.MQPort, cfg.MQUser, cfg.MQPassword, cfg.MQQueueName)
	defer rabbitMQ.Close()

	fileService := service.NewFileService(*fileRepo, rabbitMQ)
	fileHandler := handlers.NewFileHandler(fileService)

	// Set up HTTP server and routes
	r := mux.NewRouter()
	r.Handle("/files", fileHandler.GetFiles()).Methods("GET")
	r.Handle("/upload", fileHandler.UploadFile()).Methods("POST")

	// Start the server
	log.Println("Starting server on port: ", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cfg.ServerPort), r))
}
