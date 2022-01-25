package service

import (
	"github.com/kolosek/pkg/model/domain"
	"github.com/kolosek/pkg/repository"
	"github.com/kolosek/pkg/service/log"
)

//UserService interface
type UserService interface {
	FindByID(id int64) (*domain.User, error)
	Save(user *domain.User) error
	FindFirstN(count int) ([]*domain.User, error)
}

//User service structure
type User struct {
	userRepo repository.User
	logger   *log.Logger
}

// NewUser creates a new User service.
func NewUser(userRepo repository.User, logger *log.Logger) UserService {
	return &User{
		userRepo: userRepo,
		logger:   logger,
	}
}

// FindByID returns a user with the specified ID.
func (u *User) FindByID(id int64) (*domain.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Save creates a new user, or updates an existing user in the database
func (u *User) Save(user *domain.User) error {
	existing, err := u.userRepo.FindByEmail(user.Email)
	if err != nil {
		u.logger.WithError(err).Error("Error during user retrieval")
		return err
	}

	// If the user exists, update ID field
	if existing != nil {
		user.ID = existing.ID
		return nil
	}

	//Create user if they dont exist
	if err := u.userRepo.Create(user); err != nil {
		u.logger.WithError(err).Error("Error during user creation")
		return err
	}

	return nil
}

//FindFirstN users
func (u *User) FindFirstN(count int) ([]*domain.User, error) {
	return u.userRepo.FindFirstN(count)
}
