package api

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/api/deliveryDriver"
	"icecreamshop/internal/api/flavor"
	"icecreamshop/internal/api/myAccount"
	"icecreamshop/internal/api/myOrders"
	"icecreamshop/internal/api/order"
	"icecreamshop/internal/api/user"
	"icecreamshop/internal/middleware"
	"icecreamshop/internal/storage"
	"os"
)

type Server struct {
	Store storage.Storage
}

func NewServer(store storage.Storage) *Server {
	return &Server{store}
}

func (server *Server) Start() error {
	router := server.SetupRouter()
	return router.Run("localhost:8080")
}

func (server *Server) SetupRouter() *gin.Engine {
	middle := middleware.NewMiddleware(server.Store)

	if os.Getenv("API_ENV") == "testing" {
		gin.SetMode(gin.TestMode)
	}

	router := gin.New()

	flavor.RegisterRoutes(router, server.Store, middle)
	order.RegisterRoutes(router, server.Store, middle)
	myOrders.RegisterRoutes(router, server.Store, middle)
	deliveryDriver.RegisterRoutes(router, server.Store, middle)
	user.RegisterRoutes(router, server.Store, middle)
	myAccount.RegisterRoutes(router, server.Store, middle)

	return router
}
