package router

//login---------will do using env

import (
	"fmt"
	"habba/models"
	"habba/scripts"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)
type VerifyRequest struct {
    UTR    string `json:"utr" binding:"required"`
   
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


func AdminRoutes(r *gin.Engine,db *gorm.DB)*gin.Engine {
	admin := r.Group("/admin")
	{
	admin.GET("/", func(c *gin.Context) {
    var data []models.Booking

    // Fetch all bookings
    if err := db.Find(&data).Error; err != nil {
        c.JSON(500, gin.H{
            "error": err.Error(),
        })
        return
    }
	total := len(data)
pendingVerify := 0
totalEntered := 0
totalcoupans:=0

for _, d := range data {
    if !d.IsVerified {
        pendingVerify=pendingVerify+d.Kannadigas+d.NonKannadigas
    }
    if d.Status == "present" {
        totalEntered=totalEntered+d.NonKannadigas+d.Kannadigas;
    }
	totalcoupans=totalcoupans+d.Kannadigas+d.NonKannadigas;
	
}

    c.JSON(200, gin.H{
        "message": "Hello Admin!",
        "data":    data,
		"totalcoupons" :totalcoupans,
		"total":total,
		"pendingVerify":pendingVerify,
		"totalEntered":totalEntered,
    })
})
admin.POST("/verifybook",func(c*gin.Context){
 var req VerifyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }

    var booking models.Booking
    if err := db.Where("utr = ?", req.UTR).First(&booking).Error; err != nil {
        c.JSON(404, gin.H{"error": "Booking not found"})
        return
    }

     booking.IsVerified=true
 if err := db.Save(&booking).Error; err != nil {
     c.JSON(500, gin.H{"error": "Failed to update booking"})
     return
    }
    scripts.EmailInvitation(booking.Email,booking)

    c.JSON(200, gin.H{
        "message": "Booking updated successfully",
        "data":    booking,
    })


})
jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	admin.POST("/adminauth", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

	
		adminName := os.Getenv("Adminname")
		adminPass := os.Getenv("Adminpass")

		if req.Username != adminName || req.Password != adminPass {
         
			c.JSON(502, gin.H{"error": "Invalid credentials"})
			return
		}
 
 

		// create JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.Username,
			"exp":      time.Now().Add(time.Hour * 1).Unix(), // expires in 1 hr
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
            fmt.Println(err.Error())
			c.JSON(400, gin.H{"error": "Failed to create token"})
			return
		}

		c.JSON(200, gin.H{"token": tokenString})
	})

	admin.POST("/verifyadmintoken", func(c *gin.Context) {
		
		var Token = c.GetHeader("Authorization")

		// parse token
		parsedToken, err := jwt.Parse(Token, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !parsedToken.Valid {
			c.JSON(502, gin.H{"error": "Invalid token"})
			return
		}

		c.JSON(200, gin.H{"message": "Token valid"})
	})
	admin.PUT("/statuscount/:utr", func(ctx *gin.Context) {
    utr := ctx.Param("utr")


    var body struct {
        StatusCount int `json:"statuscount"`
    }
    if err := ctx.ShouldBindJSON(&body); err != nil {
        ctx.JSON(400, gin.H{"error": "invalid request"})
        return
    }
    result := db.Model(&models.Booking{}).Where("utr = ?", utr).Update("status_count", body.StatusCount)
    if result.Error != nil {
		fmt.Println(result.Error)
        ctx.JSON(500, gin.H{"error": result.Error.Error()})
        return
    }

    if result.RowsAffected == 0 {
        ctx.JSON(404, gin.H{"error": "user not found"})
        return
    }

    ctx.JSON(200, gin.H{
        "message":      "statuscount updated",
        "utr":          utr,
        "statuscount":  body.StatusCount,
    })
})

admin.POST("/enter",func(c*gin.Context){
 var req VerifyRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid request"})
        return
    }

    var booking models.Booking
    if err := db.Where("utr = ?", req.UTR).First(&booking).Error; err != nil {
        c.JSON(404, gin.H{"error": "Booking not found"})
        return
    }

 booking.Status = "Present"
 if err := db.Save(&booking).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to update booking"})
        return
    }

  

    c.JSON(200, gin.H{
        "message": "Booking updated successfully",
        "data":    booking,
    })


})

	
		


		//stats




	}
	return r
}
