package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite3 driver
)

// Initializes the SQLite database, ensures required tables exist,
// and returns a ready-to-use *sql.DB connection
// Called once at application startup; if anything fails, the app should not continue.
func initDB() *sql.DB {

	// Open the SQLite database file, fail if it cannot be opened
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}

	// Schema definition:
	// - guest information
	// - notifications: tracks which billing reminders have already been sent
	// 		- UNIQUE constraint on (guest_id, period_number) to guarantee that
	//        same billing period cannot be emailed twice
	schema := `
	CREATE TABLE IF NOT EXISTS guests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		room_number TEXT NOT NULL,
		daily_rate INTEGER NOT NULL,
		check_in_date DATE NOT NULL,
		contact TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		guest_id INTEGER NOT NULL,
		
		period_number INTEGER NOT NULL,
		sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE (guest_id, period_number)
	);`

	// Execute schema creation, fail if the database is unusable
	// App should exit and log the error
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}

	// Return the database connection
	return db
}
