package order

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/middleware"
	"icecreamshop/internal/storage"
)

func RegisterRoutes(router *gin.Engine, storage storage.Storage, middleware *middleware.Middleware) {
	orders := newHandler(storage)

	ordersGroup := router.Group("/orders", middleware.AuthenticateAdmin)
	{
		ordersGroup.GET("", orders.GetAllOrders)
		ordersGroup.GET("/:id", orders.GetOrderByID)
		ordersGroup.PUT("/:id/delivery-driver", orders.AssignDeliveryDriverToOrder)
		ordersGroup.DELETE("/:id/delivery-driver", orders.DeleteDeliveryDriverFromOrder)
	}
}
