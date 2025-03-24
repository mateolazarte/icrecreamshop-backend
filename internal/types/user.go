package types

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"icecreamshop/internal/messageErrors"
)

type User struct {
	ID             uint     `json:"id" gorm:"primaryKey; autoIncrement"`
	Email          string   `json:"email" gorm:"unique; not null"`
	Name           string   `json:"name" gorm:"not null"`
	LastName       string   `json:"lastName" gorm:"not null"`
	Password       string   `json:"-" gorm:"not null"`
	Orders         []Order  `json:"order" gorm:"foreignKey:UserID"`
	Permissions    []string `json:"permissions" gorm:"-"`
	RawPermissions string   `json:"-" gorm:"column:permissions; type:jsonb; default:'[]'"`
}

type SignUpInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Password string `json:"password"`
}

func (u *User) IsDeliveryDriver() bool {
	for _, rol := range u.Permissions {
		if rol == "delivery" {
			return true
		}
	}
	return false
}

func (u *User) IsAdmin() bool {
	for _, rol := range u.Permissions {
		if rol == "admin" {
			return true
		}
	}
	return false
}

func (u *User) Validate() error {
	if err := u.ValidateUserDataWithoutPassword(); err != nil {
		return err
	}
	if len(u.Password) < 8 {
		return errors.New(messageErrors.PasswordIsTooShort)
	}
	return nil
}

func (u *User) ValidateUserDataWithoutPassword() error {
	if u.Email == "" {
		return errors.New(messageErrors.EmailIsRequired)
	}
	if u.Name == "" {
		return errors.New(messageErrors.FirstNameIsRequired)
	}
	if u.LastName == "" {
		return errors.New(messageErrors.LastNameIsRequired)
	}
	return nil
}

func (u User) IsEqualTo(anotherUser User) bool {
	if u.ID != anotherUser.ID {
		return false
	}
	if u.Email != anotherUser.Email {
		return false
	}
	if u.Name != anotherUser.Name {
		return false
	}
	if u.LastName != anotherUser.LastName {
		return false
	}
	if u.Password != anotherUser.Password {
		return false
	}
	if len(u.Orders) != len(anotherUser.Orders) {
		return false
	}
	if len(u.Permissions) != len(anotherUser.Permissions) {
		return false
	}
	for i := range u.Orders {
		if !u.Orders[i].IsEqualTo(anotherUser.Orders[i]) {
			return false
		}
	}
	for i := range u.Permissions {
		if u.Permissions[i] != anotherUser.Permissions[i] {
			return false
		}
	}
	return true
}

// BeforeSave is executed when Gorm is about to save new data in the database.
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// Serializes Permissions from slice to JSON
	if u.Permissions != nil {
		raw, err := json.Marshal(u.Permissions)
		if err != nil {
			return err
		}
		u.RawPermissions = string(raw)
	}
	return nil
}

// AfterFind is executed just after Gorm finds data from the database.
func (u *User) AfterFind(tx *gorm.DB) (err error) {
	// Deserializes RawPermissions from JSON to Slice
	if u.RawPermissions != "" {
		err := json.Unmarshal([]byte(u.RawPermissions), &u.Permissions)
		if err != nil {
			return err
		}
	}

	// Makes sure that orders is never nil.
	if u.Orders == nil {
		u.Orders = []Order{}
	}
	return nil
}
