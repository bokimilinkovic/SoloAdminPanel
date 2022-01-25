package service

import (
	"net/http"
	"time"

	"github.com/kolosek/pkg/model/domain"
	"github.com/kolosek/pkg/model/dto"
	"github.com/kolosek/pkg/service/log"
)

const (
	cookieName            = "auth"
	cookiePath            = "/"
	cookieExpiredYearsAgo = 5
)

type googleAuthenticator interface {
	Authenticate(authCode string) (*domain.User, error)
}

type jwtGenerator interface {
	Generate(user *domain.User) (string, error)
}

type cookieGenerator interface {
	HTTPOnly(name, value, path string) *http.Cookie
}

//Auth structure
type Auth struct {
	googleAuth      googleAuthenticator
	jwtGenerator    jwtGenerator
	cookieGenerator cookieGenerator
	userService     UserService
	logger          *log.Logger
}

//NewAuth integrated google authenticatork, jwt generator, cookie generator, user service and logger
func NewAuth(ga googleAuthenticator, jwt jwtGenerator, cg cookieGenerator, userService UserService, l *log.Logger) *Auth {
	return &Auth{
		googleAuth:      ga,
		jwtGenerator:    jwt,
		cookieGenerator: cg,
		userService:     userService,
		logger:          l,
	}
}

// WithGoogle authenticates a user with Google.
// User is also saved to the database - either created or updated, if it already existed.
func (a *Auth) WithGoogle(request *dto.AuthRequest) (*dto.User, *http.Cookie, error) {
	usr, err := a.googleAuth.Authenticate(request.AuthCode)
	if err != nil {
		a.logger.WithError(err).Error("Couldn't authenticate user with google")
		return nil, nil, err
	}

	//Here user should be saved to database if not exists
	if err := a.userService.Save(usr); err != nil {
		return nil, nil, err
	}

	// Generate jwt token
	token, err := a.jwtGenerator.Generate(usr)
	if err != nil {
		a.logger.WithError(err).Error("couldn't generate JWT for user's Google auth")
		return nil, nil, err
	}

	user := &dto.User{
		ID:        uint(usr.ID),
		Email:     usr.Email,
		Picture:   usr.ImageURL,
		FirstName: usr.FirstName,
		Lastname:  usr.LastName,
	}

	cookie := a.cookieGenerator.HTTPOnly(cookieName, token, cookiePath)
	return user, cookie, nil
}

//Logout service method
func (a *Auth) Logout() *http.Cookie {
	cookie := a.cookieGenerator.HTTPOnly(cookieName, "", cookiePath)
	cookie.Expires = time.Now().AddDate(-cookieExpiredYearsAgo, 0, 0)
	return cookie
}
