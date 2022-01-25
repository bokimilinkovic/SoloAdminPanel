package cookie_test

import (
	"testing"
	"time"

	"github.com/kolosek/pkg/service/cookie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_HTTPOnly(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	domain := ""
	duration, err := time.ParseDuration("1h")
	require.NoError(err)

	generator := cookie.NewGenerator(domain, duration)

	name := "auth"
	value := "y2.y3.y4"
	path := "/path"
	cookie := generator.HTTPOnly(name, value, path)

	assert.Equal(name, cookie.Name)
	assert.Equal(value, cookie.Value)
	assert.Equal(path, cookie.Path)
	assert.True(time.Now().Before(cookie.Expires))
	// Adding time.Second fixes the tests on Windows.
	assert.True(time.Now().Add(duration).Add(time.Second).After(cookie.Expires))
}
