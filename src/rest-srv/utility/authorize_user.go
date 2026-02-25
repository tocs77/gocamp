package utility

import (
	"errors"
	"slices"
)

type ContextKey string

func AuthorizeUser(userRole string, allowedRoles ...string) (bool, error) {
	if !slices.Contains(allowedRoles, userRole) {
		return false, ErrorHandler(errors.New("unauthorized"), "unauthorized")
	}
	return true, nil
}
