package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

func GenerateTokenFromUserEmail(email string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString
}

func ParseToken(tokenString string) *jwt.Token {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if ok == false {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return token
}
