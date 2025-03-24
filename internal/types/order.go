package types

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"icecreamshop/internal/messageErrors"
)

type IceCreamTubPrice struct {
	Weight uint `json:"weight" gorm:"not null"`
	Price  uint `json:"price" gorm:"not null"`
}
type IceCreamTub struct {
	ID         uint     `json:"id" gorm:"primaryKey; autoIncrement"`
	Weight     uint     `json:"weight" gorm:"not null"`
	Flavors    []string `json:"flavor" gorm:"-"`
	RawFlavors string   `json:"-" gorm:"column:flavor; type:jsonb"`
	OrderID    uint     `json:"order_id" gorm:"not null"`
}

type Order struct {
	ID               uint          `json:"id" gorm:"primaryKey; autoIncrement"`
	Address          string        `json:"address" gorm:"not null"`
	IceCreamTubs     []IceCreamTub `json:"iceCreamTubs" gorm:"foreignKey:OrderID"`
	UserID           uint          `json:"userID" gorm:"not null"`
	DeliveryDriverID uint          `json:"deliveryDriverID"`
	PaymentState     string        `json:"state" gorm:"not null"`
	TotalCost        uint          `json:"totalCost" gorm:"not null"`
}

func (p *IceCreamTub) Validate() error {
	if p.Weight == 0 {
		return errors.New(messageErrors.WeightCannotBeZero)
	}
	if len(p.Flavors) == 0 || len(p.Flavors) > 4 {
		return errors.New(messageErrors.InvalidAmountOfFlavors)
	}
	return nil
}

func (p *Order) Validate() error {
	if p.Address == "" {
		return errors.New(messageErrors.AddressIsRequired)
	}
	return nil
}

func (p IceCreamTub) IsEqualTo(pote IceCreamTub) bool {

	if p.ID != pote.ID {
		return false
	}
	if p.Weight != pote.Weight {
		return false
	}
	if p.OrderID != pote.OrderID {
		return false
	}
	if len(p.Flavors) != len(pote.Flavors) {
		return false
	}
	for i := range p.Flavors {
		if p.Flavors[i] != pote.Flavors[i] {
			return false
		}
	}

	return true
}

func (p Order) IsEqualTo(pedido Order) bool {

	if p.ID != pedido.ID {
		return false
	}
	if p.Address != pedido.Address {
		return false
	}
	if p.UserID != pedido.UserID {
		return false
	}
	if p.DeliveryDriverID != pedido.DeliveryDriverID {
		return false
	}
	if p.PaymentState != pedido.PaymentState {
		return false
	}
	if p.TotalCost != pedido.TotalCost {
		return false
	}

	if len(p.IceCreamTubs) != len(pedido.IceCreamTubs) {
		return false
	}
	for i := range p.IceCreamTubs {
		if !p.IceCreamTubs[i].IsEqualTo(pedido.IceCreamTubs[i]) {
			return false
		}
	}

	return true
}

// BeforeSave is executed when Gorm is about to save new data in the database.
func (p *IceCreamTub) BeforeSave(tx *gorm.DB) (err error) {
	// Serializes Flavors from slice to JSON
	if p.Flavors != nil {
		raw, err := json.Marshal(p.Flavors)
		if err != nil {
			return err
		}
		p.RawFlavors = string(raw)
		print(p.RawFlavors)
	}
	return nil
}

// AfterFind is executed just after Gorm finds data from the database.
func (p *IceCreamTub) AfterFind(tx *gorm.DB) (err error) {

	// Deserializes RawFlavors from JSON to Slice
	if p.RawFlavors != "" {
		err := json.Unmarshal([]byte(p.RawFlavors), &p.Flavors)
		if err != nil {
			return err
		}
		p.RawFlavors = ""
	}

	return nil
}

func (p *IceCreamTub) AfterCreate(tx *gorm.DB) (err error) {
	p.RawFlavors = ""
	return nil
}

func (p *Order) AfterCreate(tx *gorm.DB) (err error) {
	if p.IceCreamTubs == nil {
		p.IceCreamTubs = []IceCreamTub{}
	}
	return nil
}

func (p *Order) AfterFind(tx *gorm.DB) (err error) {
	if p.IceCreamTubs == nil {
		p.IceCreamTubs = []IceCreamTub{}
	}
	return nil
}
