package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwtlib "github.com/dgrijalva/jwt-go"
	"github.com/kolosek/pkg/model/domain"
	"github.com/kolosek/pkg/server/middleware"
	"github.com/kolosek/pkg/service/jwt"
	"github.com/kolosek/pkg/service/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserLoader_FromCookie(t *testing.T) {
	assert := assert.New(t)

	next := func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello")
	}

	t.Run("failed to set user if token is nil", func(t *testing.T) {
		ctx, _, userLoader := setupUserLoader("/")
		h := userLoader.FromCookie()(next)
		assert.EqualError(h(ctx), "code=400, message=No cookie present")
		assert.Nil(ctx.Get("user"))
	})

	t.Run("failed to set user if token can't be cast to *jwtlib.Token", func(t *testing.T) {
		ctx, _, userLoader := setupUserLoader("/")
		cookie := &http.Cookie{Name: "auth", Value: "this can be whatever"}
		ctx.Request().AddCookie(cookie)
		h := userLoader.FromCookie()(next)
		assert.EqualError(h(ctx), "code=400, message=token contains an invalid number of segments")
		assert.Nil(ctx.Get("user"))
	})

	t.Run("failed to set user if token claims are nil", func(t *testing.T) {
		ctx, _, userLoader := setupUserLoader("/")
		mockToken := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, nil)
		tokenString, _ := mockToken.SignedString("testic")
		cookie := &http.Cookie{Name: "auth", Value: tokenString}
		ctx.Request().AddCookie(cookie)
		h := userLoader.FromCookie()(next)
		assert.EqualError(h(ctx), "code=400, message=token contains an invalid number of segments")
		assert.Nil(ctx.Get("user"))
	})

	t.Run("failed to set user if userFinder returns an error", func(t *testing.T) {
		ctx, userSvc, userLoader := setupUserLoader("/")
		mockToken := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, &jwt.Claims{Subject: 21, ExpiresAt: time.Now().Add(1000 * time.Second).Unix(), Issuer: "Kolosek SoLo", IssuedAt: time.Now().Unix(), NotBefore: time.Now().Unix()})
		tokenString, err := mockToken.SignedString([]byte("testic"))
		assert.NoError(err)
		cookie := &http.Cookie{Name: "auth", Value: tokenString}
		ctx.Request().AddCookie(cookie)
		userSvc.On("FindByID", int64(21)).Return(nil, errors.New("unexpected"))
		h := userLoader.FromCookie()(next)
		assert.EqualError(h(ctx), "unexpected")
		assert.Nil(ctx.Get("user"))
	})

	t.Run("failed to set user if userFinder returns a nil user", func(t *testing.T) {
		ctx, userSvc, userLoader := setupUserLoader("/")
		mockToken := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, &jwt.Claims{Subject: 21, ExpiresAt: time.Now().Add(1000 * time.Second).Unix(), Issuer: "Kolosek SoLo", IssuedAt: time.Now().Unix(), NotBefore: time.Now().Unix()})
		tokenString, _ := mockToken.SignedString([]byte("testic"))
		cookie := &http.Cookie{Name: "auth", Value: tokenString}
		ctx.Request().AddCookie(cookie)
		userSvc.On("FindByID", int64(21)).Return(nil, nil)
		h := userLoader.FromCookie()(next)
		assert.EqualError(h(ctx), "code=404, message=user not found")
		assert.Nil(ctx.Get("user"))
	})

	t.Run("user is set if no errpr occurred", func(t *testing.T) {
		ctx, userSvc, userLoader := setupUserLoader("/")
		mockuser := &domain.User{ID: 21}
		mockToken := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, &jwt.Claims{Subject: 21, ExpiresAt: time.Now().Add(1000 * time.Second).Unix(), Issuer: "Kolosek SoLo", IssuedAt: time.Now().Unix(), NotBefore: time.Now().Unix()})
		tokenString, _ := mockToken.SignedString([]byte("testic"))
		cookie := &http.Cookie{Name: "auth", Value: tokenString}
		ctx.Request().AddCookie(cookie)
		userSvc.On("FindByID", int64(21)).Return(mockuser, nil)
		h := userLoader.FromCookie()(next)
		assert.NoError(h(ctx))
		assert.Equal(ctx.Get("user"), mockuser)
	})
}

func setupUserLoader(path string) (echo.Context, *mock.User, *middleware.UserLoader) {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	res := httptest.NewRecorder()
	e := echo.New()
	ctx := e.NewContext(req, res)
	userService := &mock.User{}
	userLoader := middleware.NewUserLoaders(userService, []byte("testic"))
	return ctx, userService, userLoader
}
