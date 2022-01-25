package repository

import "github.com/kolosek/pkg/model/domain"

// User holds methods for management of the users in the DB.
type User interface {
	FindByID(id int64) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Exists(user *domain.User) error
	Create(user *domain.User) error
	Update(user *domain.User) error
	FindFirstN(count int) ([]*domain.User, error)
}
