package main

import "icecreamshop/internal/types"

var flavorDDL types.Flavor = types.Flavor{ID: "ddl", Name: "Dulce de leche", Type: "Dulce de leches"}
var flavorMRC types.Flavor = types.Flavor{ID: "mrc", Name: "Chocolate marroc", Type: "Chocolates"}
var flavorTRM types.Flavor = types.Flavor{ID: "trm", Name: "Tramontana", Type: "Cremas"}
var flavorFRT types.Flavor = types.Flavor{ID: "frt", Name: "Frutilla al agua", Type: "Al agua"}

func initialFlavors() []types.Flavor {
	return []types.Flavor{flavorDDL, flavorMRC, flavorTRM, flavorFRT}
}

func initialUsers() []types.User {
	return []types.User{
		{
			ID:          1,
			Email:       "abcde@gmail.com",
			Name:        "abcde",
			LastName:    "xyz",
			Password:    "$2a$10$xQy8YTOUh6GST9zO1cfmZeV4iPi1I5TLEr5WnTE7Y/XNHgLbqEeFO", //hash for "admin123"
			Orders:      []types.Order{},
			Permissions: []string{"admin"},
		},
	}
}

func initialPrices() map[uint]uint {
	return map[uint]uint{
		250:  3,
		500:  5,
		1000: 10,
	}
}
