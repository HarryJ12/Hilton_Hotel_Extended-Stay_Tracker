package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func deleteGuest(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		_, err := db.Exec("DELETE FROM notifications WHERE guest_id = ?", id)
		if err != nil {
			c.JSON(500, gin.H{"error": "db error"})
			return
		}

		res, err := db.Exec("DELETE FROM guests WHERE id = ?", id)
		if err != nil {
			c.JSON(500, gin.H{"error": "db error"})
			return
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			c.JSON(404, gin.H{"error": "guest not found"})
			return
		}

		c.Status(204)
	}
}

func getGuests(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		defer rows.Close()

		now := time.Now().UTC()
		var guests []Guest

		for rows.Next() {
			var g Guest

			err := rows.Scan(
				&g.ID, &g.Name, &g.RoomNumber, &g.DailyRate,
				&g.CheckInDate, &g.Contact,
			)
			if err != nil {
				continue
			}

			g.WeeksStayed = weeksStayed(g.CheckInDate, now)

			guests = append(guests, g)
		}

		c.JSON(http.StatusOK, guests)

	}
}

func createGuest(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name        string `json:"name"`
			Contact     string `json:"contact"`
			RoomNumber  string `json:"room_number"`
			DailyRate   int    `json:"daily_rate"`
			CheckInDate string `json:"check_in_date"`
		}

		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		checkIn, err := time.Parse("2006-01-02", input.CheckInDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"})
			return
		}

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

		c.Status(http.StatusCreated)
		runAgent()
	}
}

func weeksStayed(checkIn, now time.Time) int {
	if now.Before(checkIn) {
		return 0
	}

	days := int(now.Sub(checkIn).Hours() / 24)
	weeks := days / 7

	return weeks
}
