package server

import (
	"strings"
	"time"

	"github.com/kolosek/pkg/service/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/tylerb/graceful.v1"
)

const (
	echoLoggerFormat = `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
		`"method":"${method}","uri":"${uri}","path":"${path}","status":${status},"referer":"${referer}",` +
		`"user_agent":"${user_agent}","latency":${latency},"latency_human":"${latency_human}"` +
		`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n"
)

type REST struct {
	address string
	engine  *echo.Echo
	logger  *log.Logger
}

func New(cfg Config, logger *log.Logger) *REST {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		Skipper: func(ctx echo.Context) bool {
			return strings.HasPrefix(ctx.Request().RequestURI, "/swagger")
		},
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: echoLoggerFormat,
		Output: logger.Out,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// AllowCredentials: cfg.CORS.AllowCredentials,
		AllowHeaders: cfg.CORS.Headers,
		AllowMethods: cfg.CORS.Methods,
		AllowOrigins: []string{"*"},
	}))

	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout

	e.Debug = cfg.Debug
	e.HideBanner = true

	e.Server.Addr = cfg.Address

	server := &REST{
		address: cfg.Address,
		engine:  e,
		logger:  logger,
	}

	return server
}

// SetErrorHandler sets the error handler.
func (r *REST) SetErrorHandler(errorHandler echo.HTTPErrorHandler) {
	r.engine.HTTPErrorHandler = errorHandler
}

// SetValidation sets the validator and binder that validate incoming payload.
func (r *REST) SetValidation(validator echo.Validator, binder echo.Binder) {
	r.engine.Validator = validator
	r.engine.Binder = binder
}

// SetBasicAuth sets basic authentication.
func (r *REST) SetBasicAuth(validator middleware.BasicAuthValidator, skipper middleware.Skipper) {
	r.engine.Use(
		middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
			Validator: validator,
			Skipper:   skipper,
		}),
	)
}

// SetupRoutes setups the server routes.
func (r *REST) SetupRoutes() *echo.Group {
	return r.engine.Group("")
}

// Run runs the REST server
func (r *REST) Run() {
	r.logger.Infof("server running on: %s", r.address)
	r.logger.Fatal(graceful.ListenAndServe(r.engine.Server, 5*time.Second))
}
