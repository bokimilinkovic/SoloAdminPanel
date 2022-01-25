package mock

import (
	"context"

	"github.com/kolosek/pkg/service/auth/google"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type OidcAuthenticator struct {
	mock.Mock
}

// Exchange mocks implementation of the real method.
func (a *OidcAuthenticator) Exchange(ctx context.Context, authCode string) (*oauth2.Token, error) {
	args := a.Called(ctx, authCode)

	if args.Get(0) != nil {
		return args.Get(0).(*oauth2.Token), args.Error(1)
	}

	return nil, args.Error(1)
}

// Verify mocks implementation of the real method.
func (a *OidcAuthenticator) Verify(ctx context.Context, oauth2Token *oauth2.Token) (*google.Claims, error) {
	args := a.Called(ctx, oauth2Token)

	if args.Get(0) != nil {
		return args.Get(0).(*google.Claims), args.Error(1)
	}

	return nil, args.Error(1)
}
