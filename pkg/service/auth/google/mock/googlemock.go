package mock

import (
	"github.com/kolosek/pkg/model/domain"
	"github.com/stretchr/testify/mock"
)

// Authenticator mocks a `google.Authenticator` object.
type Authenticator struct {
	mock.Mock
}

// Authenticate mocks implementation of the real method.
func (a *Authenticator) Authenticate(authCode string) (*domain.User, error) {
	args := a.Called(authCode)

	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}

	return nil, args.Error(1)
}
