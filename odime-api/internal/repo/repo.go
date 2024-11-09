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
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	// Ensure the files table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS files (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		path TEXT NOT NULL,
		file_type VARCHAR(10) NOT NULL
	)`)
	if err != nil {
		return nil, err
	}

	return &FileRepository{db: db}, nil
}

func (r *FileRepository) GetFiles() ([]models.File, error) {
	rows, err := r.db.Query("SELECT id, name, path, file_type FROM files")
	if err != nil {
		return nil, fmt.Errorf("failed to get files: %w", err)
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
		if err := rows.Scan(&file.ID, &file.Name, &file.Path, &file.FileType); err != nil {
			return nil, fmt.Errorf("failed to scan file row: %w", err)
		}
		files = append(files, file)
	}

	// Check for errors after looping through rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return files, nil
}

func (r *FileRepository) SaveFile(file models.File) error {
	query := `INSERT INTO files (name, path, file_type) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, file.Name, file.Path, file.FileType)
	if err != nil {
		return fmt.Errorf("failed to insert file: %s", file.Path)
	}
	return nil
}

func (r *FileRepository) Close() {
	if err := r.db.Close(); err != nil {
		log.Println("Error closing the database connection:", err)
	}
}
