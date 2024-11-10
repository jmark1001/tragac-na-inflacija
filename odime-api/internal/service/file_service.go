package service

import (
	"odime-api/internal/rabbitmq/publisher"
	"odime-api/internal/repo"
	"odime-api/pkg/models"
)

type FileService struct {
	repo  repo.FileRepository
	queue publisher.MessageQueue
}

func NewFileService(repo repo.FileRepository, queue publisher.MessageQueue) *FileService {
	return &FileService{repo: repo, queue: queue}
}

func (s *FileService) GetFiles() ([]models.File, error) {
	return s.repo.GetFiles()
}

func (s *FileService) SaveAndPublish(file models.File) error {
	if err := s.repo.SaveFile(file); err != nil {
		return err
	}
	print("file saved")
	return s.queue.Publish(file)
}

func (s *FileService) ProcessConsumedFile(file models.File, expense *models.Expense) error {
	if err := s.repo.UpdateFile(file); err != nil {
		return err
	}
	if expense != nil {
		if err := s.repo.SaveExpense(*expense); err != nil {
			return err
		}
	}
	return nil
}
