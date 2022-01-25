package google

import (
	"context"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Oidc struct {
	oauthConfig  *oauth2.Config
	oidcVerifier *oidc.IDTokenVerifier
}

func NewOidc2(cfg *oauth2.Config) (*Oidc, error) {
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, errors.New("Couldn't initialize oidc client without client ID and secret")
	}

	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return nil, errors.New("Couldn't initialize OIDC provider " + err.Error())
	}

	oidcConfig := &oidc.Config{
		ClientID:          cfg.ClientID,
		SkipClientIDCheck: false,
		SkipExpiryCheck:   false,
	}
	v := provider.Verifier(oidcConfig)

	return &Oidc{
		oauthConfig:  cfg,
		oidcVerifier: v,
	}, nil
}

func (o *Oidc) Exchange(ctx context.Context, authCode string) (*oauth2.Token, error) {
	oauth2Token, err := o.oauthConfig.Exchange(ctx, authCode)
	if err != nil {
		return nil, err
	}

	return oauth2Token, nil
}

func (o *Oidc) Verify(ctx context.Context, oauth2Token *oauth2.Token) (*Claims, error) {
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("No id_token field in oatuh2 Token")
	}

	idToken, err := o.oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	var claims *Claims
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	return claims, nil
}
