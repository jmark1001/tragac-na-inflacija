package models

type File struct {
	ID                 uint   `json:"id" gorm:"primaryKey"`
	ReceiptID          int64  `json:"receipt_id" gorm:"unique" gorm:"not null"`
	SharedPath         string `json:"shared_path" gorm:"not null"`
	Path               string `json:"path" gorm:"not null"`
	Status             string `json:"status" gorm:"not null"`
	UploadedTimestamp  int64  `json:"uploaded_timestamp"`
	ProcessedTimestamp int64  `json:"processed_timestamp"`
}
