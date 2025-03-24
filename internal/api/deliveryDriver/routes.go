package deliveryDriver

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/middleware"
	"icecreamshop/internal/storage"
)

func RegisterRoutes(router *gin.Engine, storage storage.Storage, middleware *middleware.Middleware) {
	handler := newHandler(storage)

	deliveryDriversGroup := router.Group("/delivery-drivers", middleware.AuthenticateUser)
	{
		deliveryDriversGroup.GET("", middleware.AuthenticateAdmin, handler.GetAllDeliveryDrivers)
		deliveryDriversGroup.POST("", middleware.AuthenticateAdmin, handler.AddDeliveryDriver)
		deliveryDriversGroup.GET("/:id", middleware.AuthenticateAdmin, handler.GetDeliveryDriverByID)
	}
}
