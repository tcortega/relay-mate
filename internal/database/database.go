package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB opens (and if needed creates) the SQLite database at the given DSN.
// It also creates the necessary tables if they do not exist.
func InitDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := createContactsTable(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// createContactsTable ensures the "contacts" table exists.
func createContactsTable(db *sql.DB) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS contacts (
        jid TEXT PRIMARY KEY,
        last_starter_sent TEXT
    );`
	_, err := db.Exec(createTableSQL)
	return err
}
