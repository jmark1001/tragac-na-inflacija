package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"log"
	"net/http"
	"odime-api/internal/service"
	"odime-api/pkg/models"
	"time"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (h *FileHandler) GetFiles() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		files, err := h.fileService.GetFiles()
		if err != nil {
			http.Error(w, "Error fetching files", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	})
}

func (h *FileHandler) SaveAndPublish() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var fileRequest struct {
			FilePath string `json:"file_path"` // Expecting file_path in the request
		}

		// Decode the incoming JSON body containing file path
		err := json.NewDecoder(r.Body).Decode(&fileRequest)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		file := models.File{
			ReceiptID:         uuid.New().ID(),
			Path:              fileRequest.FilePath,
			Status:            "pending",
			UploadedTimestamp: time.Now().Unix(),
		}

		if err := h.fileService.SaveAndPublish(file); err != nil {
			http.Error(w, "Error uploading file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(map[string]string{"status": "file uploaded"})
		if err != nil {
			log.Fatalf("Error uploading file: %s", err.Error())
			return
		}
	})
}
