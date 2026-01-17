package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin" // Import Gin, a HTTP web framework used for routing and middleware
)

// Deletes a guest from the system + associated notifications
func deleteGuest(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get guest ID from URL parameter
		id := c.Param("id")

		// Remove all notification records tied to this guest if any exist, checking for errors
		_, err := db.Exec("DELETE FROM notifications WHERE guest_id = ?", id)
		if err != nil {
			c.JSON(500, gin.H{"error": "db error"})
			return
		}

		// Remove the guest record itself, checking for errors
		res, err := db.Exec("DELETE FROM guests WHERE id = ?", id)
		if err != nil {
			c.JSON(500, gin.H{"error": "db error"})
			return
		}

		// Check if any row was actually deleted, if not, return 404 meaning not found
		rows, _ := res.RowsAffected()
		if rows == 0 {
			c.JSON(404, gin.H{"error": "guest not found"})
			return
		}

		c.Status(204) // 204 success: deleted, no content
	}
}

/*
Returns all guests with calculated weeks stayed
  - LEFT JOIN keeps guests even if they have no notifications
  - GROUP BY prevents duplicate rows caused by the join
*/
func getGuests(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Query to get all guests, checking for errors
		rows, err := db.Query(`
			SELECT g.id, g.name, g.room_number, g.daily_rate,
			       g.check_in_date, g.contact
			FROM guests g
			LEFT JOIN notifications n ON g.id = n.guest_id
			GROUP BY g.id
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		defer rows.Close() // Ensure rows are closed after processing

		// Retrive current time for weeks stayed calculation
		now := time.Now().UTC()
		var guests []Guest

		// Iterate through query results
		for rows.Next() {
			var g Guest

			// Scan row data into guest struct, skipping on error
			err := rows.Scan(
				&g.ID, &g.Name, &g.RoomNumber, &g.DailyRate,
				&g.CheckInDate, &g.Contact,
			)
			if err != nil {
				continue
			}

			g.WeeksStayed = weeksStayed(g.CheckInDate, now) // Calculate weeks stayed
			guests = append(guests, g)                      // Add guest to the list
		}

		c.JSON(http.StatusOK, guests) // Send guest list as JSON
		c.Status(http.StatusOK)       // 200 success: retrieved guests successfully
	}
}

// Creates a new guest entry.
func createGuest(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Input structure for binding JSON
		var input struct {
			Name        string `json:"name"`
			Contact     string `json:"contact"`
			RoomNumber  string `json:"room_number"`
			DailyRate   int    `json:"daily_rate"`
			CheckInDate string `json:"check_in_date"`
		}

		// Parse JSON input to the structure
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		// Convert date string (YYYY-MM-DD) into time.Time
		checkIn, err := time.Parse("2006-01-02", input.CheckInDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
			return
		}

		// Insert new guest into the database, checking for errors
		_, err = db.Exec(
			`INSERT INTO guests (name, room_number, daily_rate, check_in_date, contact)
			 VALUES (?, ?, ?, ?, ?)`,
			input.Name,
			input.RoomNumber,
			input.DailyRate,
			checkIn.UTC(),
			input.Contact,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}

		c.Status(http.StatusCreated) // 201 success: created
		runAgent()                   // trigger agent run to pick up new guest that may need notifications
	}
}

// Calculates full weeks between check-in and now and returns it as an integer
func weeksStayed(checkIn, now time.Time) int {
	// If check-in is in the future, return 0 weeks
	if now.Before(checkIn) {
		return 0
	}
	days := int(now.Sub(checkIn).Hours() / 24) // Convert duration to days
	weeks := days / 7                          // Calculate full weeks

	return weeks
}
