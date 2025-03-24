package myOrders

import (
	"github.com/gin-gonic/gin"
	"icecreamshop/internal/messageErrors"
	"icecreamshop/internal/services/payment"
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

// GetAllMyOrders handles the GET request to obtain all order from the user who is logged in.
func (h *handler) GetAllMyOrders(c *gin.Context) {
	userEmail, _ := c.Get("user-email")
	orders := h.Store.GetAllOrdersByUserEmail(userEmail.(string))
	c.JSON(http.StatusOK, orders)
}

// CreateOrder handles the POST request to create a new order for the user who is logged in.
func (h *handler) CreateOrder(c *gin.Context) {
	userID, _ := c.Get("user-id")

	var order types.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidJsonFormat})
		return
	}
	order.UserID = userID.(uint)

	if err := order.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Store.CreateOrder(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetMyOrderByID handles the GET request to obtain an order by id from the user who is logged in.
func (h *handler) GetMyOrderByID(c *gin.Context) {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user-id")
	order, err := h.Store.GetUserOrderByID(id, userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

// UpdateMyOrderByID handles the PUT request to update an order by id from the user who is logged in.
func (h *handler) UpdateMyOrderByID(c *gin.Context) {
	userID, _ := c.Get("user-id")

	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedOrder types.Order
	if err := c.ShouldBindJSON(&updatedOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidJsonFormat})
		return
	}

	if err := updatedOrder.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedOrder.ID = id
	updatedOrder.UserID = userID.(uint)
	order, err := h.Store.UpdateOrderByID(id, &updatedOrder)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetIceCreamTubsFromOrderByID handles the GET request to obtain all ice cream tubs from an order by order id. User must be the order's owner.
func (h *handler) GetIceCreamTubsFromOrderByID(c *gin.Context) {
	id, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user-id")
	if _, err := h.Store.GetUserOrderByID(id, userID.(uint)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	tubs, err := h.Store.GetIceCreamTubsByOrderID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tubs)
}

// AddIceCreamTubToOrderByID handles the POST request to add a new ice cream tub to an order by its id. User must be the order's owner.
func (h *handler) AddIceCreamTubToOrderByID(c *gin.Context) {
	orderID, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tub types.IceCreamTub
	if err := c.ShouldBindJSON(&tub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidJsonFormat})
		return
	}
	if err := tub.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user-id")
	if _, err := h.Store.GetUserOrderByID(orderID, userID.(uint)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messageErrors.OrderNotFound})
		return
	}

	err = h.Store.AddIceCreamTubByOrderID(orderID, &tub)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tub)
}

// DeleteIceCreamTubByIDFromOrder handles the DELETE request to delete an ice cream tub from an order by order id and tub id. User must be the order's owner and tub must exist in the order.
func (h *handler) DeleteIceCreamTubByIDFromOrder(c *gin.Context) {
	orderID, err := utils.StringToUint(c.Param("orderID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tubID, err := utils.StringToUint(c.Param("tubID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user-id")
	if _, err := h.Store.GetUserOrderByID(orderID, userID.(uint)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messageErrors.OrderNotFound})
		return
	}

	err = h.Store.DeleteIceCreamTubByOrderID(tubID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetDeliveryDriverFromOrder handles the GET request to obtain the assigned delivery driver from an order by its id. User must be order's owner.
func (h *handler) GetDeliveryDriverFromOrder(c *gin.Context) {
	orderID, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user-id")
	if _, err := h.Store.GetUserOrderByID(orderID, userID.(uint)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messageErrors.OrderNotFound})
		return
	}

	deliveryDriverID, err := h.Store.GetDeliveryDriverFromOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, deliveryDriverID)
}

// ProcessOrderPayment handles the POST request to process the order payment by its id. User must be order's owner.
func (h *handler) ProcessOrderPayment(c *gin.Context) {
	userID, _ := c.Get("user-id")

	orderID, err := utils.StringToUint(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.Store.GetUserOrderByID(orderID, userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messageErrors.OrderNotFound})
		return
	}

	var paymentData payment.PaymentRequest
	if err := c.ShouldBindJSON(&paymentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messageErrors.InvalidJsonFormat})
		return
	}

	paymentResponse, err := payment.ProcessPayment(paymentData, order.TotalCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.PaymentState = "paid"
	_, err = h.Store.UpdateOrderByID(orderID, &order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, paymentResponse)
}
