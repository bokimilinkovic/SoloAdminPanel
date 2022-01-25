package mock

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// Generator mocks a `cookie.Generator` object.
type Generator struct {
	mock.Mock
}

// Generate mocks implementation of the real method.
func (g *Generator) HTTPOnly(name, value, path string) *http.Cookie {
	args := g.Called(name, value, path)

	if args.Get(0) != nil {
		return args.Get(0).(*http.Cookie)
	}

	return nil
}
