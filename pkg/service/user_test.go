package service_test

import (
	"errors"
	"testing"

	"github.com/kolosek/pkg/model/domain"
	"github.com/kolosek/pkg/repository/database/postgres/mock"
	"github.com/kolosek/pkg/service"
	"github.com/kolosek/pkg/service/log"
	"github.com/stretchr/testify/assert"
)

var (
	email = "test@gmail.com"
	guid  = "osdk2dja-2kdjdj2"
	user  = &domain.User{
		ID:        2,
		GUID:      guid,
		FirstName: "TestUser",
		Email:     email,
	}
	unexpectedError = errors.New("unexpected error")
	testID          = int64(2)
	users           = []*domain.User{user}
)

func TestUser_FindByID(t *testing.T) {
	t.Run("User repository return an error on FindByID call", func(t *testing.T) {
		userService, userRepo := setupUserServiceTest()
		userRepo.On("FindByID", testID).Return(nil, unexpectedError)
		usr, err := userService.FindByID(testID)
		assert.Equal(t, unexpectedError, err)
		assert.Nil(t, usr)
	})

	t.Run("User successfully found by ID", func(t *testing.T) {
		userService, userRepo := setupUserServiceTest()
		userRepo.On("FindByID", testID).Return(user, nil)
		usr, err := userService.FindByID(testID)
		assert.Equal(t, usr, user)
		assert.Nil(t, err)
	})
}

func TestUser_Save(t *testing.T) {
	t.Run("User repository return an error on FindByEmail call", func(t *testing.T) {
		userService, userRepo := setupUserServiceTest()
		userRepo.On("FindByEmail", email).Return(nil, unexpectedError)
		err := userService.Save(user)
		assert.Equal(t, err, unexpectedError)
	})

	t.Run("Error creating user in database", func(t *testing.T) {
		userService, userRepo := setupUserServiceTest()
		userRepo.On("FindByEmail", email).Return(nil, nil)
		userRepo.On("Create", user).Return(unexpectedError)
		err := userService.Save(user)
		assert.Equal(t, unexpectedError, err)
	})

	t.Run("User successfully created", func(t *testing.T) {
		userService, userRepo := setupUserServiceTest()
		userRepo.On("FindByEmail", email).Return(user, nil)
		userRepo.On("Create", user).Return(nil)
		err := userService.Save(user)
		assert.Nil(t, err)
	})
}

func TestUser_FindFirstN(t *testing.T) {
	userService, userRepoMock := setupUserServiceTest()
	userRepoMock.On("FindFirstN", 1).Return(users, nil)
	usrs, err := userService.FindFirstN(1)
	assert.NoError(t, err)
	assert.Equal(t, users, usrs)
}

func setupUserServiceTest() (service.UserService, *mock.User) {
	userRepoMock := &mock.User{}
	userService := service.NewUser(userRepoMock, log.TestLogger())
	return userService, userRepoMock
}
