package myAccount

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/auth"
	"icecreamshop/internal/messageErrors"
	"icecreamshop/internal/storage"
	"icecreamshop/internal/types"
	"net/http"
)

type handler struct {
	Store storage.Storage
}

func newHandler(store storage.Storage) *handler {
	return &handler{store}
}

// SignUpUser handles the POST request to sign up a new user.
func (h *handler) SignUpUser(c *gin.Context) {
	var signUpUser types.SignUpInput
	if err := c.ShouldBindJSON(&signUpUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidJsonFormat})
		return
	}

	user := types.User{
		Email:    signUpUser.Email,
		Password: signUpUser.Password,
		Name:     signUpUser.Name,
		LastName: signUpUser.LastName,
	}
	if err := user.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Store.SignUpUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// LogInUser handles the POST request to log in an user.
// User must be already registered in the system.
func (h *handler) LogInUser(c *gin.Context) {
	var body struct {
		Email    string
		Password string
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidJsonFormat})
		return
	}

	err := h.Store.LogInUser(body.Email, body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidEmailOrPassword})
		return
	}

	tokenString := auth.GenerateTokenFromUserEmail(body.Email)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 60*60*24, "", "", false, true)

	c.JSON(http.StatusOK, nil)
}

// GetMyAccount handles the GET request to obtain all data from the user who is logged in.
// It also checks if the user is a delivery driver to return extra information if needed.
func (h *handler) GetMyAccount(c *gin.Context) {
	userEmail, _ := c.Get("user-email")
	user, err := h.Store.GetUserByEmail(userEmail.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	for _, permission := range user.Permissions {
		if permission == "repartidor" {
			deliveryDriver, _ := h.Store.GetDeliveryDriverByID(user.ID)
			response := struct {
				types.User
				DeliveryDriver types.DeliveryDriver `json:"deliveryDriver,omitempty"` // will be included only if not nil
			}{
				User:           user,
				DeliveryDriver: deliveryDriver,
			}
			c.JSON(http.StatusOK, response)
			return
		}
	}
	c.JSON(http.StatusOK, user)
}

// DeleteMyAccount handles the DELETE request to delete the account of the user who is logged in.
// Automatically logs out the user.
func (h *handler) DeleteMyAccount(c *gin.Context) {
	userID, _ := c.Get("user-id")
	h.Store.DeleteUserByID(userID.(uint))
	c.SetCookie("Authorization", "", 0, "", "", false, true)
	c.JSON(http.StatusNoContent, nil)
}

// UpdateMyAccount handles the POST request to update the account data of the user who is logged in.
func (h *handler) UpdateMyAccount(c *gin.Context) {
	userID, _ := c.Get("user-id")
	var userUpdated types.User
	if err := c.ShouldBindJSON(&userUpdated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidJsonFormat})
		return
	}

	userUpdated.ID = userID.(uint)
	if err := userUpdated.ValidateUserDataWithoutPassword(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Store.UpdateUser(userUpdated)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateDeliveryDriverData handles the POST request to update the delivery driver data of the user who is logged in.
// User must be a delivery driver.
func (h *handler) UpdateDeliveryDriverData(c *gin.Context) {
	userID, _ := c.Get("user-id")

	var updatedDeliveryDriver types.DeliveryDriver
	if err := c.ShouldBindJSON(&updatedDeliveryDriver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidJsonFormat})
		return
	}

	if err := updatedDeliveryDriver.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Store.UpdateDeliveryDriverByID(userID.(uint), &updatedDeliveryDriver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedDeliveryDriver)
}

// DeleteDeliveryDriver handles the DELETE request to delete the delivery driver who is logged in.
func (h *handler) DeleteDeliveryDriver(c *gin.Context) {
	userID, _ := c.Get("user-id")

	err := h.Store.DeleteDeliveryDriverByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
