package deliveryDriver

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/messageErrors"
	"icecreamshop/internal/storage"
	"icecreamshop/internal/types"
	"icecreamshop/internal/utils"
	"net/http"
)

type handler struct {
	Store storage.Storage
}

func newHandler(store storage.Storage) *handler {
	return &handler{store}
}

// GetAllDeliveryDrivers handles the GET request to obtain all delivery drivers.
func (h *handler) GetAllDeliveryDrivers(c *gin.Context) {
	deliveryDrivers := h.Store.GetDeliveryDrivers()
	c.JSON(http.StatusOK, deliveryDrivers)
}

// AddDeliveryDriver handles the POST request to add a new delivery driver.
// User id inputted in JSON Body must be an existing user.
func (h *handler) AddDeliveryDriver(c *gin.Context) {
	var deliveryDriver types.DeliveryDriver
	if err := c.ShouldBindJSON(&deliveryDriver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidJsonFormat})
		return
	}

	if err := deliveryDriver.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Store.AddDeliveryDriver(&deliveryDriver)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, deliveryDriver)
}

// GetDeliveryDriverByID handles the GET request to obtain a delivery driver by its user id.
func (h *handler) GetDeliveryDriverByID(c *gin.Context) {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deliveryDriver, err := h.Store.GetDeliveryDriverByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messageErrors.DeliveryDriverNotFound})
		return
	}

	c.JSON(http.StatusOK, deliveryDriver)
}
