package myAccount

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/middleware"
	"icecreamshop/internal/storage"
)

func RegisterRoutes(router *gin.Engine, storage storage.Storage, middleware *middleware.Middleware) {
	handler := newHandler(storage)

	router.POST("/signup", handler.SignUpUser)
	router.POST("/login", middleware.CheckIfNotLoggedIn, handler.LogInUser)

	accountRoutes := router.Group("/my-account", middleware.AuthenticateUser)
	{
		accountRoutes.GET("", handler.GetMyAccount)
		accountRoutes.DELETE("", handler.DeleteMyAccount)
		accountRoutes.PUT("", handler.UpdateMyAccount)
		accountRoutes.PUT("/delivery-driver", middleware.AuthenticateDeliveryDriver, handler.UpdateDeliveryDriverData)
		accountRoutes.DELETE("/delivery-driver", middleware.AuthenticateDeliveryDriver, handler.DeleteDeliveryDriver)
	}
}
