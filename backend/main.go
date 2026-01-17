package main

import (
	"time"

	// "github.com/gin-contrib/cors" // CORS middleware for Gin for local testing
	"github.com/gin-gonic/gin" // Import Gin, a HTTP web framework used for routing and middleware
)

func main() {

	// Runs the billing notification agent immediately on startup, then once every 24 hours
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			runAgent() // Execute billing notification logic
			<-ticker.C // Block until the next 24-hour tick
		}
	}()

	db := initDB() // Initialize the SQLite database

	// Initialize Gin router with logging and recovery middleware
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Serve static files and the main HTML file for the frontend
	r.Static("/static", "../frontend")
	r.GET("/", func(c *gin.Context) {
		c.File("../frontend/index.html")
	})

	// -------- CORS: FOR LOCAL TESTING ONLY ---------
	// Uncomment ONLY when running frontend on a different origin (localhost:5500)
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5500"},
	// 	AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Content-Type"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: false,
	// 	MaxAge:           12 * time.Hour,
	// }))
	// -----------------------------------------------

	// ---------------- API ROUTES -------------------
	r.POST("/guests", createGuest(db))       // Create a new guest record
	r.GET("/guests", getGuests(db))          // Retrieve all guest records
	r.DELETE("/guests/:id", deleteGuest(db)) // Delete a guest record by ID
	// -----------------------------------------------

	r.Run(":8080") // Start the Gin server on port 8080
}
