package main

import (
	"habba/router"
	"habba/scripts"
	// "log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"
)

func main() {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	db := scripts.ConnectDatabase()

	


r := gin.Default()
r.Use(cors.New(cors.Config{
	AllowOrigins:     []string{"http://localhost:3000","https://kannaddaganeshutsava.vercel.app"}, 
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           12 * time.Hour,
}))
r.Use(scripts.LimitPerIP())
r = router.UserRoutes(r,db)
r =router.AdminRoutes(r,db)
	
	r.Run(":8088")

	
	select {}
}