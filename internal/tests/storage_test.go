package tests

import (
	"github.com/stretchr/testify/assert"
	"icecreamshop/internal/messageErrors"
	"icecreamshop/internal/types"
	"icecreamshop/internal/utils"
	"slices"
	"testing"
)

/*************************/
/***** FLAVORS TESTS *****/
/*************************/

func TestGettingAllFlavors(t *testing.T) {
	store := newStorage(flavors, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	allFlavors := store.GetFlavors()

	assert.Equal(t, flavors, allFlavors)
}

func TestAddingANewFlavor(t *testing.T) {
	store := newStorage([]types.Flavor{}, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newFlavor := types.Flavor{ID: "ddl", Name: "Dulce de leche", Type: "Dulce de leches"}

	err := store.AddFlavor(newFlavor)
	allFlavors := store.GetFlavors()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(allFlavors))
	assert.True(t, slices.Contains(allFlavors, newFlavor), "New flavor should be in the collection")
}

func TestCannotAddAnExistingFlavor(t *testing.T) {
	store := newStorage(flavors, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newFlavor := types.Flavor{ID: "ddl", Name: "Dulce de leche", Type: "Dulce de leches"}

	err := store.AddFlavor(newFlavor)
	allFlavors := store.GetFlavors()

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.AlreadyExistingFlavor)
	assert.Equal(t, flavors, allFlavors)
}

func TestFilteringFlavorsByType(t *testing.T) {
	store := newStorage(flavors, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	expectedFlavors := []types.Flavor{
		{ID: "mrc", Name: "Chocolate marroc", Type: "Chocolates"},
	}
	actualFlavors := store.GetFlavorsByType("Chocolates")

	assert.Equal(t, expectedFlavors, actualFlavors)
}

func TestGettingFlavorByID(t *testing.T) {
	store := newStorage(flavors, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	expectedFlavor := types.Flavor{ID: "mrc", Name: "Chocolate marroc", Type: "Chocolates"}
	actualFlavor, err := store.GetFlavorByID("mrc")

	assert.NoError(t, err)
	assert.Equal(t, expectedFlavor, actualFlavor)
}

func TestCannotGetANonExistentFlavor(t *testing.T) {
	store := newStorage(flavors, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	_, err := store.GetFlavorByID("hello")

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.FlavorNotFound)
}

/*********************************/
/***** USER MANAGEMENT TESTS *****/
/*********************************/

func TestGettingAllUsers(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	allUsers := store.GetAllUsers()
	assert.Equal(t, users, allUsers)
}

func TestAddingANewValidUser(t *testing.T) {
	store := newStorage([]types.Flavor{}, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newUser := types.User{
		Email:    "abcde@gmail.com",
		Name:     "hello",
		LastName: "world",
		Password: "admin",
	}

	err := store.SignUpUser(&newUser)
	allUsers := store.GetAllUsers()

	assert.NoError(t, err)
	assert.True(t, utils.SliceContains(allUsers, newUser), "New user should be in the collection")
}

func TestCannotAddANewUserWithAnExistingEmail(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newUser := types.User{
		Email:    "abcde@gmail.com",
		Name:     "hello",
		LastName: "world",
		Password: "aasdjasoidjsd",
	}

	err := store.SignUpUser(&newUser)
	allUsers := store.GetAllUsers()

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.EmailAlreadyExists)
	assert.Equal(t, users, allUsers)
}

func TestAnUserHasNoOrdersWhenCreated(t *testing.T) {
	store := newStorage([]types.Flavor{}, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newUser := types.User{
		Email:    "abcde@gmail.com",
		Name:     "hello",
		LastName: "world",
		Password: "aasdjasoidjsd",
	}

	err := store.SignUpUser(&newUser)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(newUser.Orders))
}

func TestAnUserHasNoPermissionsWhenCreated(t *testing.T) {
	store := newStorage([]types.Flavor{}, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newUser := types.User{
		Email:    "abcde@gmail.com",
		Name:     "hello",
		LastName: "world",
		Password: "aasdjasoidjsd",
	}

	err := store.SignUpUser(&newUser)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(newUser.Permissions))
}

func TestAnUserHasAnIDAssignedWhenCreated(t *testing.T) {
	store := newStorage([]types.Flavor{}, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newUser := types.User{
		Email:    "abcde@gmail.com",
		Name:     "hello",
		LastName: "world",
		Password: "aasdjasoidjsd",
	}

	err := store.SignUpUser(&newUser)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), newUser.ID)
}

func TestGettingAnUserByEmail(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	user, err := store.GetUserByEmail("abcde@gmail.com")

	assert.NoError(t, err)
	assert.Equal(t, users[0], user)
}

func TestCannotGetAnUserByANonExistingEmail(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	_, err := store.GetUserByEmail("hello@gmail.com")

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.UserEmailNotFound)
}

func TestGettingAnUserByID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	user, err := store.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, users[0], user)
}

func TestCannotGetAnUserByANonExistingID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	_, err := store.GetUserByID(100)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.UserIDNotFound)
}

func TestUpdatingAnUserByID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newUserData := types.User{
		ID:       1,
		Email:    "abcde@gmail.com",
		Name:     "bruce",
		LastName: "wayne",
	}

	updatedUser, err := store.UpdateUser(newUserData)
	actualUser, _ := store.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, updatedUser, actualUser)
}

func TestCannotUpdateAnUserByANonExistingID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newUserData := types.User{
		ID:       100,
		Email:    "abcde@gmail.com",
		Name:     "bruce",
		LastName: "wayne",
	}

	_, err := store.UpdateUser(newUserData)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.UserIDNotFound)
}

func TestDeletingAnUserByID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	err := store.DeleteUserByID(1)
	allUsers := store.GetAllUsers()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(allUsers))
	assert.Equal(t, users[1:], allUsers)
}

func TestCannotDeleteAnUserByANonExistingID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	err := store.DeleteUserByID(100)
	allUsers := store.GetAllUsers()

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.UserIDNotFound)
	assert.Equal(t, users, allUsers)
}

func TestPromotingAnUserByIDToAdmin(t *testing.T) {
	store := newStorage([]types.Flavor{}, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newUser := types.User{
		Email:    "abcde@gmail.com",
		Name:     "bruce",
		LastName: "wayne",
		Password: "aasdjasoidjsd",
	}

	err := store.SignUpUser(&newUser)
	err = store.PromoteUserToAdmin(1)
	user, _ := store.GetUserByID(1)

	assert.NoError(t, err)
	assert.True(t, user.IsAdmin())
}

func TestCannotPromoteAnUserByANonExistingID(t *testing.T) {
	store := newStorage([]types.Flavor{}, []types.User{}, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	err := store.PromoteUserToAdmin(100)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.UserIDNotFound)
}

func TestCannotPromoteAnAdminToAdmin(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	err := store.PromoteUserToAdmin(1)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.UserIsAlreadyAnAdmin)
}

func TestLoggingInAnUser(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	err := store.LogInUser("abcde@gmail.com", "admin123")

	assert.NoError(t, err)
}

func TestCannotLogInAnUserWithANonExistingEmail(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	err := store.LogInUser("hello@gmail.com", "admin")

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.InvalidEmailOrPassword)
}

func TestCannotLogInAnUserWithAnIncorrectPassword(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	err := store.LogInUser("abcde@gmail.com", "aasdjasoidjsd")

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.InvalidEmailOrPassword)
}

/************************/
/***** ORDERS TESTS *****/
/************************/

func TestCreatingAnOrderForValidUser(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}

	err := store.CreateOrder(&newOrder)
	user, err := store.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(user.Orders), "User should have one order")
	assert.Equal(t, newOrder, user.Orders[0])
	assert.Equal(t, uint(1), user.Orders[0].ID, "Order should have an ID")
	assert.Equal(t, 0, len(user.Orders[0].IceCreamTubs), "Order should not have ice cream tubs")
	assert.Equal(t, uint(0), user.Orders[0].DeliveryDriverID, "Order should not have a delivery driver assigned")
	assert.Equal(t, "pending", user.Orders[0].PaymentState, "Order should be in pending state")
	assert.Equal(t, uint(0), user.Orders[0].TotalCost, "Order should not have a total cost")
}

func TestCannotCreateAnOrderForANonExistingUser(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  100,
	}

	err := store.CreateOrder(&newOrder)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.UserIDNotFound)
}

func TestGettingSystemOrderByID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	err := store.CreateOrder(&newOrder)
	order, err := store.GetOrderByID(1)
	assert.NoError(t, err)
	assert.Equal(t, newOrder, order)
}

func TestCannotGetAnOrderByANonExistingID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	_, err := store.GetOrderByID(100)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.OrderNotFound)
}

func TestGettingAllSystemOrders(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	anotherOrder := types.Order{
		Address: "Calle 456",
		UserID:  2,
	}

	err := store.CreateOrder(&newOrder)
	err = store.CreateOrder(&anotherOrder)

	allOrders := store.GetAllOrders()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(allOrders))
	assert.True(t, utils.SliceContains(allOrders, newOrder), "New order should be in the collection")
	assert.True(t, utils.SliceContains(allOrders, anotherOrder), "Another order should be in the collection")
}

func TestGettingAllUserOrdersByEmail(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	anotherOrder := types.Order{
		Address: "Calle 456",
		UserID:  2,
	}

	err := store.CreateOrder(&newOrder)
	err = store.CreateOrder(&anotherOrder)

	userOrders := store.GetAllOrdersByUserEmail("abcde@gmail.com")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(userOrders))
	assert.True(t, utils.SliceContains(userOrders, newOrder), "New order should be in the collection")
	assert.False(t, utils.SliceContains(userOrders, anotherOrder), "Another order should not be in the collection")
}

func TestGettingAnUsersOrderByUserID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}

	err := store.CreateOrder(&newOrder)
	userOrder, err := store.GetUserOrderByID(1, 1)

	assert.NoError(t, err)
	assert.Equal(t, newOrder, userOrder)
}

func TestCannotGetAnUsersOrderByUserIDForANonExistingUser(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}

	err := store.CreateOrder(&newOrder)
	_, err = store.GetUserOrderByID(1, 100)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.UserIDNotFound)
}

func TestCannotGetAnUsersOrderByUserIDForANonExistingOrder(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}

	err := store.CreateOrder(&newOrder)
	_, err = store.GetUserOrderByID(1, 2)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.OrderNotFound)
}

func TestUpdatingAnUsersOrderByID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	//User's IDs is not supposed to change
	updatedOrder := types.Order{
		ID:           1,
		UserID:       1,
		Address:      "Calle 456",
		PaymentState: "paid",
	}

	err := store.CreateOrder(&newOrder)
	actualOrder, err := store.UpdateOrderByID(newOrder.ID, &updatedOrder)

	assert.NoError(t, err)
	assert.Equal(t, updatedOrder.ID, actualOrder.ID)
	assert.Equal(t, updatedOrder.UserID, actualOrder.UserID)
	assert.Equal(t, updatedOrder.TotalCost, actualOrder.TotalCost)
	assert.Equal(t, updatedOrder.DeliveryDriverID, actualOrder.DeliveryDriverID)
	assert.ElementsMatch(t, updatedOrder.IceCreamTubs, actualOrder.IceCreamTubs)
	assert.Equal(t, "Calle 456", actualOrder.Address)
	assert.Equal(t, "paid", actualOrder.PaymentState)
}

func TestCannotUpdateAnUsersOrderByUserIDForANonExistingOrder(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	updatedOrder := types.Order{
		UserID:  1,
		Address: "Calle 456",
	}

	err := store.CreateOrder(&newOrder)
	_, err = store.UpdateOrderByID(100, &updatedOrder)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.OrderNotFound)
}

func TestAddingAnIceCreamTubToAnOrder(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  250,
		Flavors: []string{"ddl", "frt"},
	}

	err := store.CreateOrder(&newOrder)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &newIceCreamTub)
	actualOrder, err := store.GetOrderByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(actualOrder.IceCreamTubs))
	assert.Equal(t, newIceCreamTub, actualOrder.IceCreamTubs[0])
}

func TestCannotAddAnIceCreamTubToAnOrderForANonExistingOrder(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  250,
		Flavors: []string{"ddl", "frt"},
	}

	err := store.CreateOrder(&newOrder)
	err = store.AddIceCreamTubByOrderID(100, &newIceCreamTub)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.OrderNotFound)
}

func TestCannotAddAnIceCreamTubWithNonExistingFlavorsToAnOrder(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  250,
		Flavors: []string{"ddl", "frt", "non-existing-flavor"},
	}
	err := store.CreateOrder(&newOrder)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &newIceCreamTub)
	actualOrder, _ := store.GetOrderByID(1)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.NonExistingFlavors)
	assert.Equal(t, 0, len(actualOrder.IceCreamTubs))
}

func TestCannotAddAnIceCreamTubWithANonAvailableWeightToAnOrder(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  123,
		Flavors: []string{"ddl", "frt"},
	}
	err := store.CreateOrder(&newOrder)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &newIceCreamTub)
	actualOrder, _ := store.GetOrderByID(1)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.WeightNotAvailable)
	assert.Equal(t, 0, len(actualOrder.IceCreamTubs))
}

func TestTotalCostOfAnOrderIsUpdatedCorrectlyWhenAddingIceCreamTubs(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  500,
		Flavors: []string{"ddl", "frt"},
	}
	anotherIceCreamTub := types.IceCreamTub{
		Weight:  250,
		Flavors: []string{"mrc", "frt"},
	}

	err := store.CreateOrder(&newOrder)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &newIceCreamTub)
	actualOrder, _ := store.GetOrderByID(1)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &anotherIceCreamTub)
	actualOrder, _ = store.GetOrderByID(1)

	assert.NoError(t, err)
	assert.Equal(t, prices[500]+prices[250], actualOrder.TotalCost)
}

func TestGettingAllIceCreamTubsByOrderID(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  500,
		Flavors: []string{"ddl", "frt"},
	}
	anotherIceCreamTub := types.IceCreamTub{
		Weight:  250,
		Flavors: []string{"mrc"},
	}

	err := store.CreateOrder(&newOrder)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &newIceCreamTub)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &anotherIceCreamTub)
	allIceCreamTubs, err := store.GetIceCreamTubsByOrderID(newOrder.ID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(allIceCreamTubs))
	assert.True(t, utils.SliceContains(allIceCreamTubs, newIceCreamTub))
	assert.True(t, utils.SliceContains(allIceCreamTubs, anotherIceCreamTub))

}

func TestCannotGetAllIceCreamTubsByOrderIDForANonExistingOrder(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  500,
		Flavors: []string{"ddl", "frt"},
	}

	err := store.CreateOrder(&newOrder)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &newIceCreamTub)
	allIceCreamTubs, err := store.GetIceCreamTubsByOrderID(100)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.OrderNotFound)
	assert.Equal(t, 0, len(allIceCreamTubs))
}

func TestDeletingAnIceCreamTubFromAnOrder(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  500,
		Flavors: []string{"ddl", "frt"},
	}

	err := store.CreateOrder(&newOrder)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &newIceCreamTub)
	err = store.DeleteIceCreamTubByOrderID(newOrder.ID, 1)
	actualOrder, _ := store.GetOrderByID(1)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(actualOrder.IceCreamTubs))
}

func TestCannotDeleteAnIceCreamTubFromAnOrderForANonExistingOrder(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  500,
		Flavors: []string{"ddl", "frt"},
	}
	err := store.CreateOrder(&newOrder)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &newIceCreamTub)
	err = store.DeleteIceCreamTubByOrderID(1, 100)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.OrderNotFound)
}

func TestCannotDeleteAnIceCreamTubFromAnOrderForANonExistingTubID(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  500,
		Flavors: []string{"ddl", "frt"},
	}
	err := store.CreateOrder(&newOrder)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &newIceCreamTub)
	err = store.DeleteIceCreamTubByOrderID(100, 1)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.IceCreamTubNotFound)
}

func TestCannotDeleteAnIceCreamTubWithMismatchedOrderID(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	anotherOrder := types.Order{
		Address: "Calle 333",
		UserID:  2,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  500,
		Flavors: []string{"ddl", "frt"},
	}
	err := store.CreateOrder(&newOrder)
	err = store.CreateOrder(&anotherOrder)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &newIceCreamTub)
	err = store.DeleteIceCreamTubByOrderID(1, 2)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.IceCreamTubNotFound)
}

func TestTotalCostOfAnOrderIsUpdatedCorrectlyWhenDeletingIceCreamTubs(t *testing.T) {
	store := newStorage(flavors, users, prices)
	defer clearAndCloseConnection(t, store)
	newOrder := types.Order{
		Address: "Calle 123",
		UserID:  1,
	}
	newIceCreamTub := types.IceCreamTub{
		Weight:  500,
		Flavors: []string{"ddl", "frt"},
	}
	anotherIceCreamTub := types.IceCreamTub{
		Weight:  250,
		Flavors: []string{"mrc"},
	}

	err := store.CreateOrder(&newOrder)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &newIceCreamTub)
	err = store.AddIceCreamTubByOrderID(newOrder.ID, &anotherIceCreamTub)
	err = store.DeleteIceCreamTubByOrderID(newOrder.ID, 1)
	actualOrder, _ := store.GetOrderByID(1)

	assert.NoError(t, err)
	assert.Equal(t, prices[250], actualOrder.TotalCost)
}

func TestThereAreNotDeliveryDriversWhenJustInitializedTheStore(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	allDeliveryDrivers := store.GetDeliveryDrivers()

	assert.Equal(t, 0, len(allDeliveryDrivers))
}

/**********************************/
/***** DELIVERY-DRIVERS TESTS *****/
/**********************************/

func TestAddingNewDeliveryDriver(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	allDeliveryDrivers := store.GetDeliveryDrivers()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(allDeliveryDrivers))
	assert.Equal(t, newDeliveryDriver.UserID, allDeliveryDrivers[0].UserID)
	assert.Equal(t, newDeliveryDriver.Cuil, allDeliveryDrivers[0].Cuil)
	assert.Equal(t, newDeliveryDriver.Age, allDeliveryDrivers[0].Age)
	assert.ElementsMatch(t, newDeliveryDriver.Vehicles, allDeliveryDrivers[0].Vehicles)
}

func TestCannotAddANewDeliveryDriverForANonExistingUser(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   100,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	allDeliveryDrivers := store.GetDeliveryDrivers()

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.UserIDNotFound)
	assert.Equal(t, 0, len(allDeliveryDrivers))
}

func TestCannotAddANewDeliveryDriverWhenUserIsAlreadyOne(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	allDeliveryDrivers := store.GetDeliveryDrivers()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(allDeliveryDrivers))

	err = store.AddDeliveryDriver(&newDeliveryDriver)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.UserIsAlreadyADriver)
	assert.Equal(t, 1, len(allDeliveryDrivers))
}

func TestGettingDeliverDriverByID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	delivery, err := store.GetDeliveryDriverByID(1)

	assert.NoError(t, err)
	assert.Equal(t, newDeliveryDriver.UserID, delivery.UserID)
	assert.Equal(t, newDeliveryDriver.Cuil, delivery.Cuil)
	assert.Equal(t, newDeliveryDriver.Age, delivery.Age)
	assert.Equal(t, newDeliveryDriver.Vehicles, delivery.Vehicles)
}

func TestCannotGetDeliverDriverByIDForANonExistingID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	_, err = store.GetDeliveryDriverByID(100)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.DeliveryDriverNotFound)
}

func TestUpdatingDeliveryDriverByID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}
	updatedDeliveryDriver := types.DeliveryDriver{
		Cuil:     "20454545456",
		Age:      21,
		Vehicles: []string{"ABC123", "DEF456", "GHI789"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	err = store.UpdateDeliveryDriverByID(1, &updatedDeliveryDriver)
	delivery, _ := store.GetDeliveryDriverByID(1)

	assert.NoError(t, err)
	assert.Equal(t, updatedDeliveryDriver.UserID, delivery.UserID)
	assert.Equal(t, updatedDeliveryDriver.Cuil, delivery.Cuil)
	assert.Equal(t, updatedDeliveryDriver.Age, delivery.Age)
	assert.ElementsMatch(t, updatedDeliveryDriver.Vehicles, delivery.Vehicles)
}

func TestCannotUpdateDeliveryDriverByIDForANonExistingID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}
	updatedDeliveryDriver := types.DeliveryDriver{
		Cuil:     "20454545456",
		Age:      21,
		Vehicles: []string{"ABC123", "DEF456", "GHI789"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	err = store.UpdateDeliveryDriverByID(100, &updatedDeliveryDriver)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.DeliveryDriverNotFound)

	delivery, _ := store.GetDeliveryDriverByID(1)
	assert.Equal(t, newDeliveryDriver.UserID, delivery.UserID)
	assert.Equal(t, newDeliveryDriver.Cuil, delivery.Cuil)
	assert.Equal(t, newDeliveryDriver.Age, delivery.Age)
	assert.ElementsMatch(t, newDeliveryDriver.Vehicles, delivery.Vehicles)
}

func TestDeletingDeliveryDriverByID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	err = store.DeleteDeliveryDriverByID(1)
	allDeliveryDrivers := store.GetDeliveryDrivers()

	assert.NoError(t, err)
	assert.Equal(t, 0, len(allDeliveryDrivers))
}

func TestCannotDeleteDeliveryDriverByIDForANonExistingID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	err = store.DeleteDeliveryDriverByID(100)
	allDeliveryDrivers := store.GetDeliveryDrivers()

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.DeliveryDriverNotFound)
	assert.Equal(t, 1, len(allDeliveryDrivers))
}

func TestUserIsNotADeliveryDriverAnymoreWhenDeletingDeliveryDriverWithEqualID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	err = store.DeleteDeliveryDriverByID(1)
	user, err := store.GetUserByID(1)

	assert.NoError(t, err)
	assert.False(t, user.IsDeliveryDriver())
}

func TestGettingDeliveryDriverVehiclesByID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	allVehicles, err := store.GetVehiclesByDeliveryDriverID(1)

	assert.NoError(t, err)
	assert.Equal(t, []string{"ABC123", "DEF456"}, allVehicles)
}

func TestCannotGetDeliveryDriverVehiclesByIDForANonExistingID(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	_, err = store.GetVehiclesByDeliveryDriverID(100)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.DeliveryDriverNotFound)
}

func TestAssigningDeliveryDriverToUserOrder(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}
	newOrder := types.Order{
		UserID:  2,
		Address: "Calle 123",
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	err = store.CreateOrder(&newOrder)
	err = store.AssignDeliveryDriverToOrder(1, 1)
	actualOrder, _ := store.GetOrderByID(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), actualOrder.DeliveryDriverID)
}

func TestCannotAssignDeliveryDriverToUserOrderForANonExistingOrder(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}
	newOrder := types.Order{
		UserID:  2,
		Address: "Calle 123",
	}
	err := store.AddDeliveryDriver(&newDeliveryDriver)
	err = store.CreateOrder(&newOrder)
	err = store.AssignDeliveryDriverToOrder(100, 1)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.OrderNotFound)
}

func TestCannotAssignDeliveryDriverToUserOrderForANonExistingDeliveryDriver(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}
	newOrder := types.Order{
		UserID:  2,
		Address: "Calle 123",
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	err = store.CreateOrder(&newOrder)
	err = store.AssignDeliveryDriverToOrder(1, 100)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.DeliveryDriverNotFound)
}

func TestDeleteDeliveryDriverFromUserOrder(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}
	newOrder := types.Order{
		UserID:  2,
		Address: "Calle 123",
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	err = store.CreateOrder(&newOrder)
	err = store.AssignDeliveryDriverToOrder(1, 1)
	err = store.DeleteDeliveryDriverFromOrder(1)
	actualOrder, _ := store.GetOrderByID(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(0), actualOrder.DeliveryDriverID)
}

func TestCannotDeleteDeliveryDriverFromUserOrderForANonExistingOrder(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	err := store.DeleteDeliveryDriverFromOrder(100)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.OrderNotFound)
}

func TestGettingDeliveryDriverFromUserOrder(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)
	newDeliveryDriver := types.DeliveryDriver{
		UserID:   1,
		Cuil:     "20454545456",
		Age:      20,
		Vehicles: []string{"ABC123", "DEF456"},
	}
	newOrder := types.Order{
		UserID:  2,
		Address: "Calle 123",
	}

	err := store.AddDeliveryDriver(&newDeliveryDriver)
	err = store.CreateOrder(&newOrder)
	err = store.AssignDeliveryDriverToOrder(1, 1)
	deliverDriverID, _ := store.GetDeliveryDriverFromOrder(1)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), deliverDriverID)
}

func TestCannotGettingDeliveryDriverFromUserOrderForANonExistingOrder(t *testing.T) {
	store := newStorage([]types.Flavor{}, users, map[uint]uint{})
	defer clearAndCloseConnection(t, store)

	_, err := store.GetDeliveryDriverFromOrder(100)

	assert.Error(t, err)
	assert.EqualError(t, err, messageErrors.OrderNotFound)
}
