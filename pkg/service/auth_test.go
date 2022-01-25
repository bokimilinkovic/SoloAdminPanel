package service_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/kolosek/pkg/model/domain"
	"github.com/kolosek/pkg/model/dto"
	"github.com/kolosek/pkg/service"
	googlemock "github.com/kolosek/pkg/service/auth/google/mock"
	cookiemock "github.com/kolosek/pkg/service/cookie/mock"
	jwtmock "github.com/kolosek/pkg/service/jwt/mock"
	"github.com/kolosek/pkg/service/log"
	"github.com/kolosek/pkg/service/mock"
	"github.com/stretchr/testify/assert"
)

func TestAuth_WithGoogle(t *testing.T) {
	assert := assert.New(t)
	authRequest := &dto.AuthRequest{AuthCode: "42"}
	email := "jane@email.com"
	mockUser := &domain.User{
		ID:        12,
		FirstName: "Jane",
		Email:     email,
		ImageURL:  "picture_url",
	}
	_ = mockUser

	t.Run("couldn't authenticate user if google auth fails", func(t *testing.T) {
		authSvc, _, _, googleAuth, _ := setupAuthService()

		googleAuth.On("Authenticate", "42").Return(nil, errors.New("unexpected error"))

		resp, cookie, err := authSvc.WithGoogle(authRequest)
		assert.EqualError(err, "unexpected error")
		assert.Nil(resp)
		assert.Empty(cookie)
	})

	t.Run("error reading/saving user from userService", func(t *testing.T) {
		authSvc, _, _, googleAuth, userServiceMock := setupAuthService()

		googleAuth.On("Authenticate", "42").Return(mockUser, nil)
		userServiceMock.On("Save", mockUser).Return(errors.New("unexpected error"))
		resp, cookie, err := authSvc.WithGoogle(authRequest)
		assert.EqualError(err, "unexpected error")
		assert.Nil(resp)
		assert.Empty(cookie)
	})

	t.Run("error generatin token", func(t *testing.T) {
		authSvc, cookieSvc, jwtSvc, googleAuth, userServiceMock := setupAuthService()
		_ = cookieSvc
		googleAuth.On("Authenticate", "42").Return(mockUser, nil)
		userServiceMock.On("Save", mockUser).Return(nil)
		jwtSvc.On("Generate", mockUser).Return("", errors.New("error generating token"))
		resp, cookie, err := authSvc.WithGoogle(authRequest)
		assert.EqualError(err, "error generating token")
		assert.Nil(resp)
		assert.Empty(cookie)
	})

	t.Run("authentication successfull", func(t *testing.T) {
		authSvc, cookieSvc, jwtSvc, googleAuth, userServiceMock := setupAuthService()

		jwt := "y2.y3.y4"
		mockCookie := &http.Cookie{Name: "name", Value: jwt}
		expectedResp := &dto.User{
			ID:        12,
			FirstName: "Jane",
			Email:     email,
			Picture:   "picture_url",
		}

		googleAuth.On("Authenticate", "42").Return(mockUser, nil)
		userServiceMock.On("Save", mockUser).Return(nil)
		jwtSvc.On("Generate", mockUser).Return(jwt, nil)
		cookieSvc.On("HTTPOnly", "auth", jwt, "/").Return(mockCookie)

		resp, cookie, err := authSvc.WithGoogle(authRequest)
		assert.NoError(err)
		assert.Equal(expectedResp.Email, resp.Email, expectedResp.FirstName, resp.FirstName)
		assert.Equal(cookie, mockCookie)
	})
}

func TestAuth_Logout(t *testing.T) {
	authSvc, cookieSvc, _, _, _ := setupAuthService()
	expectedCookie := &http.Cookie{Name: "auth", Value: "", Path: "/"}
	cookieSvc.On("HTTPOnly", "auth", "", "/").Return(expectedCookie)
	responseCookie := authSvc.Logout()
	assert.True(t, time.Now().After(responseCookie.Expires))
}

func setupAuthService() (*service.Auth, *cookiemock.Generator, *jwtmock.Generator, *googlemock.Authenticator, *mock.User) {
	cookieGenerator := &cookiemock.Generator{}
	jwtGenerator := &jwtmock.Generator{}
	googleMock := &googlemock.Authenticator{}
	userServiceMock := &mock.User{}
	authService := service.NewAuth(googleMock, jwtGenerator, cookieGenerator, userServiceMock, log.TestLogger())
	return authService, cookieGenerator, jwtGenerator, googleMock, userServiceMock
}
