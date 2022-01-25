package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Token wraps a raw JWT and its claims.
type Token struct {
	Raw    string
	claims *jwt.StandardClaims
}

// NewToken initializes a new Token object by parsing the claims of the provided raw JWT.
func NewToken(raw string) (*Token, error) {
	claims, err := parseClaims(raw)
	if err != nil {
		return nil, err
	}

	return &Token{Raw: raw, claims: claims}, nil
}

// IsExpired checks whether the token has expired.
func (t *Token) IsExpired() bool {
	required := true
	return !t.claims.VerifyExpiresAt(time.Now().Unix(), required)
}

func parseClaims(rawToken string) (*jwt.StandardClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(rawToken, &jwt.StandardClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || claims == nil {
		return nil, errors.New("Unable to parse jwt")
	}

	return claims, nil
}
