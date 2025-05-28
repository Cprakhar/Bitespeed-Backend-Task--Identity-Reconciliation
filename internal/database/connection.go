package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB          *sql.DB
	ContactRepo *ContactRepository
)

func InitDB() error {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./contacts.db"
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Database connected successfully")
	if err := createTables(); err != nil {
		return err
	}

	ContactRepo = NewContactRepository(DB)

	return nil
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS contacts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        phone_number TEXT,
        email TEXT,
        linked_id INTEGER,
        link_precedence TEXT NOT NULL CHECK(link_precedence IN ('primary', 'secondary')),
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        deleted_at DATETIME,
        FOREIGN KEY (linked_id) REFERENCES contacts(id)
    );

    CREATE INDEX IF NOT EXISTS idx_phone ON contacts(phone_number) WHERE deleted_at IS NULL;
    CREATE INDEX IF NOT EXISTS idx_email ON contacts(email) WHERE deleted_at IS NULL;
    CREATE INDEX IF NOT EXISTS idx_linked_id ON contacts(linked_id) WHERE deleted_at IS NULL;
    `
	_, err := DB.Exec(query)
	if err != nil {
		log.Printf("Error creating tables: %v", err)
		return err
	}

	log.Println("Database tables created successfully")
	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
