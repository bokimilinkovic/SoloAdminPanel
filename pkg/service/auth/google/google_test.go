package google_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kolosek/pkg/service/auth/google"
	"github.com/kolosek/pkg/service/auth/google/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestAuth_Authenticate(t *testing.T) {
	assert := assert.New(t)

	t.Run("couldn't not exchange auth_code for token", func(t *testing.T) {
		oidc, auth, ctx := setupAuthMocks()
		oidc.On("Exchange", ctx, "42").Return(nil, errors.New("unexpected error"))
		usr, err := auth.Authenticate("42")
		assert.EqualError(err, "unexpected error")
		assert.IsType(&google.ExchangeError{}, err)
		assert.Nil(usr)
	})

	t.Run("couldn't retrieve user if oauth2 token verification produces an error", func(t *testing.T) {
		oidc, auth, ctx := setupAuthMocks()
		token := &oauth2.Token{}
		oidc.On("Exchange", ctx, "42").Return(token, nil)
		oidc.On("Verify", ctx, token).Return(nil, errors.New("unexpected error"))
		usr, err := auth.Authenticate("42")
		assert.EqualError(err, "unexpected error")
		assert.IsType(&google.VerificationError{}, err)
		assert.Nil(usr)
	})

	t.Run("user is retrieved if exchange and verification steps produce no error", func(t *testing.T) {
		oidc, auth, ctx := setupAuthMocks()
		expectedClaims := &google.Claims{
			Email:      "john.doe@example.com",
			GivenName:  "John Doe",
			PictureURL: "johndoe.jpg",
			Subject:    "42",
		}
		token := &oauth2.Token{}
		oidc.On("Exchange", ctx, "42").Return(token, nil)
		oidc.On("Verify", ctx, token).Return(expectedClaims, nil)
		usr, err := auth.Authenticate("42")
		assert.NoError(err)
		assert.EqualValues(expectedClaims.GivenName, usr.FirstName)
		assert.EqualValues(expectedClaims.Email, usr.Email)
	})
}

func setupAuthMocks() (*mock.OidcAuthenticator, *google.Auth, context.Context) {
	oidc := &mock.OidcAuthenticator{}
	auth := google.NewAuth(oidc)
	ctx := context.Background()
	return oidc, auth, ctx
}
