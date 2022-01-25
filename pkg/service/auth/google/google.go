package google

import (
	"context"

	"github.com/kolosek/pkg/model/domain"
	"golang.org/x/oauth2"
)

type oidcAuthenticator interface {
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	Verify(ctx context.Context, oauth2Token *oauth2.Token) (*Claims, error)
}

type Auth struct {
	authenticator oidcAuthenticator
}

func NewAuth(auth oidcAuthenticator) *Auth {
	return &Auth{authenticator: auth}
}

// Authenticate authenticates a user by
// 1. exchanging the auth code with Google for ID token and
// 2. verifying the received ID token.
func (a *Auth) Authenticate(authCode string) (*domain.User, error) {
	ctx := context.Background()
	oauth2Token, err := a.authenticator.Exchange(ctx, authCode)
	if err != nil {
		return nil, &ExchangeError{err}
	}

	claims, err := a.authenticator.Verify(ctx, oauth2Token)
	if err != nil {
		return nil, &VerificationError{err}
	}

	return claims.ToUser(), nil
}
