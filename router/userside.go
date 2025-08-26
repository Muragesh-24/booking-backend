package router

import (
	"fmt"
	"habba/models"
	"habba/scripts"
	// "os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.Engine,db *gorm.DB) *gin.Engine {
	user := r.Group("/user")
	{
		user.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello User!"})
		})
		//login
 user.POST("/signin", func(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		// Password string `json:"password"`
	}


	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}


	var u models.User
	if err := db.Where("email = ?", req.Email).First(&u).Error; err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	
	if !u.IsVerified {
		ctx.JSON(403, gin.H{"error": "Email not verified"})
		return
	}

	// hashSecret := os.Getenv("HASH_SECRET")
	// if hashSecret == "" {
	// 	ctx.JSON(500, gin.H{"error": "hash secret not set"})
	// 	return
	// }

	// hashedInputPwd, err := scripts.HashPassword(req.Password, hashSecret)
	// if err != nil {
	// 	ctx.JSON(500, gin.H{"error": "failed to hash password"})
	// 	return
	// }

	// if hashedInputPwd != u.Password {
	// 	ctx.JSON(401, gin.H{"error": "Invalid email or password"})
	// 	return
	// }

	token, err := scripts.Tokengeneration(u)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "signin successful",
		"token":   token,
	})
})

		
		//signup
    user.POST("/signup", func(ctx *gin.Context) {
    var u models.User
    if err := ctx.ShouldBindJSON(&u); err != nil {
        fmt.Println(err.Error())
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // hashSecret := os.Getenv("HASH_SECRET")
    // if hashSecret == "" {
    //     ctx.JSON(500, gin.H{"error": "hash secret not set"})
    //     return
    // }


    // hashedPwd, err := scripts.HashPassword(u.Password, hashSecret)
    // if err != nil {
    //     ctx.JSON(500, gin.H{"error": "failed to hash password"})
    //     return
    // }
    // u.Password = hashedPwd
    u.IsVerified = false

    
    if err := db.Create(&u).Error; err != nil {
        ctx.JSON(500, gin.H{"error": "could not save user"})
        return
    }

    token, err := scripts.Tokengeneration(u)
    if err != nil {
        ctx.JSON(500, gin.H{"error": "could not generate token"})
        return
    }
    if err := scripts.EmailVerifymail(u.Email, token); err != nil {
        fmt.Println(err.Error())
       ctx.JSON(201, gin.H{"message": "Signup successful, verification email failed,"})
        return
    }

    ctx.JSON(200, gin.H{"message": "signup successful, please verify email"})
})



//mailverify--used for token verify also
user.GET("/verify", func(ctx *gin.Context) {
    token := ctx.Query("token")
    if token == "" {
        ctx.JSON(400, gin.H{"error": "token missing"})
        return
    }

    // Parse the token -> extract user email
    claims, err := scripts.TokenVerifymail(token) // you'll implement ParseToken
    if err != nil {
        ctx.JSON(400, gin.H{"error": "invalid or expired token"})
        return
    }

    var user models.User
    if err := db.Where("email = ?", claims.Email).First(&user).Error; err != nil {
        ctx.JSON(404, gin.H{"error": "user not found"})
        return
    }

    // Mark user as verified
    user.IsVerified = true
    if err := db.Save(&user).Error; err != nil {
        ctx.JSON(500, gin.H{"error": "could not update user"})
        return
    }

   
   ctx.Redirect(302, fmt.Sprintf("https://kannaddaganeshiitk.vercel.app/auth/verified?query=%s", token))

})


user.GET("/verifytoken", func(c *gin.Context) {
    // Get the token from Authorization header
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(401, gin.H{"error": "Authorization header missing"})
        return
    }

    // Expected format: "Bearer <token>"
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        c.JSON(401, gin.H{"error": "Invalid Authorization header"})
        return
    }

    tokenStr := parts[1]

    claims, err := scripts.TokenVerifymail(tokenStr)
    if err != nil {
        c.JSON(401, gin.H{"error": "Invalid or expired token"})
        return
    }

    // Token is valid
    c.JSON(200, gin.H{
        "message": "token is valid",
        "email":   claims.Email,
        "userId":  claims.ID,
    })
})

user.GET("/mybooking",func(c *gin.Context) {

 authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(402, gin.H{"error": "Missing Authorization header"})
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == authHeader {
        c.JSON(402, gin.H{"error": "Invalid Authorization format"})
        return
    }

    claims, err := scripts.TokenVerifymail(tokenString)
    if err != nil {
        c.JSON(402, gin.H{"error": "Invalid or expired token"})
        return
    }

var bookings []models.Booking
    if err := db.Where("email = ?", claims.Email).
        Order("created_at DESC").
        Find(&bookings).Error; err != nil {
        c.JSON(400, gin.H{"error": "DB error"})
        return
    }

    // 4. Send response
    c.JSON(200, gin.H{
        "message": "Your Bookings",
        "data":    bookings,
        "total":   len(bookings),
    })
})


user.POST("/book", func(c *gin.Context) {

    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(402, gin.H{"error": "Missing Authorization header"})
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == authHeader {
        c.JSON(402, gin.H{"error": "Invalid Authorization format"})
        return
    }

    claims, err := scripts.TokenVerifymail(tokenString)
    if err != nil {
        c.JSON(402, gin.H{"error": "Invalid or expired token"})
        return
    }

    var booking models.Booking
    if err := c.ShouldBindJSON(&booking); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    booking.Email = claims.Email
    booking.CreatedAt = time.Now()
    if err := db.Create(&booking).Error; err != nil {
        fmt.Println(err.Error())
        c.JSON(401, gin.H{"error": "Failed to save booking"})
        return
    }

    c.JSON(200, gin.H{
        "message": "Booking saved successfully",
        "booking": booking,
    })
})


	}

	



//getting booked details
	return r
}
