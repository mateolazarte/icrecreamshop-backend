package types

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"icecreamshop/internal/messageErrors"
)

type DeliveryDriver struct {
	UserID      uint     `json:"userid" gorm:"primaryKey"`
	Cuil        string   `json:"cuil" gorm:"not null;unique"`
	Age         uint     `json:"age" gorm:"not null; check: Age>17"` //check
	Vehicles    []string `json:"vehicles" gorm:"-"`
	RawVehicles string   `json:"-" gorm:"column:vehicles; type:jsonb; not null"`
}

func (d *DeliveryDriver) Validate() error {
	if len(d.Cuil) > 11 || len(d.Cuil) < 10 {
		return errors.New(messageErrors.InvalidCuilFormat)
	}
	if d.Age < 18 {
		return errors.New(messageErrors.AgeMustBeGreaterThan18)
	}
	if len(d.Vehicles) == 0 {
		return errors.New(messageErrors.AtLeastOneVehicleIsRequired)
	}
	for _, vehiculo := range d.Vehicles {
		if len(vehiculo) > 7 || len(vehiculo) < 6 {
			return errors.New(messageErrors.InvalidVehicleIDFormat)
		}
	}
	return nil
}

func (d DeliveryDriver) IsEqualTo(repartidor DeliveryDriver) bool {
	if d.UserID != repartidor.UserID {
		return false
	}
	if d.Cuil != repartidor.Cuil {
		return false
	}
	if d.Age != repartidor.Age {
		return false
	}
	if len(d.Vehicles) != len(repartidor.Vehicles) {
		return false
	}
	for i := range d.Vehicles {
		if d.Vehicles[i] != repartidor.Vehicles[i] {
			return false
		}
	}
	return true
}

// BeforeSave is executed when Gorm is about to save new data in the database.
func (d *DeliveryDriver) BeforeSave(tx *gorm.DB) (err error) {
	// Serializes Vehicles from slice to JSON
	if d.Vehicles != nil {
		raw, err := json.Marshal(d.Vehicles)
		if err != nil {
			return err
		}
		d.RawVehicles = string(raw)
	}
	return nil
}

// AfterFind is executed just after Gorm finds data from the database.
func (d *DeliveryDriver) AfterFind(tx *gorm.DB) (err error) {

	// Deserializes RawVehicles from JSON to Slice
	if d.RawVehicles != "" {
		err := json.Unmarshal([]byte(d.RawVehicles), &d.Vehicles)
		if err != nil {
			return err
		}
	}

	return nil
}
