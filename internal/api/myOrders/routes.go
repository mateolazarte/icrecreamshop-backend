package myOrders

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/middleware"
	"icecreamshop/internal/storage"
)

func RegisterRoutes(router *gin.Engine, storage storage.Storage, middleware *middleware.Middleware) {
	handler := newHandler(storage)

	myOrdersGroup := router.Group("/my-orders", middleware.AuthenticateUser)
	{
		myOrdersGroup.GET("", handler.GetAllMyOrders)
		myOrdersGroup.POST("", handler.CreateOrder)
		myOrdersGroup.GET("/:id", handler.GetMyOrderByID)
		myOrdersGroup.PUT("/:id", handler.UpdateMyOrderByID)
		myOrdersGroup.GET("/:id/tubs", handler.GetIceCreamTubsFromOrderByID)
		myOrdersGroup.POST("/:id/tubs", handler.AddIceCreamTubToOrderByID)
		myOrdersGroup.DELETE("/:orderID/tubs/:tubID", handler.DeleteIceCreamTubByIDFromOrder)
		myOrdersGroup.GET("/:id/delivery-driver", handler.GetDeliveryDriverFromOrder)
		myOrdersGroup.POST("/:id/pay", handler.ProcessOrderPayment)
	}
}
