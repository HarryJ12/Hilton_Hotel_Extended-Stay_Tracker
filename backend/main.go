package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	go func() {
		ticker := time.NewTicker(24 * time.Hour)

		defer ticker.Stop()

		for {
			runAgent()
			<-ticker.C
		}
	}()

	db := initDB()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5500"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/guests", createGuest(db))
	r.GET("/guests", getGuests(db))
	r.DELETE("/guests/:id", deleteGuest(db))

	r.Run(":8080")
}
