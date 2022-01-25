package jwt_test

import (
	"testing"

	"github.com/kolosek/pkg/service/jwt"
	"github.com/stretchr/testify/assert"
)

func TestToken_IsExpired(t *testing.T) {
	assert := assert.New(t)
	validToken := `eyJhbGciOiJSUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZXhwIjozMjgxN` +
		`DgwMjgwMH0.NPj-GRIeeSm_QOpjoVMamk6wHAMBOKfqeLcfBSqkKXyOyUngwVi7_731Ulqi4FmHJ1nWMd_1lJ9wVca` +
		`cnPz6b8yrSOR1nf0uNCPZCdhf6BkPHkv1VTFr2W2_x-oHq4w1wwLGQSTvWiKl5duDtG0GRxwdC8nlGWUMUqb15PO5p` +
		`KH1o7TO8GzflqnmHvyqaHp3RahTRM26z-5LucvBo1wotvb6vIY4rHezwlTHKEzpuOi-djTiDAErc6uo8llahB0oy6C` +
		`0i9ot3cwb5y0CIxYiMty4SKsyN_lQ9gNOU15BoyD74M_yk3NgFeGvB1nGoATsFZn8d1_0UgiW7hTuxbfULw`
	expiredToken := `eyJhbGciOiJSUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiZXhwIjoxMjU` +
		`3ODk0MDAwfQ.IGAplyZFkCUkYmI3mRVDf_sR-tQEX2Umy0uilBOwbRxNu2elOrJ6-sGOGFCiYhwe9ya-spiOHBZvfd` +
		`uDtdkFOvehgz61utM3nXoo92NhQ_Nov3ts2RD-5U0HdL-5b_a-4faqYLmwF_DoJvtROumAeU0fWAH4moW7H7MU2i26` +
		`d3vAhZDC_VIZeLpK0BFvA__X9Wr0tX7J_ZAQJMUcdWNqWIMme3JYg6K-oszoa2t0oFUK8FNa1E3hyM3l8AaCNV8lYW` +
		`qWTKU0AilbhVGPztgJE1Rlcrfp-eFTZ4AljVWclktL-SKJ8TqraLGyHOXWHx_r53YNoO28e4BC6a8zsgvp9w`
	noExpirationToken := `eyJhbGciOiJSUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.ZkMS5y` +
		`B50eoeMVA8kQBi3p3MbMp6RlvlMYdkzs_C4gBEI1Kcf_bITQykMYBmAePcigYmsiQ8fYyiT6x0uIwlFnvR5CwuLZy0` +
		`cuRYZNBSfU9tTnzcSDTIIFD1uK2MllbGhgwNvowLl-YaIH3ndpUMZrtuKA_CZd44BwwwsGkBxx3RMjDN_edgQGkd6Y` +
		`3nATtEMftcrIMYeXxCAcjvOyo0C7WGyvWGyd05ep-QVBxa264BxmaF-E1xwFS_WLQvR3OJXLxPoJQh1GoQ_vA_BEMh` +
		`2Q176Z_64LVjKT89UPPsc60UFgTWXNCKD5CGQ1xLumC1kzzDdaFEhpX25R2XFVNo3w`

	t.Run("returns false for non-expired token", func(t *testing.T) {
		token, err := jwt.NewToken(validToken)
		assert.NoError(err)
		assert.False(token.IsExpired())
	})

	t.Run("returns true for expired token", func(t *testing.T) {
		token, err := jwt.NewToken(expiredToken)
		assert.NoError(err)
		assert.True(token.IsExpired())
	})

	t.Run("returns true for token without expiration date", func(t *testing.T) {
		token, err := jwt.NewToken(noExpirationToken)
		assert.NoError(err)
		assert.True(token.IsExpired())
	})
}
