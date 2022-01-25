package mock

import (
	"github.com/kolosek/pkg/model/domain"
	"github.com/stretchr/testify/mock"
)

type User struct {
	mock.Mock
}

func (u *User) FindByID(id int64) (*domain.User, error) {
	args := u.Called(id)

	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.User), nil
}

func (u *User) FindByEmail(email string) (*domain.User, error) {
	args := u.Called(email)

	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	if args.Get(0) == nil && args.Get(1) == nil {
		return nil, nil
	}

	return args.Get(0).(*domain.User), nil
}

func (u *User) Exists(user *domain.User) error {
	args := u.Called(user)

	if args.Get(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (u *User) Create(user *domain.User) error {
	args := u.Called(user)

	if args.Get(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (u *User) Update(user *domain.User) error {
	args := u.Called(user)

	if args.Get(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (u *User) FindFirstN(count int) ([]*domain.User, error) {
	args := u.Called(count)

	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.User), nil
}
