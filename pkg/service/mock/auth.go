package mock

import (
	"net/http"

	"github.com/kolosek/pkg/model/dto"
	"github.com/stretchr/testify/mock"
)

type Auth struct {
	mock.Mock
}

func (a *Auth) WithGoogle(request *dto.AuthRequest) (*dto.User, *http.Cookie, error) {
	args := a.Called(request)

	if args.Get(2) != nil {
		return nil, nil, args.Error(2)
	}

	return args.Get(0).(*dto.User), args.Get(1).(*http.Cookie), nil
}

func (a *Auth) Logout() *http.Cookie {
	args := a.Called()
	return args.Get(0).(*http.Cookie)
}
