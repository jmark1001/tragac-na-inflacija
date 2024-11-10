package models

type File struct {
	ID                 uint   `json:"id" gorm:"primaryKey"`
	ReceiptID          uint32 `json:"receipt_id" gorm:"unique" gorm:"not null"`
	Path               string `json:"path"`
	Status             string `json:"status"`
	UploadedTimestamp  int64  `json:"uploaded_timestamp"`
	ProcessedTimestamp int64  `json:"processed_timestamp"`
}
