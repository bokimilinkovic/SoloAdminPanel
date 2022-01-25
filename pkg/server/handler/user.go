package handler

import (
	"github.com/kolosek/pkg/service"
	"github.com/kolosek/pkg/service/log"
)

// User mnages users api calls.
type User struct {
	Handler
	userService service.UserService
}

// NewUser creates User handler.
func NewUser(logger *log.Logger, userService service.UserService) *User {
	return &User{
		Handler:     Handler{Logger: logger},
		userService: userService,
	}
}
