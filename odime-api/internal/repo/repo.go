package repo

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"odime-api/pkg/models"
)

type FileRepository struct {
	db *sql.DB
}

func NewRepository(host string, port int, user string, password string, dbName string) (*FileRepository, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS files (
		id SERIAL PRIMARY KEY,
		receipt_id	BIGINT UNIQUE NOT NULL,
		shared_path TEXT NOT NULL,
		path TEXT NOT NULL,
		status VARCHAR(100),
    	uploaded_timestamp BIGINT,
    	processed_timestamp BIGINT
	)`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		receipt_id	BIGINT UNIQUE NOT NULL,
		category VARCHAR(100),
		amount DECIMAL(20, 2),
    	timestamp BIGINT
	)`)
	if err != nil {
		return nil, err
	}

	return &FileRepository{db: db}, nil
}

func (r *FileRepository) GetFiles() ([]models.File, error) {
	rows, err := r.db.Query("SELECT id, name, path, file_type FROM files")
	if err != nil {
		return nil, fmt.Errorf("failed to get files: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Slice to hold the result
	var files []models.File

	// Loop through the rows and scan into File structs
	for rows.Next() {
		var file models.File
		if err := rows.Scan(&file.ID, &file.ReceiptID, &file.Path, &file.Status); err != nil {
			return nil, fmt.Errorf("failed to scan file row: %v", err)
		}
		files = append(files, file)
	}

	// Check for errors after looping through rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %v", err)
	}

	return files, nil
}

func (r *FileRepository) SaveFile(file *models.File) error {
	query := `INSERT INTO files (receipt_id, shared_path, path, status, uploaded_timestamp)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, file.ReceiptID, file.SharedPath, file.Path, file.Status, file.UploadedTimestamp).Scan(&file.ID)
	if err != nil {
		return fmt.Errorf("failed to insert file %s: %v", file.SharedPath, err)
	}
	return nil
}

func (r *FileRepository) UpdateFile(file models.File) error {
	query := `UPDATE files SET status=$1, processed_timestamp=$2 WHERE receipt_id=$3`
	_, err := r.db.Exec(query, file.Status, file.ProcessedTimestamp, file.ReceiptID)
	if err != nil {
		return fmt.Errorf("failed to update file %s: %v", file.Path, err)
	}
	return nil
}

func (r *FileRepository) SaveExpense(expense models.Expense) error {
	query := `INSERT INTO expenses (receipt_id, category, amount, timestamp) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, expense.ReceiptID, expense.Category, expense.Amount, expense.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert expense %d: %v", expense.ReceiptID, err)
	}
	return nil
}

func (r *FileRepository) Close() {
	if err := r.db.Close(); err != nil {
		log.Println("Error closing the database connection: ", err)
	}
}
