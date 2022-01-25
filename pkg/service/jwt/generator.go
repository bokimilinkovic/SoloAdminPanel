package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kolosek/pkg/model/domain"
)

// Generator is responsible for JWT generation
type Generator struct {
	jwtLifetime time.Duration
	now         func() time.Time
	secret      []byte
}

// Generator generates a new Generator object
func NewGenerator(now func() time.Time, jwtLifetime time.Duration, secret []byte) *Generator {
	return &Generator{
		jwtLifetime: jwtLifetime,
		now:         now,
		secret:      secret,
	}
}

// Generate generates a JWT from the given parameters.
func (g *Generator) Generate(user *domain.User) (string, error) {
	now := g.now()
	nowUnix := now.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, &Claims{
		ExpiresAt: now.Add(g.jwtLifetime).Unix(),
		IssuedAt:  nowUnix,
		Issuer:    iss,
		NotBefore: nowUnix,
		Subject:   uint(user.ID),
	})

	return token.SignedString(g.secret)
}
