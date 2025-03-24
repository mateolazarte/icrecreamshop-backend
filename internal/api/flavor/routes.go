package flavor

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/middleware"
	"icecreamshop/internal/storage"
)

func RegisterRoutes(router *gin.Engine, storage storage.Storage, middleware *middleware.Middleware) {
	handler := newHandler(storage)

	flavorsGroup := router.Group("/flavors")
	{
		flavorsGroup.GET("", handler.GetFlavors)
		flavorsGroup.GET("/:id", handler.GetFlavorByID)
		flavorsGroup.POST("", middleware.AuthenticateAdmin, handler.AddFlavor)
	}
}
