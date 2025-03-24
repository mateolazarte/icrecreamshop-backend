package tests

import (
	"encoding/json"
	"fmt"
	"icecreamshop/internal/services/payment"
	"icecreamshop/internal/types"
	"net/http"
	"net/http/httptest"
	"strings"
)

var flavorDDL types.Flavor = types.Flavor{ID: "ddl", Name: "Dulce de leche", Type: "Dulce de leches"}
var flavorMRC types.Flavor = types.Flavor{ID: "mrc", Name: "Chocolate marroc", Type: "Chocolates"}
var flavorTRM types.Flavor = types.Flavor{ID: "trm", Name: "Tramontana", Type: "Cremas"}
var flavorFRT types.Flavor = types.Flavor{ID: "frt", Name: "Frutilla al agua", Type: "Al agua"}

var flavors = []types.Flavor{
	flavorDDL, flavorMRC, flavorTRM, flavorFRT,
}

var adminUser = types.User{
	ID:          1,
	Email:       "abcde@gmail.com",
	Name:        "abcde",
	LastName:    "xyz",
	Password:    "$2a$10$xQy8YTOUh6GST9zO1cfmZeV4iPi1I5TLEr5WnTE7Y/XNHgLbqEeFO", //hash for "admin123"
	Orders:      []types.Order{},
	Permissions: []string{"admin"},
}

var genericUser = types.User{
	ID:          2,
	Email:       "zzzzz@gmail.com",
	Name:        "hello",
	LastName:    "world",
	Password:    "$2a$10$xQy8YTOUh6GST9zO1cfmZeV4iPi1I5TLEr5WnTE7Y/XNHgLbqEeFO", //hash for "admin123"
	Orders:      []types.Order{},
	Permissions: []string{},
}

var newDeliveryDriverForGenericUser = types.DeliveryDriver{
	UserID:   genericUser.ID,
	Cuil:     "0123456789",
	Age:      24,
	Vehicles: []string{"ABC123"},
}

var newDeliveryDriverForAdminUser = types.DeliveryDriver{
	UserID:   adminUser.ID,
	Cuil:     "9876543210",
	Age:      32,
	Vehicles: []string{"FFF111", "XYZ987"},
}

var newDeliveryDriverForNonExistingUser = types.DeliveryDriver{
	UserID:   100000000,
	Cuil:     "0123456789",
	Age:      20,
	Vehicles: []string{"ABC123"},
}

var invalidDeliveryDriver = types.DeliveryDriver{
	UserID:   2,
	Cuil:     "0123456789",
	Age:      7,
	Vehicles: []string{"TOO LONG VEHICLE ID"},
}

var users = []types.User{adminUser, genericUser}

var prices = map[uint]uint{
	250:  3,
	500:  5,
	1000: 10,
}

var newValidOrder types.Order = types.Order{
	Address:      "Calle 123",
	PaymentState: "pending",
}

var anotherNewValidOrder types.Order = types.Order{
	Address:      "Calle 789",
	PaymentState: "pending",
}

var newValidIceCreamTub types.IceCreamTub = types.IceCreamTub{
	Weight:  500,
	Flavors: []string{"ddl", "frt"},
}

var anotherNewValidIceCreamTub types.IceCreamTub = types.IceCreamTub{
	Weight:  1000,
	Flavors: []string{"ddl", "frt", "mrc"},
}

var iceCreamTubWithUnavailableFlavorsAndWeight types.IceCreamTub = types.IceCreamTub{
	Weight:  123,
	Flavors: []string{"non-existing-flavor"},
}

var invalidIceCreamTub = types.IceCreamTub{
	Weight:  0,
	Flavors: []string{"ddl", "frt"},
}

var validCreditCardPaymentRequest = payment.PaymentRequest{
	PaymentType: payment.CreditCardType,
	CreditCard:  &validCreditCard,
}

var validCreditCard = payment.CreditCard{
	Payment: payment.Payment{
		PaymentType: payment.CreditCardType,
	},
	CardNumber:      "1234567887654321",
	CVV:             "123",
	CardHolderName:  "John Doe",
	ExpirationMonth: "10",
	ExpirationYear:  "2050",
}

var validDigitalWalletPaymentRequest = payment.PaymentRequest{
	PaymentType:   payment.DigitalWalletType,
	CreditCard:    &validCreditCard,
	DigitalWallet: &validDigitalWallet,
}

var validDigitalWallet = payment.DigitalWallet{
	Payment: payment.Payment{
		PaymentType: payment.DigitalWalletType,
	},
	WalletID: "testing.alias",
}

var invalidPaymentRequest = payment.PaymentRequest{
	PaymentType: payment.DigitalWalletType,
	CreditCard:  &payment.CreditCard{},
}

// requestWithCookie receives the necessary data to make a request with a cookie value
func requestWithCookie(method, path string, structBody any, cookieName, cookieValue string) *httptest.ResponseRecorder {
	jsonBody, err := json.Marshal(structBody)
	if err != nil {
		panic(err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, strings.NewReader(string(jsonBody)))

	req.AddCookie(&http.Cookie{
		Name:  cookieName,
		Value: cookieValue,
	})

	router.ServeHTTP(w, req)
	return w
}

// requestToMakeAnOrder builds a request to make an order.
// Receives an order struct and the token from the user making it.
func requestToMakeAnOrder(orderStruct types.Order, userToken string) types.Order {
	w := requestWithCookie("POST", "/my-orders", orderStruct, "Authorization", userToken)
	var createdOrder types.Order
	err := json.Unmarshal(w.Body.Bytes(), &createdOrder)
	if err != nil {
		panic(err)
	}
	return createdOrder
}

// requestToAddATubToAnOrder builds a request to add a new tub to an order.
// Receives the tub struct, the order id and the user token who requested it.
func requestToAddATubToAnOrder(tubStruct types.IceCreamTub, orderID uint, userToken string) types.IceCreamTub {
	uri := fmt.Sprintf("/my-orders/%v/tubs", orderID)
	w := requestWithCookie("POST", uri, tubStruct, "Authorization", userToken)
	var createdTub types.IceCreamTub
	err := json.Unmarshal(w.Body.Bytes(), &createdTub)
	if err != nil {
		panic(err)
	}
	return createdTub
}

// requestToAddADeliveryDriver builds a request to add a new delivery driver.
// Receives the delivery driver struct and the admin token who requested it.
func requestToAddADeliveryDriver(deliveryDriverStruct types.DeliveryDriver, adminToken string) types.DeliveryDriver {
	w := requestWithCookie("POST", "/delivery-drivers", deliveryDriverStruct, "Authorization", adminToken)
	var createdDeliveryDriver types.DeliveryDriver
	err := json.Unmarshal(w.Body.Bytes(), &createdDeliveryDriver)
	if err != nil {
		panic(err)
	}
	return createdDeliveryDriver
}

// requestToAssignDeliveryDriverToOrder builds a request to assign a delivery driver to an order.
// Receives the delivery driver id, the order id and the admin token who requested it.
func requestToAssignDeliveryDriverToOrder(deliveryDriverID, orderID uint, authorizationToken string) {
	idStruct := struct {
		ID uint `json:"id"`
	}{deliveryDriverID}

	uri := fmt.Sprintf("/orders/%v/delivery-driver", orderID)
	_ = requestWithCookie("PUT", uri, idStruct, "Authorization", authorizationToken)
}
