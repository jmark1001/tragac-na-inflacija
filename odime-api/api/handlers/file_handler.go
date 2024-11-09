package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"odime-api/internal/service"
	"odime-api/pkg/models"
	"path/filepath"
	"strings"
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

func (h *FileHandler) UploadFile() http.Handler {
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

		// Extract file name and extension from the file path
		fileName := filepath.Base(fileRequest.FilePath)
		fileExtension := strings.ToLower(filepath.Ext(fileName))

		// Create a file model instance (without the actual file)
		file := models.File{
			Name:     fileName,
			Path:     fileRequest.FilePath,
			FileType: fileExtension,
		}

		if err := h.fileService.UploadFile(file); err != nil {
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
