package service

import (
	"odime-api/internal/queue"
	"odime-api/internal/repo"
	"odime-api/pkg/models"
)

type FileService struct {
	repo  repo.FileRepository
	queue queue.MessageQueue
}

func NewFileService(repo repo.FileRepository, queue queue.MessageQueue) *FileService {
	return &FileService{repo: repo, queue: queue}
}

func (s *FileService) GetFiles() ([]models.File, error) {
	return s.repo.GetFiles()
}

func (s *FileService) UploadFile(file models.File) error {
	if err := s.repo.SaveFile(file); err != nil {
		return err
	}
	println("File saved to db!")
	return s.queue.Publish(file)
}
