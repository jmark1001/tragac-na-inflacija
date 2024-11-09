package models

type File struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Name     string `json:"name"`
    Path     string `json:"path"`
    FileType string `json:"file_type"`
}