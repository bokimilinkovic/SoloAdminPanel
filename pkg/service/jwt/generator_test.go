package jwt_test

import (
	"testing"
	"time"

	"github.com/kolosek/pkg/model/domain"
	"github.com/kolosek/pkg/service/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	duration, err := time.ParseDuration("1h")
	require.NoError(err)

	expectedToken := "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTgwMTgxODcsImlhdCI6MTU1ODAxNDU4NywiaXNzIjoiS29sb3NlayBTb0xvIiwibmJmIjoxNTU4MDE0NTg3LCJzdWIiOjB9.yUOvKA7cAOyREi_64THUlFAsFvezUKvVAsGPwXoXgViiSVerH3QwrBicXKePhAsu"
	user := &domain.User{}

	secret := []byte("123-abc")
	mockNow := func() time.Time { return time.Unix(1558014587, 0) }
	generator := jwt.NewGenerator(mockNow, duration, secret)

	token, err := generator.Generate(user)

	assert.NoError(err)
	assert.Equal(expectedToken, token)
}
