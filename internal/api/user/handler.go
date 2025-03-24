package user

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

// GetUsers handles the GET request to obtain all users (only admins)
func (h *handler) GetUsers(c *gin.Context) {
	users := h.Store.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

// GetUserByID handles the GET request to obtain any user by id (only admins)
func (h *handler) GetUserByID(c *gin.Context) {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.Store.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUserByID handles the DELETE request to delete any user by id (only admins)
func (h *handler) DeleteUserByID(c *gin.Context) {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.Store.DeleteUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// PromoteToAdmin handles the PUT request to promote any user to admin by id (only admins)
func (h *handler) PromoteToAdmin(c *gin.Context) {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.Store.PromoteUserToAdmin(id)
	if err != nil {
		if err.Error() == messageErrors.UserIDNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"description": "The user has been promoted to admin"})
}
