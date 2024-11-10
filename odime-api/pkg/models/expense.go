package models

type Expense struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	ReceiptID uint32  `json:"receipt_id"`
	Category  string  `json:"category"`
	Amount    float32 `json:"amount"`
	Timestamp int64   `json:"timestamp"`
}