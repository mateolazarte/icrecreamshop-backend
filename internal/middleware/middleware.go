package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"icecreamshop/internal/auth"
	"icecreamshop/internal/messageErrors"
	"icecreamshop/internal/storage"
	"net/http"
	"time"
)

type Middleware struct {
	Store storage.Storage
}

func NewMiddleware(store storage.Storage) *Middleware {
	return &Middleware{Store: store}
}

// CheckIfNotLoggedIn checks if there is NOT an user logged in. Otherwise, aborts.
func (middleware *Middleware) CheckIfNotLoggedIn(c *gin.Context) {
	cookie, err := c.Cookie("Authorization")
	if err == nil && cookie != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.AlreadyLoggedIn})
		c.Abort()
		return
	}

	c.Next()
}

// AuthenticateUser authenticates if an user is logged in.
func (middleware *Middleware) AuthenticateUser(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := auth.ParseToken(tokenString)
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if tokenIsExpired(claims) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := middleware.Store.GetUserByEmail(claims["sub"].(string))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user-email", user.Email)
	c.Set("user-id", user.ID)
	c.Next()
}

// AuthenticateAdmin authenticates if an admin is logged in.
func (middleware *Middleware) AuthenticateAdmin(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := auth.ParseToken(tokenString)
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if tokenIsExpired(claims) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := middleware.Store.GetUserByEmail(claims["sub"].(string))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !user.IsAdmin() {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Set("user-email", user.Email)
	c.Set("user-id", user.ID)
	c.Next()
}

// AuthenticateDeliveryDriver authenticates if a delivery driver is logged in.
func (middleware *Middleware) AuthenticateDeliveryDriver(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token := auth.ParseToken(tokenString)
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := middleware.Store.GetUserByEmail(claims["sub"].(string))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !user.IsAdmin() && !user.IsDeliveryDriver() {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Set("user-email", user.Email)
	c.Set("user-id", user.ID)
	c.Next()
}

func tokenIsExpired(claims jwt.MapClaims) bool {
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		return true
	}
	return false
}
