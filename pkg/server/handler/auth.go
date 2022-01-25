package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kolosek/pkg/model/dto"
	"github.com/kolosek/pkg/service"
	"github.com/kolosek/pkg/service/auth/google"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type authService interface {
	WithGoogle(request *dto.AuthRequest) (*dto.User, *http.Cookie, error)
	Logout() *http.Cookie
}

// Auth handler allows users to login using OIDC and SAML, and performs a logout action.
type Auth struct {
	Handler
	authService authService
	oauthConfig *oauth2.Config
}

//NewAuth function integrated oauth2 configuration and authService
func NewAuth(authService authService, oac *oauth2.Config) *Auth {
	return &Auth{
		authService: authService,
		oauthConfig: oac,
	}
}

var (
	oauthState string
)

// Google godoc
// @Summary Login using google oauth. Exchange code for token
// @Description Login using google.
// @Accept json
// @Produce json
// @Success 200 {object}  dto.User
// @Param auth_code body dto.AuthRequest true "Send Auth code"
// @Failure 401 {object} server.HTTPError
// @Failure 401 {object} google.VerificationError
// @Failure 500 {object} server.HTTPError
// @Router /v1/authenticate/google [post]
func (a *Auth) Google(e echo.Context) error {
	req := &dto.AuthRequest{}
	if err := bindAndValidate(req, e); err != nil {
		return err
	}

	resp, cookie, err := a.authService.WithGoogle(req)
	if err != nil {
		if err == service.ErrUnauthorized {
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		}
		switch err.(type) {
		case *google.ExchangeError, *google.VerificationError:
			return echo.NewHTTPError(http.StatusUnauthorized, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	e.SetCookie(cookie)
	return e.JSON(http.StatusOK, resp)
}

// Logout godoc
// @Summary Logout logouts user.
// @Description Cookie is set to expired.
// @Success 204 {string} string  "user has been logged out successfuly"
// @Router /v1/logout [get]
func (a *Auth) Logout(e echo.Context) error {
	cookie := a.authService.Logout()
	e.SetCookie(cookie)
	return e.NoContent(http.StatusNoContent)
}

// GetCode redirects user to google oauth page.
func (a *Auth) GetCode(c echo.Context) error {
	fmt.Println(a.oauthConfig.ClientID + " " + a.oauthConfig.ClientSecret)
	oauthState = generateStateOauthCookie(c.Response().Writer)

	/*
	   AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
	   validate that it matches the the state query parameter on your redirect callback.
	*/
	u := a.oauthConfig.AuthCodeURL(oauthState)
	return c.Redirect(http.StatusTemporaryRedirect, u)
}

// Callback checks if the state returned from google is the same we provided it. It is it
// we are printing code which will be exchaned for access-token.
// THIS IS JUST A TEST HANDLER, DO NOT CALL THIS IN PRODUCTION.
func (a *Auth) Callback(c echo.Context) error {
	if c.FormValue("state") != oauthState {
		return c.Redirect(http.StatusTemporaryRedirect, "/")
	}

	code := c.FormValue("code")
	fmt.Println(code)

	return c.String(http.StatusOK, code)
}

// ExchangeInfos exchanges code for user infromation to verify everything is workingfine.
// DO NOT USE THIS IN PRODUCTION.
func (a *Auth) ExchangeInfos(c echo.Context) error {
	code := c.QueryParam("code")
	res, err := a.getUserDataFromGoogle(code)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, string(res))
}

func (a *Auth) getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := a.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}
