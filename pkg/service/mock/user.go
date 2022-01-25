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
	if args.Get(0) == nil && args.Get(1) == nil {
		return nil, nil
	}

	return args.Get(0).(*domain.User), nil
}

func (u *User) Save(user *domain.User) error {
	args := u.Called(user)

	if args.Get(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (u *User) FindByName(firstname string) ([]*domain.User, error) {
	args := u.Called(firstname)

	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.User), nil
}

func (u *User) FindUsersByEmail(email string) ([]*domain.User, error) {
	args := u.Called(email)

	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.User), nil
}
func (u *User) FindByGUID(guid string) ([]*domain.User, error) {
	args := u.Called(guid)

	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.User), nil
}

func (u *User) FindFirstN(count int) ([]*domain.User, error) {
	args := u.Called(count)

	if args.Get(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.User), nil
}

func (u *User) Search(query []string) ([]*domain.User, error) {

	return nil, nil
}
