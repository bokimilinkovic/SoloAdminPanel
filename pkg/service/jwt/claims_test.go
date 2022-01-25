package jwt_test

import (
	"testing"
	"time"

	"github.com/kolosek/pkg/service/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaims_Valid(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	hour, err := time.ParseDuration("1h")
	require.NoError(err)
	sub := uint(12)
	validClaims := jwt.Claims{
		ExpiresAt: time.Now().Add(hour).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "Kolosek SoLo",
		NotBefore: time.Now().Unix(),
		Subject:   sub,
	}

	t.Run("valid claims", func(t *testing.T) {
		err := validClaims.Valid()
		assert.NoError(err)
	})

	t.Run("expired claims", func(t *testing.T) {
		expiredClaims := validClaims
		expiredClaims.ExpiresAt = time.Now().Add(-hour).Unix()
		err := expiredClaims.Valid()
		assert.EqualError(err, "Token has expired")
	})

	t.Run("claims missing expiry date", func(t *testing.T) {
		expiredClaims := validClaims
		expiredClaims.ExpiresAt = 0
		err := expiredClaims.Valid()
		assert.EqualError(err, "Token has expired")
	})

	t.Run("claims issued in the future", func(t *testing.T) {
		futureClaims := validClaims
		futureClaims.IssuedAt = time.Now().Add(hour).Unix()
		err := futureClaims.Valid()
		assert.EqualError(err, "Token is not valid yet")
	})

	t.Run("claims missing issued date", func(t *testing.T) {
		futureClaims := validClaims
		futureClaims.IssuedAt = 0
		err := futureClaims.Valid()
		assert.EqualError(err, "Token is not valid yet")
	})

	t.Run("claims with not before date in the future", func(t *testing.T) {
		futureClaims := validClaims
		futureClaims.NotBefore = time.Now().Add(hour).Unix()
		err := futureClaims.Valid()
		assert.EqualError(err, "Token is not valid yet")
	})

	t.Run("claims missing not before date", func(t *testing.T) {
		futureClaims := validClaims
		futureClaims.NotBefore = 0
		err := futureClaims.Valid()
		assert.EqualError(err, "Token is not valid yet")
	})

	t.Run("claims without subject", func(t *testing.T) {
		subjectlessClaims := validClaims
		subjectlessClaims.Subject = 0
		err := subjectlessClaims.Valid()
		assert.EqualError(err, "missing subject")
	})
}
