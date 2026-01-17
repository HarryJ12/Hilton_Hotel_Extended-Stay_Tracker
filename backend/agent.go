package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import the SQLite3 driver
)

// runAgent is a background task that:
// 1. Scans all guests
// 2. Computes how many weeks they have stayed
// 3. Sends a billing reminder exactly once per guest per week
func runAgent() {
	// Open the database connection, log and return on error
	// Ensure the database is closed when done
	agentDB, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Println("Agent error: cannot open db")
		return
	}
	defer agentDB.Close()

	// Set SQLite pragmas for WAL mode and busy timeout
	agentDB.Exec("PRAGMA journal_mode=WAL;")
	agentDB.Exec("PRAGMA busy_timeout = 5000;")

	now := time.Now().UTC() // Get the current UTC time

	// Query all guests from the database, log and return on error
	// Ensure rows are closed when done
	rows, err := agentDB.Query(`
		SELECT id, name, room_number, daily_rate, check_in_date, contact
		FROM guests
	`)
	if err != nil {
		log.Println("Agent error: read failed")
		return
	}
	defer rows.Close()

	// Define a candidate struct to hold guest info for notification
	type candidate struct {
		id      int
		name    string
		room    string
		rate    int
		week    int
		contact string
	}
	var toNotify []candidate

	// Iterate over each guest row to determine if they need notification
	for rows.Next() {
		var (
			id      int
			name    string
			room    string
			rate    int
			checkIn time.Time
			contact string
		)

		// Scan row into variables
		if err := rows.Scan(
			&id,
			&name,
			&room,
			&rate,
			&checkIn,
			&contact,
		); err != nil {
			continue // Skip malformed rows instead of crashing the agent
		}

		// Calculate weeks stayed and check if notification is needed
		weeks := weeksStayed(checkIn, now)
		if weeks < 1 {
			continue
		}

		// Check if notification for this week already exists
		toNotify = append(toNotify, candidate{
			id:      id,
			name:    name,
			room:    room,
			rate:    rate,
			week:    weeks,
			contact: contact,
		})
	}

	// Send notifications
	for _, g := range toNotify {
		// Insert a notification record only if it doesn't already exist, log and continue on error
		// This guarantees that notifications are sent only once per week
		res, err := agentDB.Exec(`
			INSERT OR IGNORE INTO notifications (guest_id, period_number)
			VALUES (?, ?)
		`, g.id, g.week)
		if err != nil {
			log.Println("Agent error: insert failed")
			continue
		}

		// If no row was inserted, the notification was already sent
		affected, _ := res.RowsAffected()
		if affected != 1 {
			continue
		}

		subject := "Extended-Stay Guest Billing Reminder" // Email subject

		// Email body content
		body := fmt.Sprintf(
			"Weekly Billing Reminder for: \n"+"%s\n\n"+
				"Room: %s\n"+
				"Weeks Stayed: %d\n"+
				"Daily Rate: $%d\n"+
				"Contact Information: %s\n",
			g.name,
			g.room,
			g.week,
			g.rate,
			g.contact,
		)

		// Send billing reminder email, log error if fails
		if err := sendEmail(subject, body); err != nil {
			log.Println("Agent error: email send failed")
			continue
		}
	}
}
