package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"icecreamshop/internal/auth"
	"icecreamshop/internal/messageErrors"
	"icecreamshop/internal/types"
	"icecreamshop/internal/utils"
	"net/http"
	"testing"
)

/*********************************/
/***** USER MANAGEMENT TESTS *****/
/*********************************/

func TestSigningUpANewUser(t *testing.T) {
	setup()
	newUser := types.SignUpInput{
		Email:    "newuser@gmail.com",
		Name:     "hello",
		LastName: "world",
		Password: "valid-password",
	}

	w := requestWithCookie("POST", "/signup", newUser, "", "")

	expectedUser := types.User{
		ID:       3,
		Email:    "newuser@gmail.com",
		Name:     "hello",
		LastName: "world",
		Password: "",
	}
	var createdUser types.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, expectedUser, createdUser)

	obtainedUser, err := sv.Store.GetUserByID(3)
	obtainedUser.Password = ""
	assert.NoError(t, err)
	assert.True(t, expectedUser.IsEqualTo(obtainedUser))

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotSignUpAnUserWithEmptyName(t *testing.T) {
	setup()
	userWithoutName := types.User{
		Email:    "newuser@gmail.com",
		Name:     "",
		LastName: "world",
		Password: "valid-password",
	}

	w := requestWithCookie("POST", "/signup", userWithoutName, "", "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"%s"}`, messageErrors.FirstNameIsRequired))

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotSignUpAnUserWithEmptyLastName(t *testing.T) {
	setup()
	userWithoutLastName := types.User{
		Email:    "newuser@gmail.com",
		Name:     "hello",
		LastName: "",
		Password: "valid-password",
	}

	w := requestWithCookie("POST", "/signup", userWithoutLastName, "", "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"%s"}`, messageErrors.LastNameIsRequired))

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotSignUpAnUserWithEmptyEmail(t *testing.T) {
	setup()
	userWithoutEmail := types.User{
		Email:    "",
		Name:     "hello",
		LastName: "world",
		Password: "valid-password",
	}

	w := requestWithCookie("POST", "/signup", userWithoutEmail, "", "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"%s"}`, messageErrors.EmailIsRequired))

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotSignUpAnUserWithShortPassword(t *testing.T) {
	setup()
	userWithShortPassword := types.User{
		Email:    "newuser@gmail.com",
		Name:     "hello",
		LastName: "world",
		Password: "short",
	}

	w := requestWithCookie("POST", "/signup", userWithShortPassword, "", "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"error":"%s"}`, messageErrors.PasswordIsTooShort))

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotSignUpAnUserWithInvalidJsonFormat(t *testing.T) {
	setup()
	userWithInvalidJsonFormat := "I should be an user"

	w := requestWithCookie("POST", "/signup", userWithInvalidJsonFormat, "", "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidJsonFormat), w.Body.String())

	clearAndCloseConnection(t, sv.Store)
}

func TestANewTokenIsCreatedWhenLoggingInAnExistingUser(t *testing.T) {
	setup()
	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    adminUser.Email,
		Password: "admin123",
	}

	w := requestWithCookie("POST", "/login", credentials, "", "")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "null", w.Body.String())
	assert.Equal(t, 1, len(w.Result().Cookies()))
	assert.Equal(t, "Authorization", w.Result().Cookies()[0].Name)

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotLogInWithWrongCredentials(t *testing.T) {
	setup()
	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    "non-registered-email@gmail.com",
		Password: "admin123",
	}

	w := requestWithCookie("POST", "/login", credentials, "", "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidEmailOrPassword), w.Body.String())
	assert.Equal(t, 0, len(w.Result().Cookies()))

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotLogInWithInvalidJsonFormat(t *testing.T) {
	setup()
	credentials := "I should be a credential struct"

	w := requestWithCookie("POST", "/login", credentials, "", "")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidJsonFormat), w.Body.String())
	assert.Equal(t, 0, len(w.Result().Cookies()))

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotLoginIfAlreadyLoggedIn(t *testing.T) {
	setup()
	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    adminUser.Email,
		Password: "admin123",
	}

	token := auth.GenerateTokenFromUserEmail(credentials.Email)

	w := requestWithCookie("POST", "/login", credentials, "Authorization", token)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.AlreadyLoggedIn), w.Body.String())

	clearAndCloseConnection(t, sv.Store)
}

func TestGettingAllUsersWhenAnAdminIsLoggedIn(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("GET", "/users", nil, "Authorization", token)

	var obtainedUsers []types.User
	err := json.Unmarshal(w.Body.Bytes(), &obtainedUsers)

	var expectedUsers []types.User
	expectedUsers = sv.Store.GetAllUsers()
	for i := range expectedUsers {
		expectedUsers[i].Password = ""
	}

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, utils.SlicesAreEqual(expectedUsers, obtainedUsers))

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotGetAllUsersWhenAnAdminIsNotLoggedIn(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(genericUser.Email) //not an admin
	w := requestWithCookie("GET", "/users", nil, "Authorization", token)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestGettingAnUserByIDWhenAnAdminIsLoggedIn(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("GET", "/users/2", nil, "Authorization", token)

	var obtainedUser types.User
	err := json.Unmarshal(w.Body.Bytes(), &obtainedUser)

	expectedUser := users[1]
	expectedUser.Password = ""

	assert.NoError(t, err)
	assert.True(t, expectedUser.IsEqualTo(obtainedUser))
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotGetAnUserByIDWhenAnAdminIsNotLoggedIn(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(genericUser.Email) //not an admin
	w := requestWithCookie("GET", "/users/1", nil, "Authorization", token)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotGetUserByIDWhenIDIsNotAnInteger(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("GET", "/users/not-an-integer", nil, "Authorization", token)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotGetUserByIDWhenIDDoesNotExist(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(adminUser.Email) //an admin
	w := requestWithCookie("GET", "/users/100", nil, "Authorization", token)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.UserIDNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCanDeleteAnUserByID(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("DELETE", "/users/1", nil, "Authorization", token)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotDeleteAnUserByIDWhenUserPermissionIsNotAdmin(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("DELETE", "/users/1", nil, "Authorization", token)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotDeleteAnUserByIDWhenIDIsNotAnInteger(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("DELETE", "/users/hola", nil, "Authorization", token)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotDeleteAnUserByIDWhenIDDoesNotExist(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("DELETE", "/users/100", nil, "Authorization", token)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.UserIDNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCanPromoteAnUserToAdmin(t *testing.T) {
	setup()
	userInDB, _ := sv.Store.GetUserByID(2)
	assert.False(t, userInDB.IsAdmin())

	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("PUT", "/users/2/admin", nil, "Authorization", token)
	userInDB, _ = sv.Store.GetUserByID(2)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, userInDB.IsAdmin())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotPromoteANonExistingUserToAdmin(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("PUT", "/users/1000/admin", nil, "Authorization", token)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.UserIDNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserIDMustBeAnIntegerToPromoteThem(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("PUT", "/users/not-integer/admin", nil, "Authorization", token)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotPromoteAnAdminToAdmin(t *testing.T) {
	setup()
	userInDB, _ := sv.Store.GetUserByID(1)
	assert.True(t, userInDB.IsAdmin())

	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("PUT", "/users/1/admin", nil, "Authorization", token)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.UserIsAlreadyAnAdmin), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanGetTheirAccount(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("GET", "/my-account", nil, "Authorization", token)

	var expectedUser types.User
	expectedUser = genericUser
	expectedUser.Password = ""

	var obtainedUser types.User
	err := json.Unmarshal(w.Body.Bytes(), &obtainedUser)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedUser, obtainedUser)
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotGetAnAccountIfNotLoggedIn(t *testing.T) {
	setup()
	w := requestWithCookie("GET", "/my-account", nil, "", "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanDeleteTheirAccount(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("DELETE", "/my-account", nil, "Authorization", token)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
	_, err := sv.Store.GetUserByEmail("abcde@gmail.com")
	assert.EqualError(t, err, messageErrors.UserEmailNotFound)
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotDeleteTheirAccountIfNotLoggedIn(t *testing.T) {
	setup()
	w := requestWithCookie("DELETE", "/my-account", nil, "", "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanUpdateTheirAccount(t *testing.T) {
	setup()
	newUserData := types.User{
		Email:    "abcde@gmail.com",
		Name:     "bruce",
		LastName: "wayne",
	}
	userInDb, _ := sv.Store.GetUserByEmail(adminUser.Email)

	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("PUT", "/my-account", newUserData, "Authorization", token)

	var createdUser types.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)

	expectedUser := types.User{
		ID:          userInDb.ID,
		Email:       "abcde@gmail.com",
		Name:        "bruce",
		LastName:    "wayne",
		Orders:      userInDb.Orders,
		Permissions: userInDb.Permissions,
	}
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, createdUser)
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotUpdateTheirAccountWithInvalidData(t *testing.T) {
	setup()
	newUserData := types.User{
		Email:    "",
		Name:     "bruce",
		LastName: "wayne",
	}
	token := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("PUT", "/my-account", newUserData, "Authorization", token)

	_, err := sv.Store.GetUserByEmail("")

	assert.EqualError(t, err, messageErrors.UserEmailNotFound)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.EmailIsRequired), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotUpdateTheirAccountWithInvalidJsonFormat(t *testing.T) {
	setup()
	invalidUser := "not an user struct"

	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("PUT", "/my-account", invalidUser, "Authorization", token)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidJsonFormat), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

/*************************/
/***** FLAVORS TESTS *****/
/*************************/

func TestGettingAllFlavorsFromServer(t *testing.T) {
	setup()
	w := requestWithCookie("GET", "/flavors", nil, "", "")

	var actualFlavors []types.Flavor
	err := json.Unmarshal(w.Body.Bytes(), &actualFlavors)

	assert.NoError(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, flavors, actualFlavors)

	clearAndCloseConnection(t, sv.Store)
}

func TestGettingFlavorsFromServerFilteringByType(t *testing.T) {
	setup()
	w := requestWithCookie("GET", "/flavors?type=Chocolates", nil, "", "")

	expectedFlavors := []types.Flavor{flavorMRC}

	var actualFlavors []types.Flavor
	err := json.Unmarshal(w.Body.Bytes(), &actualFlavors)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, utils.SlicesAreEqual(expectedFlavors, actualFlavors))

	clearAndCloseConnection(t, sv.Store)
}

func TestGettingAFlavorByIDFromServer(t *testing.T) {
	setup()
	w := requestWithCookie("GET", "/flavors/ddl", nil, "", "")

	var actualFlavor types.Flavor
	err := json.Unmarshal(w.Body.Bytes(), &actualFlavor)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, flavorDDL, actualFlavor)

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotGetANonExistingFlavorFromServer(t *testing.T) {
	setup()
	w := requestWithCookie("GET", "/flavors/non-existing-flavor", nil, "", "")

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.FlavorNotFound), w.Body.String())

	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCanAddANewFlavor(t *testing.T) {
	setup()
	newFlavor := types.Flavor{
		ID: "ore", Name: "Oreo", Type: "Cremas",
	}

	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("POST", "/flavors", newFlavor, "Authorization", token)

	var obtainedFlavor types.Flavor
	err := json.Unmarshal(w.Body.Bytes(), &obtainedFlavor)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, newFlavor, obtainedFlavor)

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotAddANewFlavorWithInvalidJsonFormat(t *testing.T) {
	setup()
	newInvalidFlavor := "not a flavor struct"

	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("POST", "/flavors", newInvalidFlavor, "Authorization", token)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidJsonFormat), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotAddANewFlavorWithInvalidData(t *testing.T) {
	setup()
	newFlavor := types.Flavor{
		ID: "", Name: "Oreo", Type: "Cremas",
	}
	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("POST", "/flavors", newFlavor, "Authorization", token)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.FlavorIdIsRequired), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotAddANewFlavorWithAnExistingID(t *testing.T) {
	setup()
	newFlavor := types.Flavor{
		ID: "ddl", Name: "Oreo", Type: "Cremas",
	}
	token := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("POST", "/flavors", newFlavor, "Authorization", token)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.AlreadyExistingFlavor), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

/*****************************/
/***** USER ORDERS TESTS *****/
/*****************************/

func TestAnUserCanMakeAnOrder(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("POST", "/my-orders", newValidOrder, "Authorization", token)

	expectedOrder := types.Order{
		ID:           1,
		UserID:       genericUser.ID,
		Address:      "Calle 123",
		PaymentState: "pending",
	}

	var createdOrder types.Order
	err := json.Unmarshal(w.Body.Bytes(), &createdOrder)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.True(t, expectedOrder.IsEqualTo(createdOrder))

	orderInDB, err := sv.Store.GetOrderByID(1)
	assert.NoError(t, err)
	assert.True(t, expectedOrder.IsEqualTo(orderInDB))
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotMakeAnOrderWithInvalidJsonFormat(t *testing.T) {
	setup()
	ordersBeforeRequest := sv.Store.GetAllOrdersByUserEmail(genericUser.Email)
	newInvalidOrder := "Not an order struct"
	token := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("POST", "/my-orders", newInvalidOrder, "Authorization", token)

	orders := sv.Store.GetAllOrdersByUserEmail(genericUser.Email)
	assert.Equal(t, len(ordersBeforeRequest), len(orders))
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidJsonFormat), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotMakeAnOrderWithEmptyAddress(t *testing.T) {
	setup()
	newOrder := types.Order{
		Address: "",
	}
	token := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("POST", "/my-orders", newOrder, "Authorization", token)

	orders := sv.Store.GetAllOrdersByUserEmail(genericUser.Email)
	assert.Equal(t, 0, len(orders))
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.AddressIsRequired), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanGetTheirOrderByID(t *testing.T) {
	setup()
	//creating order
	token := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("POST", "/my-orders", newValidOrder, "Authorization", token)
	var createdOrder types.Order
	err := json.Unmarshal(w.Body.Bytes(), &createdOrder)

	//getting order
	uri := fmt.Sprintf("/my-orders/%v", createdOrder.ID)
	w = requestWithCookie("GET", uri, nil, "Authorization", token)
	var obtainedOrder types.Order
	err = json.Unmarshal(w.Body.Bytes(), &obtainedOrder)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, err)
	assert.Equal(t, createdOrder, obtainedOrder)
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserMustBeLoggedInToGetTheirOrder(t *testing.T) {
	setup()
	w := requestWithCookie("GET", "/my-orders/1", nil, "", "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotGetAnOrderWhenIDIsNotAnInteger(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("GET", "/my-orders/not-an-integer", nil, "Authorization", token)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotGetAnotherUserOrder(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	tokenAnotherUser := auth.GenerateTokenFromUserEmail(adminUser.Email)

	//creating order
	w := requestWithCookie("POST", "/my-orders", newValidOrder, "Authorization", tokenAnotherUser)
	var createdOrder types.Order
	err := json.Unmarshal(w.Body.Bytes(), &createdOrder)

	//getting order
	uri := fmt.Sprintf("/my-orders/%v", createdOrder.ID)
	w = requestWithCookie("GET", uri, nil, "Authorization", tokenUser)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotGetANonExistingOrder(t *testing.T) {
	setup()
	token := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("GET", "/my-orders/1000", nil, "Authorization", token)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanGetAllTheirOrders(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	firstOrder := requestToMakeAnOrder(newValidOrder, tokenUser)
	secondOrder := requestToMakeAnOrder(anotherNewValidOrder, tokenUser)

	w := requestWithCookie("GET", "/my-orders", nil, "Authorization", tokenUser)
	var obtainedOrders []types.Order
	err := json.Unmarshal(w.Body.Bytes(), &obtainedOrders)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(obtainedOrders))
	assert.Equal(t, firstOrder, obtainedOrders[0])
	assert.Equal(t, secondOrder, obtainedOrders[1])
}

func TestAnUserMustBeLoggedInToGetAllTheirOrders(t *testing.T) {
	setup()
	w := requestWithCookie("GET", "/my-orders", nil, "", "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserDoesNotHaveAnotherUsersOrdersWhenGettingAllOrders(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	tokenAnotherUser := auth.GenerateTokenFromUserEmail(adminUser.Email)
	_ = requestToMakeAnOrder(newValidOrder, tokenAnotherUser)

	w := requestWithCookie("GET", "/my-orders", nil, "Authorization", tokenUser)
	var obtainedOrders []types.Order
	err := json.Unmarshal(w.Body.Bytes(), &obtainedOrders)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(obtainedOrders))
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanUpdateTheirOrderByID(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)

	updatedData := types.Order{
		Address:      "Calle 1000",
		PaymentState: "paid",
	}
	uri := fmt.Sprintf("/my-orders/%v", order.ID)
	w := requestWithCookie("PUT", uri, updatedData, "Authorization", tokenUser)
	var updatedOrder types.Order
	err := json.Unmarshal(w.Body.Bytes(), &updatedOrder)

	expectedOrder := order
	expectedOrder.Address = updatedData.Address
	expectedOrder.PaymentState = updatedData.PaymentState

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, updatedOrder, expectedOrder)
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotUpdateTheirOrderIfNotLoggedIn(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)

	updatedData := types.Order{
		Address:      "Calle 1000",
		PaymentState: "paid",
	}
	uri := fmt.Sprintf("/my-orders/%v", order.ID)
	w := requestWithCookie("PUT", uri, updatedData, "", "")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())

	orderInDB, _ := sv.Store.GetOrderByID(order.ID)
	assert.NotEqual(t, orderInDB.Address, updatedData.Address)
	assert.NotEqual(t, orderInDB.PaymentState, updatedData.PaymentState)

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotUpdateAnOrderWhenIDIsNotAnInteger(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	updatedData := types.Order{
		Address:      "Calle 1000",
		PaymentState: "paid",
	}
	w := requestWithCookie("PUT", "/my-orders/not-an-integer", updatedData, "Authorization", tokenUser)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotUpdateAnOrderWithInvalidJsonFormat(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)

	updatedData := "not an order struct"
	uri := fmt.Sprintf("/my-orders/%v", order.ID)
	w := requestWithCookie("PUT", uri, updatedData, "Authorization", tokenUser)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidJsonFormat), w.Body.String())

	orderInDB, _ := sv.Store.GetOrderByID(order.ID)
	assert.Equal(t, orderInDB.Address, newValidOrder.Address)
	assert.Equal(t, orderInDB.PaymentState, newValidOrder.PaymentState)

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotUpdateAnOrderWithEmptyData(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)

	updatedData := types.Order{
		Address:      "",
		PaymentState: "paid",
	}
	uri := fmt.Sprintf("/my-orders/%v", order.ID)
	w := requestWithCookie("PUT", uri, updatedData, "Authorization", tokenUser)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.AddressIsRequired), w.Body.String())

	orderInDB, _ := sv.Store.GetOrderByID(order.ID)
	assert.Equal(t, orderInDB.Address, newValidOrder.Address)
	assert.Equal(t, orderInDB.PaymentState, newValidOrder.PaymentState)

	clearAndCloseConnection(t, sv.Store)
}

func TestCannotUpdateAnOrderThatDoesNotExist(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	updatedData := types.Order{
		Address:      "Calle 1000",
		PaymentState: "paid",
	}
	w := requestWithCookie("PUT", "/my-orders/1000", updatedData, "Authorization", tokenUser)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestCannotUpdateAnOrderThatDoesNotBelongToUser(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	tokenAnotherUser := auth.GenerateTokenFromUserEmail(adminUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenAnotherUser)
	updatedData := types.Order{
		Address:      "Calle 1000",
		PaymentState: "paid",
	}
	uri := fmt.Sprintf("/my-orders/%v", order.ID)
	w := requestWithCookie("PUT", uri, updatedData, "Authorization", tokenUser)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

/********************************/
/***** ICE CREAM TUBS TESTS *****/
/********************************/

func TestAnUserCanAddAnIceCreamTubToTheirOrder(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)

	uri := fmt.Sprintf("/my-orders/%v/tubs", order.ID)
	w := requestWithCookie("POST", uri, newValidIceCreamTub, "Authorization", tokenUser)

	orderTubs, err := sv.Store.GetIceCreamTubsByOrderID(order.ID)
	var tubObtained types.IceCreamTub
	_ = json.Unmarshal(w.Body.Bytes(), &tubObtained)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, 1, len(orderTubs))
	assert.Equal(t, orderTubs[0], tubObtained)

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotAddAnIceCreamTubToAnotherUsersOrder(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	tokenAnotherUser := auth.GenerateTokenFromUserEmail(adminUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenAnotherUser)

	uri := fmt.Sprintf("/my-orders/%v/tubs", order.ID)
	w := requestWithCookie("POST", uri, newValidIceCreamTub, "Authorization", tokenUser)

	tubs, _ := sv.Store.GetIceCreamTubsByOrderID(order.ID)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())
	assert.Equal(t, 0, len(tubs))

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotAddAnIceCreamTubWhenOrderIDIsNotAnInteger(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie(
		"POST",
		"/my-orders/not-an-integer/tubs",
		newValidIceCreamTub,
		"Authorization",
		tokenUser,
	)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotAddAnIceCreamTubWithInvalidJsonFormat(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	uri := fmt.Sprintf("/my-orders/%v/tubs", order.ID)
	w := requestWithCookie("POST", uri, "not an ice cream tub struct", "Authorization", tokenUser)

	tubs, _ := sv.Store.GetIceCreamTubsByOrderID(order.ID)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidJsonFormat), w.Body.String())
	assert.Equal(t, 0, len(tubs))

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotAddAnIceCreamTubWithInvalidData(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	uri := fmt.Sprintf("/my-orders/%v/tubs", order.ID)
	w := requestWithCookie("POST", uri, invalidIceCreamTub, "Authorization", tokenUser)

	tubs, _ := sv.Store.GetIceCreamTubsByOrderID(order.ID)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.WeightCannotBeZero), w.Body.String())
	assert.Equal(t, 0, len(tubs))

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotAddAnIceCreamTubWithANonExistingOrderID(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	uri := fmt.Sprintf("/my-orders/%v/tubs", 1000)
	w := requestWithCookie("POST", uri, newValidIceCreamTub, "Authorization", tokenUser)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotAddAnIceCreamTubWithNonAvailableFlavorsOrWeight(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	uri := fmt.Sprintf("/my-orders/%v/tubs", order.ID)
	w := requestWithCookie("POST", uri, iceCreamTubWithUnavailableFlavorsAndWeight, "Authorization", tokenUser)

	tubs, _ := sv.Store.GetIceCreamTubsByOrderID(order.ID)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, 0, len(tubs))

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanGetTheirIceCreamTubFromTheirOrder(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	uri := fmt.Sprintf("/my-orders/%v/tubs", order.ID)
	w := requestWithCookie("POST", uri, newValidIceCreamTub, "Authorization", tokenUser)
	w = requestWithCookie("POST", uri, anotherNewValidIceCreamTub, "Authorization", tokenUser)
	w = requestWithCookie("GET", uri, nil, "Authorization", tokenUser)

	tubsInDB, _ := sv.Store.GetIceCreamTubsByOrderID(order.ID)

	var obtainedTubs []types.IceCreamTub
	err := json.Unmarshal(w.Body.Bytes(), &obtainedTubs)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 2, len(obtainedTubs))
	assert.ElementsMatch(t, tubsInDB, obtainedTubs)

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotGetTheirIceCreamTubsFromOtherOrder(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	tokenAnotherUser := auth.GenerateTokenFromUserEmail(adminUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenAnotherUser)
	uri := fmt.Sprintf("/my-orders/%v/tubs", order.ID)
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotGetTheirIceCreamTubsWhenOrderIDIsNotAnInteger(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	uri := fmt.Sprintf("/my-orders/%v/tubs", "NOT-AN-INTEGER")
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanDeleteAnIceCreamTubFromTheirOrder(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	tub := requestToAddATubToAnOrder(newValidIceCreamTub, order.ID, tokenUser)
	tubsInDBBeforeDelete, _ := sv.Store.GetIceCreamTubsByOrderID(order.ID)

	uri := fmt.Sprintf("/my-orders/%v/tubs/%v", order.ID, tub.ID)
	w := requestWithCookie("DELETE", uri, nil, "Authorization", tokenUser)

	tubsInDBAfterDelete, _ := sv.Store.GetIceCreamTubsByOrderID(order.ID)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, 1, len(tubsInDBBeforeDelete))
	assert.Equal(t, 0, len(tubsInDBAfterDelete))
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotDeleteAnIceCreamTubFromAnotherUsersOrder(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	tokenAnotherUser := auth.GenerateTokenFromUserEmail(adminUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenAnotherUser)
	tub := requestToAddATubToAnOrder(newValidIceCreamTub, order.ID, tokenAnotherUser)
	tubsBeforeDeletion, _ := sv.Store.GetIceCreamTubsByOrderID(order.ID)

	uri := fmt.Sprintf("/my-orders/%v/tubs/%v", order.ID, tub.ID)
	w := requestWithCookie("DELETE", uri, nil, "Authorization", tokenUser)
	tubsAfterDeletion, _ := sv.Store.GetIceCreamTubsByOrderID(order.ID)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())
	assert.ElementsMatch(t, tubsBeforeDeletion, tubsAfterDeletion)

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotDeleteAnIceCreamTubWhenOrderIDIsNotAnInteger(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("DELETE", "/my-orders/NOT-AN-INTEGER/tubs/1", nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotDeleteAnIceCreamTubWhenTubIDIsNotAnInteger(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("DELETE", "/my-orders/1/tubs/NOT-AN-INTEGER", nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotDeleteAnIceCreamTubWhenTubIDDoesNotExist(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	uri := fmt.Sprintf("/my-orders/%v/tubs/1000000", order.ID)
	w := requestWithCookie("DELETE", uri, nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.IceCreamTubNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

/**********************************/
/***** DELIVERY-DRIVERS TESTS *****/
/**********************************/

func TestAnAdminCanAddANewDeliveryDriver(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("POST", "/delivery-drivers", newDeliveryDriverForGenericUser, "Authorization", tokenAdmin)

	deliveryDriverInDB, err := sv.Store.GetDeliveryDriverByID(newDeliveryDriverForGenericUser.UserID)
	assert.Nil(t, err)

	var deliveryDriverObtained types.DeliveryDriver
	err = json.Unmarshal(w.Body.Bytes(), &deliveryDriverObtained)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.True(t, newDeliveryDriverForGenericUser.IsEqualTo(deliveryDriverInDB))
	assert.Equal(t, newDeliveryDriverForGenericUser, deliveryDriverObtained)
	clearAndCloseConnection(t, sv.Store)
}

func TestANonAdminCannotAddANewDeliveryDriver(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("POST", "/delivery-drivers", newDeliveryDriverForGenericUser, "Authorization", tokenUser)

	_, err := sv.Store.GetDeliveryDriverByID(newDeliveryDriverForGenericUser.UserID)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	assert.EqualError(t, err, messageErrors.DeliveryDriverNotFound)

	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotAddANewDeliveryDriverWithInvalidJsonFormat(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("POST", "/delivery-drivers", "NOT A DELIVERY DRIVER STRUCT", "Authorization", tokenAdmin)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidJsonFormat), w.Body.String())

	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotAddAnInvalidNewDeliveryDriver(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("POST", "/delivery-drivers", invalidDeliveryDriver, "Authorization", tokenAdmin)

	_, err := sv.Store.GetDeliveryDriverByID(newDeliveryDriverForGenericUser.UserID)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.AgeMustBeGreaterThan18), w.Body.String())
	assert.EqualError(t, err, messageErrors.DeliveryDriverNotFound)

	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotAddANewDeliveryDriverFromANonExistingUser(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	w := requestWithCookie("POST", "/delivery-drivers", newDeliveryDriverForNonExistingUser, "Authorization", tokenAdmin)

	_, err := sv.Store.GetDeliveryDriverByID(newDeliveryDriverForNonExistingUser.UserID)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.UserIDNotFound), w.Body.String())
	assert.EqualError(t, err, messageErrors.DeliveryDriverNotFound)
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCanGetAllDeliveryDrivers(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	deliveryDriver := requestToAddADeliveryDriver(newDeliveryDriverForGenericUser, tokenAdmin)
	anotherDeliveryDriver := requestToAddADeliveryDriver(newDeliveryDriverForAdminUser, tokenAdmin)

	w := requestWithCookie("GET", "/delivery-drivers", nil, "Authorization", tokenAdmin)

	expectedDeliveryDrivers := []types.DeliveryDriver{deliveryDriver, anotherDeliveryDriver}
	var deliveryDriversObtained []types.DeliveryDriver
	err := json.Unmarshal(w.Body.Bytes(), &deliveryDriversObtained)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, utils.SlicesAreEqual(expectedDeliveryDrivers, deliveryDriversObtained))
	clearAndCloseConnection(t, sv.Store)
}

func TestANonAdminCannotGetAllDeliveryDrivers(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("GET", "/delivery-drivers", nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())

	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCanGetADeliveryDriverByID(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	deliveryDriver := requestToAddADeliveryDriver(newDeliveryDriverForGenericUser, tokenAdmin)

	uri := fmt.Sprintf("/delivery-drivers/%v", deliveryDriver.UserID)
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenAdmin)

	var deliveryDriverObtained types.DeliveryDriver
	err := json.Unmarshal(w.Body.Bytes(), &deliveryDriverObtained)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, deliveryDriver, deliveryDriverObtained)

	clearAndCloseConnection(t, sv.Store)
}

func TestANonAdminCannotGetADeliveryDriverByID(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenNotAnAdmin := auth.GenerateTokenFromUserEmail(genericUser.Email)
	deliveryDriver := requestToAddADeliveryDriver(newDeliveryDriverForGenericUser, tokenAdmin)
	uri := fmt.Sprintf("/delivery-drivers/%v", deliveryDriver.UserID)

	w := requestWithCookie("GET", uri, nil, "Authorization", tokenNotAnAdmin)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotGetADeliveryDriverWhenIDIsNotAnInteger(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	uri := fmt.Sprintf("/delivery-drivers/%v", "NOT-AN-INTEGER")
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenAdmin)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotGetADeliveryDriverWhenIDDoesNotExist(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	uri := fmt.Sprintf("/delivery-drivers/%v", 10000000)
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenAdmin)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.DeliveryDriverNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanUpdateTheirDeliveryDriverData(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	_ = requestToAddADeliveryDriver(newDeliveryDriverForGenericUser, tokenAdmin)
	updatedDeliveryDriverData := types.DeliveryDriver{
		Cuil:     newDeliveryDriverForGenericUser.Cuil,
		Age:      25,
		Vehicles: []string{"ABC123", "AAA111"},
	}
	w := requestWithCookie("PUT", "/my-account/delivery-driver", updatedDeliveryDriverData, "Authorization", tokenUser)

	var deliveryDriverObtained types.DeliveryDriver
	err := json.Unmarshal(w.Body.Bytes(), &deliveryDriverObtained)

	expectedDeliveryDriver := types.DeliveryDriver{
		UserID:   genericUser.ID,
		Cuil:     updatedDeliveryDriverData.Cuil,
		Age:      25,
		Vehicles: []string{"ABC123", "AAA111"},
	}

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedDeliveryDriver, deliveryDriverObtained)
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotUpdateTheirDeliveryDriverDataWithInvalidJsonFormat(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	_ = requestToAddADeliveryDriver(newDeliveryDriverForGenericUser, tokenAdmin)
	w := requestWithCookie("PUT", "/my-account/delivery-driver", "Not a delivery driver struct", "Authorization", tokenUser)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidJsonFormat), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotUpdateTheirDeliveryDriverDataWithInvalidData(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	deliveryDriver := requestToAddADeliveryDriver(newDeliveryDriverForGenericUser, tokenAdmin)
	w := requestWithCookie("PUT", "/my-account/delivery-driver", invalidDeliveryDriver, "Authorization", tokenUser)

	deliveryDriverInDB, _ := sv.Store.GetDeliveryDriverByID(deliveryDriver.UserID)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.AgeMustBeGreaterThan18), w.Body.String())
	assert.True(t, deliveryDriver.IsEqualTo(deliveryDriverInDB))

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanDeleteThemselvesAsDeliveryDriver(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	_ = requestToAddADeliveryDriver(newDeliveryDriverForGenericUser, tokenAdmin)
	w := requestWithCookie("DELETE", "/my-account/delivery-driver", nil, "Authorization", tokenUser)

	_, err := sv.Store.GetDeliveryDriverByID(genericUser.ID)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
	assert.EqualError(t, err, messageErrors.DeliveryDriverNotFound)
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotDeleteThemselvesAsDeliveryDriverWhenTheyAreNotDrivers(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("DELETE", "/my-account/delivery-driver", nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

/*****************************************/
/***** ORDERS ADMIN MANAGEMENT TESTS *****/
/*****************************************/

func TestAnAdminCanGetAllOrders(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	_ = requestToMakeAnOrder(newValidOrder, tokenUser)
	_ = requestToMakeAnOrder(anotherNewValidOrder, tokenAdmin)
	w := requestWithCookie("GET", "/orders", nil, "Authorization", tokenAdmin)

	ordersInDB := sv.Store.GetAllOrders()

	var obtainedOrders []types.Order
	err := json.Unmarshal(w.Body.Bytes(), &obtainedOrders)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.ElementsMatch(t, ordersInDB, obtainedOrders)
	clearAndCloseConnection(t, sv.Store)
}

func TestANonAdminUserCannotGetAllOrders(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	w := requestWithCookie("GET", "/orders", nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCanGetAnyOrderByID(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	uri := fmt.Sprintf("/orders/%v", order.ID)
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenAdmin)

	orderInDB, err := sv.Store.GetOrderByID(order.ID)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, order.ID, orderInDB.ID)
	clearAndCloseConnection(t, sv.Store)
}

func TestANonAdminUserCannotGetAnyOrderByID(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenAdmin)
	uri := fmt.Sprintf("/orders/%v", order.ID)
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Empty(t, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotGetAnyOrderWhenOrderIDIsNotAnInteger(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	uri := fmt.Sprintf("/orders/%v", "NOT AN INTEGER")
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenAdmin)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotGetAnyOrderWhenOrderIDDoesNotExist(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	uri := fmt.Sprintf("/orders/%v", "1000000")
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenAdmin)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCanAssignADeliveryDriverToAnOrder(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	_ = requestToAddADeliveryDriver(newDeliveryDriverForAdminUser, tokenAdmin)
	deliveryDriverID := struct {
		ID uint `json:"id"`
	}{adminUser.ID}

	uri := fmt.Sprintf("/orders/%v/delivery-driver", order.ID)
	w := requestWithCookie("PUT", uri, deliveryDriverID, "Authorization", tokenAdmin)

	orderInDB, err := sv.Store.GetOrderByID(order.ID)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
	assert.Equal(t, deliveryDriverID.ID, orderInDB.ID)
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotAssignADeliveryDriverWhenOrderIDIsNotAnInteger(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	_ = requestToAddADeliveryDriver(newDeliveryDriverForAdminUser, tokenAdmin)
	deliveryDriverID := struct {
		ID uint `json:"id"`
	}{adminUser.ID}

	uri := fmt.Sprintf("/orders/%v/delivery-driver", "NOT AN INTEGER")
	w := requestWithCookie("PUT", uri, deliveryDriverID, "Authorization", tokenAdmin)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotAssignADeliveryDriverWithInvalidJsonFormat(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	_ = requestToAddADeliveryDriver(newDeliveryDriverForAdminUser, tokenAdmin)

	uri := fmt.Sprintf("/orders/%v/delivery-driver", order.ID)
	w := requestWithCookie("PUT", uri, "Invalid Json Format", "Authorization", tokenAdmin)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidJsonFormat), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotAssignANonExistingDeliveryDriverToAnOrder(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	deliveryDriverID := struct {
		ID uint `json:"id"`
	}{100000}

	uri := fmt.Sprintf("/orders/%v/delivery-driver", order.ID)
	w := requestWithCookie("PUT", uri, deliveryDriverID, "Authorization", tokenAdmin)

	orderInDB, err := sv.Store.GetOrderByID(order.ID)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.DeliveryDriverNotFound), w.Body.String())
	assert.NotEqual(t, 100000, orderInDB.ID)
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCanDeleteAnAssignedDeliveryDriverFromAnOrder(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	deliveryDriver := requestToAddADeliveryDriver(newDeliveryDriverForAdminUser, tokenAdmin)
	requestToAssignDeliveryDriverToOrder(deliveryDriver.UserID, order.ID, tokenAdmin)

	uri := fmt.Sprintf("/orders/%v/delivery-driver", order.ID)
	w := requestWithCookie("DELETE", uri, nil, "Authorization", tokenAdmin)

	orderInDB, err := sv.Store.GetOrderByID(order.ID)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
	assert.Equal(t, uint(0), orderInDB.DeliveryDriverID)
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotDeleteAnAssignedDeliveryWhenOrderIDIsNotAnInteger(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	deliveryDriver := requestToAddADeliveryDriver(newDeliveryDriverForAdminUser, tokenAdmin)
	requestToAssignDeliveryDriverToOrder(deliveryDriver.UserID, order.ID, tokenAdmin)

	uri := fmt.Sprintf("/orders/%v/delivery-driver", "NOT AN INTEGER")
	w := requestWithCookie("DELETE", uri, nil, "Authorization", tokenAdmin)

	orderInDB, err := sv.Store.GetOrderByID(order.ID)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	assert.Equal(t, deliveryDriver.UserID, orderInDB.DeliveryDriverID)
	clearAndCloseConnection(t, sv.Store)
}

func TestAnAdminCannotDeleteAnAssignedDeliveryFromANonExistingOrder(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)

	uri := fmt.Sprintf("/orders/%v/delivery-driver", 100000)
	w := requestWithCookie("DELETE", uri, nil, "Authorization", tokenAdmin)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanGetTheAssignedDeliveryDriverFromTheirOrder(t *testing.T) {
	setup()
	tokenAdmin := auth.GenerateTokenFromUserEmail(adminUser.Email)
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	deliveryDriver := requestToAddADeliveryDriver(newDeliveryDriverForAdminUser, tokenAdmin)
	requestToAssignDeliveryDriverToOrder(deliveryDriver.UserID, order.ID, tokenAdmin)

	uri := fmt.Sprintf("/my-orders/%v/delivery-driver", order.ID)
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenUser)

	expectedID := fmt.Sprint(deliveryDriver.UserID)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedID, w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotGetTheAssignedDeliveryDriverFromANonExistingOrder(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)

	uri := fmt.Sprintf("/my-orders/%v/delivery-driver", 100000)
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotGetTheAssignedDeliveryDriverWhenOrderIDIsNotAnInteger(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)

	uri := fmt.Sprintf("/my-orders/%v/delivery-driver", "NOT AN INTEGER")
	w := requestWithCookie("GET", uri, nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

/********************************/
/***** ORDERS PAYMENT TESTS *****/
/********************************/

func TestAnUserCanPayTheirOrderUsingAValidCreditCard(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	_ = requestToAddATubToAnOrder(newValidIceCreamTub, order.ID, tokenUser)
	_ = requestToAddATubToAnOrder(anotherNewValidIceCreamTub, order.ID, tokenUser)

	uri := fmt.Sprintf("/my-orders/%v/pay", order.ID)
	w := requestWithCookie("POST", uri, validCreditCardPaymentRequest, "Authorization", tokenUser)

	orderInDB, err := sv.Store.GetOrderByID(order.ID)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Equal(t, "paid", orderInDB.PaymentState)

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCanPayTheirOrderUsingAValidDigitalWallet(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	_ = requestToAddATubToAnOrder(newValidIceCreamTub, order.ID, tokenUser)
	_ = requestToAddATubToAnOrder(anotherNewValidIceCreamTub, order.ID, tokenUser)

	uri := fmt.Sprintf("/my-orders/%v/pay", order.ID)
	w := requestWithCookie("POST", uri, validDigitalWalletPaymentRequest, "Authorization", tokenUser)

	orderInDB, err := sv.Store.GetOrderByID(order.ID)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Equal(t, "paid", orderInDB.PaymentState)

	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotPayANonExistingOrder(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	uri := fmt.Sprintf("/my-orders/%v/pay", 100000)
	w := requestWithCookie("POST", uri, nil, "Authorization", tokenUser)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.OrderNotFound), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotPayWhenOrderIDIsNotAnInteger(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	uri := fmt.Sprintf("/my-orders/%v/pay", "NOT AN INTEGER")
	w := requestWithCookie("POST", uri, nil, "Authorization", tokenUser)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.MustBeAnInteger), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotPayWhenPaymentDataHasNotAValidJsonFormat(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	uri := fmt.Sprintf("/my-orders/%v/pay", order.ID)

	w := requestWithCookie("POST", uri, "NOT A PAYMENT STRUCT", "Authorization", tokenUser)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidJsonFormat), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}

func TestAnUserCannotPayWithInvalidPaymentData(t *testing.T) {
	setup()
	tokenUser := auth.GenerateTokenFromUserEmail(genericUser.Email)
	order := requestToMakeAnOrder(newValidOrder, tokenUser)
	uri := fmt.Sprintf("/my-orders/%v/pay", order.ID)
	w := requestWithCookie("POST", uri, invalidPaymentRequest, "Authorization", tokenUser)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, utils.CreateJsonSingletonString("error", messageErrors.InvalidPaymentData), w.Body.String())
	clearAndCloseConnection(t, sv.Store)
}
