package types

import (
	"errors"
	"icecreamshop/internal/messageErrors"
)

const (
	dulceDeLeches string = "Dulce de leches"
	chocolates    string = "Chocolates"
	cremas        string = "Cremas"
)

type Flavor struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
	Type string `json:"type" gorm:"not null"`
}

func (f *Flavor) Validate() error {
	if f.ID == "" {

		return errors.New(messageErrors.FlavorIdIsRequired)
	}
	if f.Name == "" {
		return errors.New(messageErrors.FlavorNameIsRequired)
	}
	if f.Type == "" {
		return errors.New(messageErrors.FlavorTypeIsRequired)
	}
	return nil
}

func (f Flavor) IsEqualTo(flavor Flavor) bool {
	if f.ID != flavor.ID {
		return false
	}
	if f.Name != flavor.Name {
		return false
	}
	if f.Type != flavor.Type {
		return false
	}
	return true
}
