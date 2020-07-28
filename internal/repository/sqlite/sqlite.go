package sqlite

import (
	"database/sql"
)

// SQLiteRepository
type SQLiteRepository struct {
	db *sql.DB
}

// New creates new instance of SQLiteRepository
func New(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

// InitDB creates DB from scratch if it does not exist yet
func InitDB(db *sql.DB) error {
	stmt := `
	CREATE TABLE IF NOT EXISTS message (id INTEGER PRIMARY KEY AUTOINCREMENT, message TEXT, create_at INTEGER);
	`
	_, err := db.Exec(stmt)
	return err
}
