package storage

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"icecreamshop/internal/messageErrors"
	"icecreamshop/internal/types"
	"icecreamshop/internal/utils"
	"os"
)

type DbStorage struct {
	DB *gorm.DB
}

func setDBPath() string {
	if os.Getenv("API_ENV") == "testing" {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("TEST_DB_HOST"),
			os.Getenv("TEST_DB_PORT"),
			os.Getenv("TEST_DB_USER"),
			os.Getenv("TEST_DB_PASSWORD"),
			os.Getenv("TEST_DB_NAME"),
		)
	} else {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
	}
}

func NewDBStorage(flavors []types.Flavor, users []types.User, prices map[uint]uint) *DbStorage {
	dbPath := setDBPath()

	db, err := gorm.Open(postgres.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	err = db.AutoMigrate(&types.User{}, &types.DeliveryDriver{}, &types.Order{}, &types.Flavor{}, &types.IceCreamTub{}, &types.IceCreamTubPrice{})
	if err != nil {
		panic("failed to automigrate data")
	}

	var iceCreamTubPrices []types.IceCreamTubPrice
	for key, value := range prices {
		iceCreamTubPrices = append(iceCreamTubPrices, types.IceCreamTubPrice{key, value})
	}

	db.Create(&flavors)
	db.Create(&users)
	db.Create(&iceCreamTubPrices)

	db.Exec("SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));")

	return &DbStorage{DB: db}
}

/*******************/
/***** FLAVORS *****/
/*******************/

func (dbStorage *DbStorage) GetFlavors() []types.Flavor {
	var flavors []types.Flavor
	dbStorage.DB.Find(&flavors)
	return flavors
}

func (dbStorage *DbStorage) GetFlavorsByType(kind string) []types.Flavor {
	var flavors []types.Flavor
	dbStorage.DB.Where("type = ?", kind).Find(&flavors)
	return flavors
}

func (dbStorage *DbStorage) GetFlavorByID(flavorID string) (types.Flavor, error) {
	var flavor types.Flavor
	err := dbStorage.DB.First(&flavor, "id = ?", flavorID).Error
	if err != nil {
		return types.Flavor{}, errors.New(messageErrors.FlavorNotFound)
	}
	return flavor, nil
}

func (dbStorage *DbStorage) AddFlavor(flavor types.Flavor) error {
	err := dbStorage.DB.Create(&flavor).Error
	if err != nil {
		return errors.New(messageErrors.AlreadyExistingFlavor)
	}
	return nil
}

/******************/
/***** ORDERS *****/
/******************/

func (dbStorage *DbStorage) CreateOrder(order *types.Order) error {
	err := dbStorage.DB.First(&types.User{}, order.UserID).Error
	if err != nil {
		return errors.New(messageErrors.UserIDNotFound)
	}
	order.PaymentState = "pending"
	err = dbStorage.DB.Create(&order).Error
	if err != nil {
		print(err.Error())
	}
	return nil
}

func (dbStorage *DbStorage) GetOrderByID(idOrder uint) (types.Order, error) {
	var order types.Order
	err := dbStorage.DB.Preload("IceCreamTubs").First(&order, idOrder).Error
	if err != nil {
		return types.Order{}, errors.New(messageErrors.OrderNotFound)
	}
	return order, nil
}

func (dbStorage *DbStorage) GetAllOrders() []types.Order {
	var orders []types.Order
	dbStorage.DB.Find(&orders)
	return orders
}

func (dbStorage *DbStorage) GetAllOrdersByUserEmail(email string) []types.Order {
	var user types.User
	dbStorage.DB.Preload("Orders").Where("email = ?", email).First(&user)
	return user.Orders
}

func (dbStorage *DbStorage) GetUserOrderByID(idOrder uint, idUser uint) (types.Order, error) {
	err := dbStorage.DB.First(&types.User{}, idUser).Error
	if err != nil {
		return types.Order{}, errors.New(messageErrors.UserIDNotFound)
	}
	order, err := dbStorage.GetOrderByID(idOrder)
	if err != nil {
		return types.Order{}, errors.New(messageErrors.OrderNotFound)
	}
	if order.UserID != idUser {
		return types.Order{}, errors.New(messageErrors.OrderNotFound)
	}
	return order, nil
}

func (dbStorage *DbStorage) UpdateOrderByID(idOrder uint, order *types.Order) (types.Order, error) {
	oldPedido, err := dbStorage.GetOrderByID(idOrder)
	if err != nil {
		return types.Order{}, errors.New(messageErrors.OrderNotFound)
	}
	if oldPedido.UserID != order.UserID {
		return types.Order{}, errors.New(messageErrors.OrderNotFound)
	}
	oldPedido.PaymentState = order.PaymentState
	oldPedido.Address = order.Address
	err = dbStorage.DB.Save(&oldPedido).Error
	if err != nil {
		return types.Order{}, errors.New(messageErrors.OrderNotFound)
	}
	oldPedido.Address = order.Address
	return oldPedido, nil
}

func (dbStorage *DbStorage) GetIceCreamTubsByOrderID(idOrder uint) ([]types.IceCreamTub, error) {
	err := dbStorage.DB.First(&types.Order{}, idOrder).Error
	if err != nil {
		return []types.IceCreamTub{}, errors.New(messageErrors.OrderNotFound)
	}
	var tubs []types.IceCreamTub
	err = dbStorage.DB.Find(&tubs, "order_id = ?", idOrder).Error
	if err != nil {
		return []types.IceCreamTub{}, errors.New(messageErrors.OrderNotFound)
	}
	return tubs, nil
}

func (dbStorage *DbStorage) AddIceCreamTubByOrderID(orderID uint, tub *types.IceCreamTub) error {
	if ok := areFlavorIDsRegisteredInDB(tub.Flavors, dbStorage.DB); !ok {
		return errors.New(messageErrors.NonExistingFlavors)
	}

	var price uint
	err := dbStorage.DB.Model(&types.IceCreamTubPrice{}).Select("price").Where("weight=?", tub.Weight).First(&price).Error
	if err != nil {
		return errors.New(messageErrors.WeightNotAvailable)
	}

	order, err := dbStorage.GetOrderByID(orderID)
	if err != nil {
		return errors.New(messageErrors.OrderNotFound)
	}

	tub.OrderID = order.ID
	order.TotalCost += price
	err = dbStorage.DB.Create(&tub).Error
	if err != nil {
		return errors.New("Couldn't create IceCreamTub")
	}

	err = dbStorage.DB.Model(&types.Order{}).Where("id=?", order.ID).Update("total_cost", order.TotalCost).Error
	if err != nil {
		return errors.New("Couldn't update IceCreamTub")
	}
	return nil
}

func (dbStorage *DbStorage) DeleteIceCreamTubByOrderID(idTub uint, idOrder uint) error {
	var order types.Order
	err := dbStorage.DB.First(&order, idOrder).Error
	if err != nil {
		return errors.New(messageErrors.OrderNotFound)
	}
	var tub types.IceCreamTub
	err = dbStorage.DB.First(&tub, idTub).Error
	if err != nil {
		return errors.New(messageErrors.IceCreamTubNotFound)
	}
	if tub.OrderID != idOrder {
		return errors.New(messageErrors.IceCreamTubNotFound)
	}
	var price uint
	dbStorage.DB.Model(&types.IceCreamTubPrice{}).Select("price").Where("weight=?", tub.Weight).First(&price)
	result := dbStorage.DB.Delete(&types.IceCreamTub{}, idTub)
	if result.RowsAffected == 0 {
		return errors.New(messageErrors.IceCreamTubNotFound)
	}
	order.TotalCost -= price
	err = dbStorage.DB.Model(&types.Order{}).Where("id=?", order.ID).Update("total_cost", order.TotalCost).Error
	if err != nil {
		return errors.New("Couldn't update IceCreamTub")
	}
	return nil
}

/****************************/
/***** DELIVERY DRIVERS *****/
/****************************/

func (dbStorage *DbStorage) GetDeliveryDrivers() []types.DeliveryDriver {
	var deliveryDrivers []types.DeliveryDriver
	dbStorage.DB.Find(&deliveryDrivers)
	return deliveryDrivers
}

func (dbStorage *DbStorage) GetDeliveryDriverByID(idUser uint) (types.DeliveryDriver, error) {
	var deliveryDriver types.DeliveryDriver
	err := dbStorage.DB.First(&deliveryDriver, "user_id = ?", idUser).Error
	if err != nil {
		return types.DeliveryDriver{}, errors.New(messageErrors.DeliveryDriverNotFound)
	}
	return deliveryDriver, nil
}

func (dbStorage *DbStorage) UpdateDeliveryDriverByID(idUser uint, deliveryDriver *types.DeliveryDriver) error {
	deliveryDriver.UserID = idUser
	res := dbStorage.DB.Model(&deliveryDriver).Where("user_id=?", idUser).Updates(deliveryDriver)
	if res.Error != nil || res.RowsAffected == 0 {
		return errors.New(messageErrors.DeliveryDriverNotFound)
	}
	return nil
}

func (dbStorage *DbStorage) DeleteDeliveryDriverByID(idUser uint) error {
	res := dbStorage.DB.Delete(&types.DeliveryDriver{}, "user_id=?", idUser)
	if res.RowsAffected == 0 {
		return errors.New(messageErrors.DeliveryDriverNotFound)
	}
	var user types.User
	err := dbStorage.DB.First(&user, idUser).Error
	if err != nil {
		return errors.New(messageErrors.UserIDNotFound)
	}

	user.Permissions = utils.DeletePermission(user.Permissions, "delivery")
	err = dbStorage.DB.Save(&user).Error
	if err != nil {
		return errors.New("Couldn't update user")
	}

	return nil
}

func (dbStorage *DbStorage) GetVehiclesByDeliveryDriverID(idUser uint) ([]string, error) {
	var deliveryDriver types.DeliveryDriver
	err := dbStorage.DB.First(&deliveryDriver, "user_id=?", idUser).Error
	if err != nil {
		return []string{}, errors.New(messageErrors.DeliveryDriverNotFound)
	}
	return deliveryDriver.Vehicles, nil
}

func (dbStorage *DbStorage) AddDeliveryDriver(deliveryDriver *types.DeliveryDriver) error {
	var user types.User
	err := dbStorage.DB.First(&user, deliveryDriver.UserID).Error
	if err != nil {
		return errors.New(messageErrors.UserIDNotFound)
	}

	if user.IsDeliveryDriver() {
		return errors.New(messageErrors.UserIsAlreadyADriver)
	}
	dbStorage.DB.Create(&deliveryDriver)

	user.Permissions = append(user.Permissions, "delivery")
	err = dbStorage.DB.Save(&user).Error
	if err != nil {
		return errors.New("Couldn't update user")
	}

	return nil
}

func (dbStorage *DbStorage) AssignDeliveryDriverToOrder(idOrder uint, idDeliveryDriver uint) error {
	oldOrder, err := dbStorage.GetOrderByID(idOrder)
	if err != nil {
		return errors.New(messageErrors.OrderNotFound)
	}

	err = dbStorage.DB.First(&types.DeliveryDriver{}, "user_id=?", idDeliveryDriver).Error
	if err != nil {
		return errors.New(messageErrors.DeliveryDriverNotFound)
	}

	err = dbStorage.DB.Model(&oldOrder).Update("DeliveryDriverID", idDeliveryDriver).Error
	if err != nil {
		return errors.New(messageErrors.OrderNotFound)
	}

	return nil
}

func (dbStorage *DbStorage) DeleteDeliveryDriverFromOrder(idOrder uint) error {
	oldOrder, err := dbStorage.GetOrderByID(idOrder)
	if err != nil {
		return errors.New(messageErrors.OrderNotFound)
	}
	err = dbStorage.DB.Model(&oldOrder).Update("DeliveryDriverID", uint(0)).Error
	if err != nil {
		return errors.New(messageErrors.OrderNotFound)
	}
	return nil
}

func (dbStorage *DbStorage) GetDeliveryDriverFromOrder(idOrder uint) (uint, error) {
	order, err := dbStorage.GetOrderByID(idOrder)
	if err != nil {
		return 0, errors.New(messageErrors.OrderNotFound)
	}
	return order.DeliveryDriverID, nil
}

/*****************/
/***** USERS *****/
/*****************/

func (dbStorage *DbStorage) SignUpUser(newUser *types.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
	if err != nil {
		return errors.New(messageErrors.ErrorWhileProcessingRequest)
	}

	newUser.Password = string(hashedPassword)
	err = dbStorage.DB.Create(&newUser).Error
	if err != nil {
		return errors.New(messageErrors.EmailAlreadyExists)
	}
	return nil
}

func (dbStorage *DbStorage) LogInUser(email string, password string) error {
	user, err := dbStorage.GetUserByEmail(email)
	if err != nil {
		return errors.New(messageErrors.InvalidEmailOrPassword)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New(messageErrors.InvalidEmailOrPassword)
	}
	return nil
}

func (dbStorage *DbStorage) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	err := dbStorage.DB.Model(&user).Where("Email=?", email).First(&user).Error
	if err != nil {
		return user, errors.New(messageErrors.UserEmailNotFound)
	}
	return user, nil
}

func (dbStorage *DbStorage) GetAllUsers() []types.User {
	var users []types.User
	dbStorage.DB.Find(&users)
	return users
}

func (dbStorage *DbStorage) GetUserByID(idUser uint) (types.User, error) {
	var user types.User
	err := dbStorage.DB.Preload("Orders").First(&user, idUser).Error
	if err != nil {
		return user, errors.New(messageErrors.UserIDNotFound)
	}
	return user, nil
}

func (dbStorage *DbStorage) DeleteUserByID(idUser uint) error {
	result := dbStorage.DB.Delete(&types.User{}, "ID=?", idUser)
	if result.RowsAffected == 0 {
		return errors.New(messageErrors.UserIDNotFound)
	}
	return nil
}

func (dbStorage *DbStorage) UpdateUser(updatedUser types.User) (types.User, error) {
	oldUser, err := dbStorage.GetUserByID(updatedUser.ID)
	if err != nil {
		return updatedUser, errors.New(messageErrors.UserIDNotFound)
	}
	oldUser.Email = updatedUser.Email
	oldUser.Name = updatedUser.Name
	oldUser.LastName = updatedUser.LastName
	err = dbStorage.DB.Save(&oldUser).Error
	if err != nil {
		return updatedUser, errors.New(messageErrors.UserIDNotFound)
	}
	return oldUser, nil
}

func (dbStorage *DbStorage) PromoteUserToAdmin(idUser uint) error {
	user, err := dbStorage.GetUserByID(idUser)
	if err != nil {
		return errors.New(messageErrors.UserIDNotFound)
	}
	if user.IsAdmin() {
		return errors.New(messageErrors.UserIsAlreadyAnAdmin)
	}
	user.Permissions = append(user.Permissions, "admin")
	err = dbStorage.DB.Save(&user).Error
	if err != nil {
		return errors.New("User Not Found")
	}
	return nil
}

// Others

func (dbStorage *DbStorage) CleanDB() error {
	if os.Getenv("API_ENV") == "testing" {
		return dbStorage.DB.Exec(
			"TRUNCATE TABLE users, delivery_drivers, orders, flavors, ice_cream_tubs, ice_cream_tub_prices RESTART IDENTITY CASCADE",
		).Error
	}
	return errors.New("API_ENV must be set to testing in order to completely clean DB")
}

func (dbStorage *DbStorage) Close() error {
	sqlDB, err := dbStorage.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Auxiliary functions

func isFlavorIDRegisteredInDB(flavorID string, db *gorm.DB) bool {
	var flavor types.Flavor
	err := db.First(&flavor, "ID=?", flavorID).Error
	if err != nil {
		return false
	}
	return true
}

func areFlavorIDsRegisteredInDB(flavorIDs []string, db *gorm.DB) bool {
	for _, flavorID := range flavorIDs {
		if !isFlavorIDRegisteredInDB(flavorID, db) {
			return false
		}
	}
	return true
}
