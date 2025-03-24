package user

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/middleware"
	"icecreamshop/internal/storage"
)

func RegisterRoutes(router *gin.Engine, storage storage.Storage, middleware *middleware.Middleware) {
	handler := newHandler(storage)

	userRoutes := router.Group("/users", middleware.AuthenticateAdmin)
	{
		userRoutes.GET("", handler.GetUsers)
		userRoutes.GET("/:id", handler.GetUserByID)
		userRoutes.DELETE("/:id", handler.DeleteUserByID)
		userRoutes.PUT("/:id/admin", handler.PromoteToAdmin)
	}
}
