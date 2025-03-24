package storage

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"icecreamshop/internal/messageErrors"
	"icecreamshop/internal/types"
	"icecreamshop/internal/utils"
)

type Memory struct {
	Flavors         []types.Flavor
	Users           []types.User
	DeliveryDrivers []types.DeliveryDriver
	Orders          []types.Order
	Prices          map[uint]uint
	idOrders        uint
	idUsers         uint
	idTubs          uint
}

func NewMemoryStorage(flavors []types.Flavor, users []types.User, prices map[uint]uint) *Memory {

	// Copying slices
	flavorsCopy := append([]types.Flavor(nil), flavors...)
	usersCopy := append([]types.User(nil), users...)

	// Copying map
	pricesCopy := make(map[uint]uint)
	for k, v := range prices {
		pricesCopy[k] = v
	}

	return &Memory{
		Flavors:         flavorsCopy,
		Users:           usersCopy,
		DeliveryDrivers: []types.DeliveryDriver{},
		Orders:          []types.Order{},
		Prices:          prices,
		idOrders:        1,
		idUsers:         uint(len(users) + 1),
		idTubs:          1,
	}
}

/*******************/
/***** FLAVORS *****/
/*******************/

func (memory *Memory) GetFlavors() []types.Flavor {
	return memory.Flavors
}

func (memory *Memory) GetFlavorsByType(kind string) []types.Flavor {
	var flavors []types.Flavor
	for _, flavor := range memory.Flavors {
		if flavor.Type == kind {
			flavors = append(flavors, flavor)
		}
	}
	return flavors
}

func (memory *Memory) GetFlavorByID(idFlavor string) (types.Flavor, error) {
	for _, flavor := range memory.Flavors {
		if flavor.ID == idFlavor {
			return flavor, nil
		}
	}
	return types.Flavor{}, errors.New(messageErrors.FlavorNotFound)
}

func (memory *Memory) AddFlavor(newFlavor types.Flavor) error {
	for _, flavor := range memory.Flavors {
		if flavor.ID == newFlavor.ID {
			return errors.New(messageErrors.AlreadyExistingFlavor)
		}
	}
	memory.Flavors = append(memory.Flavors, newFlavor)
	return nil
}

/******************/
/***** ORDERS *****/
/******************/

func (memory *Memory) CreateOrder(order *types.Order) error {
	order.ID = memory.idOrders
	order.PaymentState = "pending"

	for i := range memory.Users {
		if memory.Users[i].ID == order.UserID {
			memory.Users[i].Orders = append(memory.Users[i].Orders, *order)
			memory.Orders = append(memory.Orders, *order)
			memory.idOrders++
			return nil
		}
	}

	return errors.New(messageErrors.UserIDNotFound)
}

func (memory *Memory) GetOrderByID(idOrder uint) (types.Order, error) {
	for _, order := range memory.Orders {
		if order.ID == idOrder {
			return order, nil
		}
	}
	return types.Order{}, errors.New(messageErrors.OrderNotFound)
}

func (memory *Memory) GetAllOrders() []types.Order {
	return memory.Orders
}

func (memory *Memory) GetAllOrdersByUserEmail(email string) []types.Order {
	for _, user := range memory.Users {
		if user.Email == email {
			return user.Orders
		}
	}
	return []types.Order{}
}

func (memory *Memory) GetUserOrderByID(orderID uint, userID uint) (types.Order, error) {
	for _, user := range memory.Users {
		if user.ID == userID {
			for _, order := range user.Orders {
				if order.ID == orderID {
					return memory.GetOrderByID(orderID)
				}
			}
			return types.Order{}, errors.New(messageErrors.OrderNotFound)
		}
	}
	return types.Order{}, errors.New(messageErrors.UserIDNotFound)
}

func (memory *Memory) UpdateOrderByID(orderID uint, updatedOrder *types.Order) (types.Order, error) {
	for i, order := range memory.Orders {
		if order.ID == orderID {
			if order.UserID == updatedOrder.UserID {
				memory.Orders[i].Address = updatedOrder.Address
				memory.Orders[i].PaymentState = updatedOrder.PaymentState
				return memory.Orders[i], nil
			}
			return types.Order{}, errors.New(messageErrors.OrderNotFound)
		}
	}

	return types.Order{}, errors.New(messageErrors.OrderNotFound)
}

func (memory *Memory) GetIceCreamTubsByOrderID(idOrder uint) ([]types.IceCreamTub, error) {
	for _, order := range memory.Orders {
		if order.ID == idOrder {
			return order.IceCreamTubs, nil
		}
	}
	return []types.IceCreamTub{}, errors.New(messageErrors.OrderNotFound)
}

func (memory *Memory) AddIceCreamTubByOrderID(idOrder uint, iceCreamTub *types.IceCreamTub) error {

	if ok := areFlavorIDsRegisteredInMemory(iceCreamTub.Flavors, memory.Flavors); !ok {
		return errors.New(messageErrors.NonExistingFlavors)
	}

	price, ok := memory.Prices[iceCreamTub.Weight]
	if !ok {
		return errors.New(messageErrors.WeightNotAvailable)
	}

	for i := 0; i < len(memory.Orders); i++ {
		if memory.Orders[i].ID == idOrder {
			iceCreamTub.ID = memory.idTubs
			memory.Orders[i].IceCreamTubs = append(memory.Orders[i].IceCreamTubs, *iceCreamTub)
			memory.idTubs++
			memory.Orders[i].TotalCost += price
			return nil
		}
	}

	return errors.New(messageErrors.OrderNotFound)
}

func (memory *Memory) DeleteIceCreamTubByOrderID(tubID uint, orderID uint) error {
	for j := 0; j < len(memory.Orders); j++ {
		if memory.Orders[j].ID == orderID {
			for i := 0; i < len(memory.Orders[j].IceCreamTubs); i++ {
				if memory.Orders[j].IceCreamTubs[i].ID == tubID {
					memory.Orders[j].TotalCost -= memory.Prices[memory.Orders[j].IceCreamTubs[i].Weight]
					memory.Orders[j].IceCreamTubs = append(memory.Orders[j].IceCreamTubs[:i], memory.Orders[j].IceCreamTubs[i+1:]...)
					return nil
				}
			}
			return errors.New(messageErrors.IceCreamTubNotFound)
		}
	}
	return errors.New(messageErrors.OrderNotFound)
}

/****************************/
/***** DELIVERY DRIVERS *****/
/****************************/

func (memory *Memory) GetDeliveryDrivers() []types.DeliveryDriver {
	return memory.DeliveryDrivers
}

func (memory *Memory) GetDeliveryDriverByID(idUser uint) (types.DeliveryDriver, error) {
	for _, deliveryDriver := range memory.DeliveryDrivers {
		if deliveryDriver.UserID == idUser {
			return deliveryDriver, nil
		}
	}
	return types.DeliveryDriver{}, errors.New(messageErrors.DeliveryDriverNotFound)
}

func (memory *Memory) UpdateDeliveryDriverByID(idUser uint, deliveryDriver *types.DeliveryDriver) error {
	for i := 0; i < len(memory.DeliveryDrivers); i++ {
		if memory.DeliveryDrivers[i].UserID == idUser {
			deliveryDriver.UserID = memory.DeliveryDrivers[i].UserID
			memory.DeliveryDrivers[i] = *deliveryDriver
			return nil
		}
	}
	return errors.New(messageErrors.DeliveryDriverNotFound)
}

func (memory *Memory) DeleteDeliveryDriverByID(idUser uint) error {
	for i, deliveryDriver := range memory.DeliveryDrivers {
		if deliveryDriver.UserID == idUser {
			memory.DeliveryDrivers = append(memory.DeliveryDrivers[:i], memory.DeliveryDrivers[i+1:]...)
			for i, user := range memory.Users {
				if user.ID == deliveryDriver.UserID {
					memory.Users[i].Permissions = utils.DeletePermission(memory.Users[i].Permissions, "delivery")
				}
			}
			return nil
		}
	}
	return errors.New(messageErrors.DeliveryDriverNotFound)
}

func (memory *Memory) GetVehiclesByDeliveryDriverID(idUser uint) ([]string, error) {
	for _, deliveryDriver := range memory.DeliveryDrivers {
		if deliveryDriver.UserID == idUser {
			return deliveryDriver.Vehicles, nil
		}
	}
	return []string{}, errors.New(messageErrors.DeliveryDriverNotFound)
}

func (memory *Memory) AddDeliveryDriver(deliveryDriver *types.DeliveryDriver) error {
	for i, user := range memory.Users {
		if user.ID == deliveryDriver.UserID {
			if user.IsDeliveryDriver() {
				return errors.New(messageErrors.UserIsAlreadyADriver)
			}
			memory.DeliveryDrivers = append(memory.DeliveryDrivers, *deliveryDriver)
			memory.Users[i].Permissions = append(user.Permissions, "delivery")
			return nil
		}
	}
	return errors.New(messageErrors.UserIDNotFound)
}

func (memory *Memory) AssignDeliveryDriverToOrder(orderID uint, deliveryDriverID uint) error {
	if !isDelvieryDriverIDRegisteredInMemory(deliveryDriverID, memory.DeliveryDrivers) {
		return errors.New(messageErrors.DeliveryDriverNotFound)
	}

	for i := 0; i < len(memory.Orders); i++ {
		if memory.Orders[i].ID == orderID {
			memory.Orders[i].DeliveryDriverID = deliveryDriverID
			return nil
		}
	}

	return errors.New(messageErrors.OrderNotFound)
}

func (memory *Memory) DeleteDeliveryDriverFromOrder(idOrder uint) error {
	for i := 0; i < len(memory.Orders); i++ {
		if memory.Orders[i].ID == idOrder {
			memory.Orders[i].DeliveryDriverID = 0
			return nil
		}
	}
	return errors.New(messageErrors.OrderNotFound)
}

func (memory *Memory) GetDeliveryDriverFromOrder(idOrder uint) (uint, error) {
	for _, order := range memory.Orders {
		if order.ID == idOrder {
			return order.DeliveryDriverID, nil
		}
	}
	return 0, errors.New(messageErrors.OrderNotFound)
}

/*****************/
/***** USERS *****/
/*****************/

func (memory *Memory) GetAllUsers() []types.User {
	return memory.Users
}

func (memory *Memory) SignUpUser(newUser *types.User) error {
	for _, user := range memory.Users {
		if user.Email == newUser.Email {
			return errors.New(messageErrors.EmailAlreadyExists)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
	if err != nil {
		return errors.New(messageErrors.ErrorWhileProcessingRequest)
	}
	newUser.Password = string(hashedPassword)

	newUser.ID = memory.idUsers
	memory.idUsers++
	memory.Users = append(memory.Users, *newUser)

	return nil
}

func (memory *Memory) LogInUser(email string, password string) error {
	for _, user := range memory.Users {
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return errors.New(messageErrors.InvalidEmailOrPassword)
			}
			return nil
		}
	}
	return errors.New(messageErrors.InvalidEmailOrPassword)
}

func (memory *Memory) GetUserByEmail(email string) (types.User, error) {
	for _, user := range memory.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return types.User{}, errors.New(messageErrors.UserEmailNotFound)
}

func (memory *Memory) GetUserByID(userID uint) (types.User, error) {
	for _, user := range memory.Users {
		if user.ID == userID {
			return user, nil
		}
	}
	return types.User{}, errors.New(messageErrors.UserIDNotFound)
}

func (memory *Memory) DeleteUserByID(userID uint) error {
	for i := 0; i < len(memory.Users); i++ {
		if memory.Users[i].ID == userID {
			if memory.Users[i].IsDeliveryDriver() {
				memory.DeleteDeliveryDriverByID(userID)
			}
			memory.Users = append(memory.Users[:i], memory.Users[i+1:]...)
			return nil
		}
	}
	return errors.New(messageErrors.UserIDNotFound)
}

func (memory *Memory) UpdateUser(updatedUser types.User) (types.User, error) {
	for i := 0; i < len(memory.Users); i++ {
		if memory.Users[i].ID == updatedUser.ID {
			memory.Users[i].Email = updatedUser.Email
			memory.Users[i].Name = updatedUser.Name
			memory.Users[i].LastName = updatedUser.LastName
			return memory.Users[i], nil
		}
	}
	return types.User{}, errors.New(messageErrors.UserIDNotFound)
}

func (memory *Memory) PromoteUserToAdmin(idUser uint) error {
	for i := 0; i < len(memory.Users); i++ {
		if memory.Users[i].ID == idUser {
			if memory.Users[i].IsAdmin() {
				return errors.New(messageErrors.UserIsAlreadyAnAdmin)
			}
			memory.Users[i].Permissions = append(memory.Users[i].Permissions, "admin")
			return nil
		}
	}
	return errors.New(messageErrors.UserIDNotFound)
}

// Others

func (memory *Memory) CleanDB() error {
	//Do nothing
	return nil
}

func (memory *Memory) Close() error {
	//Do nothing
	return nil
}

// Auxiliary functions

func isDelvieryDriverIDRegisteredInMemory(idUser uint, deliveryDrivers []types.DeliveryDriver) bool {
	for _, deliveryDriver := range deliveryDrivers {
		if deliveryDriver.UserID == idUser {
			return true
		}
	}
	return false
}

func isFlavorIDRegisteredInMemory(flavorID string, flavors []types.Flavor) bool {
	for _, flavor := range flavors {
		if flavor.ID == flavorID {
			return true
		}
	}
	return false
}

func areFlavorIDsRegisteredInMemory(flavorIDs []string, flavors []types.Flavor) bool {
	for _, flavorID := range flavorIDs {
		if !isFlavorIDRegisteredInMemory(flavorID, flavors) {
			return false
		}
	}
	return true
}
