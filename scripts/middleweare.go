package scripts

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"habba/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)



func Tokengeneration(u models.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", fmt.Errorf("JWT_SECRET not set")
	}
	claims := Claims{
		Email: u.Email,
		Roll:  u.RollNumber,
		Name:  u.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return token.SignedString([]byte(jwtSecret))

}

func HashPassword(password, secret string) (string, error) {
    hash := sha256.New()
    hash.Write([]byte(password + secret))
    return hex.EncodeToString(hash.Sum(nil)), nil
}


type Claims struct {
	Email string `json:"email"`
	Roll  string    `json:"roll"`
	Name  string `json:"name"`
	jwt.RegisteredClaims
}


func TokenVerifymail(tokenStr string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET not set")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		
		return nil, fmt.Errorf("invalid or expired token")
	}

	return claims, nil
}