package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"odime-api/internal/service"
	"odime-api/pkg/models"
	"os"
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
		// Parse the multipart form containing the file
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
			return
		}

		// Retrieve the file from the request
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving the file", http.StatusBadRequest)
			return
		}
		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		// Define the path to save the file in the container's /data directory
		filePath := "/data/" + handler.Filename

		// Create the file in the /data directory
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Could not save file", http.StatusInternalServerError)
			return
		}
		defer func(dst *os.File) {
			err := dst.Close()
			if err != nil {

			}
		}(dst)

		// Copy the uploaded file to the destination
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Error saving file", http.StatusInternalServerError)
			return
		}

		fileRecord := models.File{
			ReceiptID:         int64(uuid.New().ID()),
			SharedPath:        filePath,
			Path:              filePath,
			Status:            "pending",
			UploadedTimestamp: time.Now().Unix(),
		}

		if err := h.fileService.SaveAndPublish(&fileRecord); err != nil {
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
