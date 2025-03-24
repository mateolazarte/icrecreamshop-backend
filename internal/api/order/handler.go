package order

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/messageErrors"
	"icecreamshop/internal/storage"
	"icecreamshop/internal/utils"
	"net/http"
)

type handler struct {
	Store storage.Storage
}

func newHandler(store storage.Storage) *handler {
	return &handler{store}
}

// GetAllOrders handles the GET request to obtain all order from all users (only admins).
func (h *handler) GetAllOrders(c *gin.Context) {
	orders := h.Store.GetAllOrders()
	c.JSON(http.StatusOK, orders)
}

// GetOrderByID handles the GET request to obtain any order by ID (only admins).
func (h *handler) GetOrderByID(c *gin.Context) {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.Store.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

// AssignDeliveryDriverToOrder handles the PUT request to assign a delivery driver to an order by its ID (only admins).
func (h *handler) AssignDeliveryDriverToOrder(c *gin.Context) {
	orderID, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deliveryDriverID := struct {
		ID uint `json:"id"`
	}{}
	if err := c.ShouldBindJSON(&deliveryDriverID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidJsonFormat})
		return
	}

	err = h.Store.AssignDeliveryDriverToOrder(orderID, deliveryDriverID.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DeleteDeliveryDriverFromOrder handles the DELETE request to delete a delivery driver from an order by its ID (only admins).
func (h *handler) DeleteDeliveryDriverFromOrder(c *gin.Context) {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.Store.DeleteDeliveryDriverFromOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
