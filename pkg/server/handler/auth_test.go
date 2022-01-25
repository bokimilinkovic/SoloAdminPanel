package handler_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/kolosek/pkg/model/dto"
	"github.com/kolosek/pkg/server/handler"
	"github.com/kolosek/pkg/service"
	"github.com/kolosek/pkg/service/auth/google"
	"github.com/kolosek/pkg/service/log"
	"github.com/kolosek/pkg/service/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuth_Google(t *testing.T) {
	assert := assert.New(t)

	var (
		cookie = &http.Cookie{Name: "authcookie", Value: "123"}
		user   = &dto.User{ID: 1, FirstName: "Test"}
	)

	type suite struct {
		auth        *handler.Auth
		authService *mock.Auth
		logger      *log.Logger
	}

	setup := func() *suite {
		logger := log.TestLogger()
		authService := &mock.Auth{}
		auth := handler.NewAuth(authService, nil)

		return &suite{
			auth:        auth,
			authService: authService,
			logger:      logger,
		}
	}

	tests := []struct {
		name        string
		queryParams string
		payload     string
		response    string
		err         interface{}
		mocks       func(suite *suite)
	}{
		{
			name:     "WithGoole returns error, which will be returned as default Interal server error",
			err:      echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error"),
			payload:  `{"auth_code":"123"}`,
			response: "",
			mocks: func(suite *suite) {
				badRequest := &dto.AuthRequest{AuthCode: "123"}
				suite.authService.On("WithGoogle", badRequest).Return(nil, nil, errors.New("Bad code provided"))
			},
		},
		{
			name:    "WithGoogle returns Unauthorized erorr",
			err:     echo.NewHTTPError(http.StatusUnauthorized, service.ErrUnauthorized),
			payload: `{"auth_code":"123"}`,
			mocks: func(suite *suite) {
				badRequest := &dto.AuthRequest{AuthCode: "123"}
				suite.authService.On("WithGoogle", badRequest).Return(nil, nil, service.ErrUnauthorized)
			},
		},
		{
			name:    "WithGoogle returns ExchangeError",
			err:     echo.NewHTTPError(http.StatusUnauthorized, &google.ExchangeError{}),
			payload: `{"auth_code":"123"}`,
			mocks: func(suite *suite) {
				excError := &google.ExchangeError{}
				badRequest := &dto.AuthRequest{AuthCode: "123"}
				suite.authService.On("WithGoogle", badRequest).Return(nil, nil, excError)
			},
		},
		{
			name:    "WithGoogle returns VerificationError",
			err:     echo.NewHTTPError(http.StatusUnauthorized, &google.VerificationError{}),
			payload: `{"auth_code":"123"}`,
			mocks: func(suite *suite) {
				excError := &google.VerificationError{}
				badRequest := &dto.AuthRequest{AuthCode: "123"}
				suite.authService.On("WithGoogle", badRequest).Return(nil, nil, excError)
			},
		},
		{
			name:     "WithGoogle returns response and cookie",
			err:      nil,
			payload:  `{"auth_code":"123"}`,
			response: `{"email":"", "first_name":"Test", "id":1, "last_name":"", "picture":""}`,
			mocks: func(suite *suite) {
				goodRequest := &dto.AuthRequest{AuthCode: "123"}
				suite.authService.On("WithGoogle", goodRequest).Return(user, cookie, nil)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			suite := setup()
			test.mocks(suite)

			ctx, rec := setupEchoServer(t, test.queryParams, test.payload)

			err := suite.auth.Google(ctx)
			if test.err != nil {
				assert.Equal(test.err, err)
				return
			}
			receivedCookie := rec.Result().Cookies()[0]
			assert.Equal(cookie.Name, receivedCookie.Name, cookie.Value, receivedCookie.Value)
			assert.JSONEq(test.response, rec.Body.String())
		})
	}
}

func TestAuth_Logout(t *testing.T) {
	assert := assert.New(t)
	ctx, rec := setupEchoServer(t, "", "")
	authSvc := mock.Auth{}
	authHandler := handler.NewAuth(&authSvc, nil)

	authSvc.On("Logout").Return(&http.Cookie{})

	err := authHandler.Logout(ctx)
	assert.NoError(err)
	assert.Equal(http.StatusNoContent, rec.Code)
}
