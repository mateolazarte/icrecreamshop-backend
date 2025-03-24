package utils

import (
	"errors"
	"fmt"
	"icecreamshop/internal/messageErrors"
	"strconv"
)

type ComparabableType[T any] interface {
	IsEqualTo(T) bool
}

type Validator interface {
	Validate() error
}

func StringToUint(s string) (uint, error) {
	integer, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, errors.New(messageErrors.MustBeAnInteger)
	}
	return uint(integer), nil
}

// SlicesAreEqual checks if two slices contains the same elements.
func SlicesAreEqual[T ComparabableType[T]](s, t []T) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if !s[i].IsEqualTo(t[i]) {
			return false
		}
	}
	return true
}

// SliceContains is true when slice contains the given element.
func SliceContains[T ComparabableType[T]](s []T, v T) bool {
	for _, a := range s {
		if a.IsEqualTo(v) {
			return true
		}
	}
	return false
}

// DeletePermission returns a slice of permissions without the deleted permission.
func DeletePermission(ps []string, permission string) []string {
	for i, p := range ps {
		if p == permission {
			return append(ps[:i], ps[i+1:]...)
		}
	}
	return ps
}

func CreateJsonSingletonString(field, value string) string {
	return fmt.Sprintf(`{"%s":"%s"}`, field, value)
}
