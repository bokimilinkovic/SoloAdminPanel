package mock

import (
	"github.com/kolosek/pkg/model/domain"
	"github.com/stretchr/testify/mock"
)

// Generator mocks a `jwt.Generator` object.
type Generator struct {
	mock.Mock
}

// Generate mocks implementation of the real method.
func (g *Generator) Generate(user *domain.User) (string, error) {
	args := g.Called(user)

	if args.Get(0) != nil {
		return args.Get(0).(string), args.Error(1)
	}

	return "", args.Error(1)
}
