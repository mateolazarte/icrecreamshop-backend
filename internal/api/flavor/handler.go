package flavor

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/messageErrors"
	"icecreamshop/internal/storage"
	"icecreamshop/internal/types"
	"net/http"
)

type handler struct {
	Store storage.Storage
}

func newHandler(storage storage.Storage) *handler {
	return &handler{storage}
}

// GetFlavors handles the GET request to obtain all flavor.
func (handler *handler) GetFlavors(c *gin.Context) {
	kind := c.Query("type")
	if kind != "" {
		flavors := handler.Store.GetFlavorsByType(kind)
		c.JSON(http.StatusOK, flavors)
		return
	}
	flavors := handler.Store.GetFlavors()
	c.JSON(http.StatusOK, flavors)
}

// GetFlavorByID handles the GET request to obtain a flavor by ID.
func (handler *handler) GetFlavorByID(c *gin.Context) {
	id := c.Param("id")
	flavor, err := handler.Store.GetFlavorByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messageErrors.FlavorNotFound})
		return
	}
	c.JSON(http.StatusOK, flavor)
}

// AddFlavor handles the POST request to add a flavor (only admins).
func (handler *handler) AddFlavor(c *gin.Context) {
	var flavor types.Flavor
	if err := c.ShouldBindJSON(&flavor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidJsonFormat})
		return
	}

	if err := flavor.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := handler.Store.AddFlavor(flavor)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, flavor)
}
