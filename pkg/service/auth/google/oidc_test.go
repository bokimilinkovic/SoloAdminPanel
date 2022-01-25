package google_test

import (
	"context"
	"testing"

	"github.com/kolosek/pkg/service/auth/google"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var (
	testCfg = oauth2.Config{
		ClientID:     "test",
		ClientSecret: "test",
		RedirectURL:  "",
	}
)

func TestOidc_NewOidc2(t *testing.T) {
	t.Run("Couldn't initialize oidc client without client ID", func(t *testing.T) {
		cfg := &oauth2.Config{}
		oidc, err := google.NewOidc2(cfg)
		assert.Error(t, err)
		assert.Nil(t, oidc)

	})

	t.Run("oidc initializer", func(t *testing.T) {
		cfg := &oauth2.Config{ClientID: "test", ClientSecret: "test"}
		oidc, err := google.NewOidc2(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, oidc)

	})
}

func TestOidc_Exchange(t *testing.T) {
	t.Run("Cannot exchange authCode for token", func(t *testing.T) {
		oidc, err := google.NewOidc2(&testCfg)
		assert.NoError(t, err)
		token, err := oidc.Exchange(context.Background(), "bad token")
		assert.Error(t, err)
		assert.Nil(t, token)
	})
}
