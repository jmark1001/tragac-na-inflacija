package models

type Data struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

type FileMessage struct {
	ReceiptID uint32 `json:"receipt_id"`
	Path      string `json:"path"`
	Status    string `json:"status"`
	Data      Data   `json:"data"`
}
