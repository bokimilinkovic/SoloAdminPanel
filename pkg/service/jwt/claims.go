package jwt

import (
	"errors"
	"time"
)

const (
	iss = "Kolosek SoLo"
)

type Claims struct {
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	Issuer    string `json:"iss"`
	NotBefore int64  `json:"nbf"`
	Subject   uint   `json:"sub"`
}

// Valid is used to validate token claims while parsing the token.
// It implements the jwt.Claims interface.
func (c *Claims) Valid() error {
	now := time.Now().Unix()
	if now > c.ExpiresAt {
		return errors.New("Token has expired")
	}

	// The value is 0 if it was not set
	if (c.NotBefore == 0 || c.NotBefore > now) || (c.IssuedAt == 0 || c.IssuedAt > now) {
		return errors.New("Token is not valid yet")
	}

	if c.Issuer != iss {
		return errors.New("Invalid issuer")
	}

	if c.Subject == 0 {
		return errors.New("missing subject")
	}

	return nil
}
