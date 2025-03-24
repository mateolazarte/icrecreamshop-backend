package storage

import (
	"icecreamshop/internal/types"
)

// Storage interface declares the methods needed for the api to work with de database.
type Storage interface {
	// GetFlavors obtains all flavors.
	GetFlavors() []types.Flavor
	// GetFlavorsByType obtains all flavors filtered by type
	GetFlavorsByType(kind string) []types.Flavor
	// GetFlavorByID obtains a flavor by its ID
	GetFlavorByID(idFlavor string) (types.Flavor, error)
	// AddFlavor adds a new flavor
	AddFlavor(flavor types.Flavor) error

	// GetAllOrders obtains all orders from all users
	GetAllOrders() []types.Order
	// GetAllOrdersByUserEmail obtains all orders from an user by their email
	GetAllOrdersByUserEmail(email string) []types.Order
	// CreateOrder creates a new order for an user.
	// The order struct inputted must include the user id.
	CreateOrder(order *types.Order) error
	// GetOrderByID obtains an order by its id.
	GetOrderByID(idOrder uint) (types.Order, error)
	// GetUserOrderByID obtains an order from an user.
	// Method checks if the user is the order's owner. Otherwise, it will return an error.
	GetUserOrderByID(orderID uint, userID uint) (types.Order, error)
	// UpdateOrderByID updates an order by its id.
	// The order struct inputted must include the new data, but it does not need the order id
	UpdateOrderByID(idOrder uint, order *types.Order) (types.Order, error)
	// GetIceCreamTubsByOrderID obtains all ice cream tubs from an order by its id.
	GetIceCreamTubsByOrderID(idOrder uint) ([]types.IceCreamTub, error)
	// AddIceCreamTubByOrderID adds a new ice cream tub to an order by its id.
	AddIceCreamTubByOrderID(idOrder uint, iceCreamTub *types.IceCreamTub) error
	// DeleteIceCreamTubByOrderID deletes an ice cream tub from an order.
	DeleteIceCreamTubByOrderID(tubID uint, orderID uint) error

	// GetDeliveryDrivers obtains all delivery drivers.
	GetDeliveryDrivers() []types.DeliveryDriver
	// AddDeliveryDriver adds a new delivery driver.
	AddDeliveryDriver(deliveryDriver *types.DeliveryDriver) error
	// GetDeliveryDriverByID obtains a delivery driver by their user id.
	GetDeliveryDriverByID(idUser uint) (types.DeliveryDriver, error)
	// UpdateDeliveryDriverByID updates a delivery driver by their user id.
	// The delivery driver struct inputted must include the new data, but it does not need the user id.
	UpdateDeliveryDriverByID(idUser uint, deliveryDriver *types.DeliveryDriver) error
	// DeleteDeliveryDriverByID deletes a delivery driver by their user id.
	DeleteDeliveryDriverByID(idUser uint) error
	// GetVehiclesByDeliveryDriverID obtains all vehicles from a delivery driver by their id.
	GetVehiclesByDeliveryDriverID(idUser uint) ([]string, error)
	// AssignDeliveryDriverToOrder assigns a delivery driver id to an order.
	AssignDeliveryDriverToOrder(orderID uint, deliveryDriverID uint) error
	// DeleteDeliveryDriverFromOrder deletes the delivery driver id from an order.
	// No delivery driver id assigned is represented by zero.
	DeleteDeliveryDriverFromOrder(idOrder uint) error
	// GetDeliveryDriverFromOrder obtains the delivery driver id assigned to an order.
	GetDeliveryDriverFromOrder(idOrder uint) (uint, error)

	// SignUpUser signs up a new user.
	SignUpUser(newUser *types.User) error
	// LogInUser logs in an user by inputting their email and password.
	// If successful, error will be nil.
	LogInUser(email string, password string) error
	// GetUserByEmail obtains an user by its email.
	GetUserByEmail(email string) (types.User, error)
	// GetAllUsers obtains all users.
	GetAllUsers() []types.User
	// GetUserByID obtains an user by its id.
	GetUserByID(userID uint) (types.User, error)
	// DeleteUserByID delete an user by its id.
	DeleteUserByID(userID uint) error
	// UpdateUser updates an user.
	// The user struct inputted must include the user id to change.
	UpdateUser(updatedUser types.User) (types.User, error)
	// PromoteUserToAdmin promotes an user to admin by its id.
	PromoteUserToAdmin(idUser uint) error

	// Close closes db connection if needed.
	Close() error
	// CleanDB cleans db data completely, only for testing.
	CleanDB() error
}
