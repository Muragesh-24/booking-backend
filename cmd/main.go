package main

import (
	"habba/router"
	"habba/scripts"
	"log"
	"os"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load(".env")
	db := scripts.ConnectDatabase()
	scripts.RabbitMQConnection()
	scripts.StartEmailVerificationConsumer()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{

			
			"https://kannaddabalag.muragesh.tech",
			"https://booking.kannadda.muragesh.tech",
			"http://localhost:3000",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(scripts.LimitPerIP())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r = router.UserRoutes(r, db)
	r = router.AdminRoutes(r, db)

	log.Printf("Server running on http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
