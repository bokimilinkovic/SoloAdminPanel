package middleware

import (
	"errors"
	"net/http"

	jwtlib "github.com/dgrijalva/jwt-go"
	"github.com/kolosek/pkg/model/domain"
	"github.com/kolosek/pkg/service/jwt"
	"github.com/labstack/echo/v4"
)

// Errors that can occur in middleware.
var (
	ErrUserLoad = errors.New("unable to load user")
)

const (
	ctxKeyUser        = "user"
	ctxKeyQueriedUser = "queried_user"
)

type userFinder interface {
	FindByID(id int64) (*domain.User, error)
}

// UserLoader is used for pre-loading a user from the database,
// and setting it in context for future handler use.
type UserLoader struct {
	userFinder  userFinder
	tokenSecret []byte
}

func NewUserLoaders(userFinder userFinder, tokenSecret []byte) *UserLoader {
	return &UserLoader{
		userFinder:  userFinder,
		tokenSecret: tokenSecret,
	}
}

// FromCookie loads the user by using the 'sub' claim found in the JWT, which is obtained
func (l *UserLoader) FromCookie() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Request().Cookie("auth")
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "No cookie present")
			}
			token := cookie.Value
			claims := &jwt.Claims{}

			_, err = jwtlib.ParseWithClaims(token, claims, func(token *jwtlib.Token) (interface{}, error) {
				return l.tokenSecret, nil
			})
			if err != nil {
				if err == jwtlib.ErrSignatureInvalid {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid signature")
				}
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			if err := l.setUser(c, int64(claims.Subject), ctxKeyUser); err != nil {
				return err
			}

			return next(c)
		}
	}
}

func (l *UserLoader) setUser(ctx echo.Context, userID int64, ctxKey string) error {
	user, err := l.userFinder.FindByID(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	ctx.Set(ctxKey, user)
	return nil
}
