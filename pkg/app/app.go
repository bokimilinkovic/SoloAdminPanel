package app

import (
	"fmt"
	"time"

	_ "github.com/kolosek/api/openapi" // Auto generated code for swagger docs.
	"github.com/kolosek/pkg/repository/database"
	"github.com/kolosek/pkg/repository/database/postgres"
	"github.com/kolosek/pkg/server"
	"github.com/kolosek/pkg/server/handler"
	"github.com/kolosek/pkg/service"
	"github.com/kolosek/pkg/service/auth/google"
	"github.com/kolosek/pkg/service/cookie"
	"github.com/kolosek/pkg/service/jwt"
	"github.com/kolosek/pkg/service/log"
	"github.com/kolosek/pkg/service/validation"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/oauth2"
	o2gl "golang.org/x/oauth2/google"
)

// Run runs the application
// nolint: govet
func Run(cfg *Config) {
	logger, err := log.New(cfg.Logger)
	if err != nil {
		fmt.Println("Error creating logger: " + err.Error())
	}

	restServer := server.New(cfg.Server, logger)

	validator, err := validation.NewValidator()
	if err != nil {
		logger.WithError(err).Fatal("Couldn't create validator")
	}
	restServer.SetValidation(validator, &validation.Binder{})
	restServer.SetErrorHandler(server.ErrorHandler(logger))

	// Setup database
	db, err := database.CreateDBConnection(&cfg.DatabaseConfig, logger)
	if err != nil {
		logger.WithError(err).Fatal("Can not connect to database")
	}

	userStore := postgres.NewUser(db)

	// ******* THIS USE ONLY ON TEST ENVIRONMENT, ON LOCAL DATABASE ********
	// err = migrateTables(userStore)
	// if err != nil {
	// 	fmt.Println("Cannot migrate tables")
	// 	return
	// }

	oauth2 := &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		Endpoint:     o2gl.Endpoint,
		RedirectURL:  cfg.Google.RedirectURI,
		Scopes:       cfg.Google.GoogleScopes,
	}

	// Setup services
	oidc, err := oidc(cfg, oauth2, logger)
	if err != nil {
		logger.WithError(err).Fatal("Couldn't initialize Google OIDC authenticator")
	}

	googleAuth := google.NewAuth(oidc)
	jwtGenerator := jwt.NewGenerator(time.Now, cfg.JWTLifetime, []byte(cfg.JWTSecret))
	cookieGenerator := cookie.NewGenerator(cfg.CookieDomain, cfg.JWTLifetime)
	userService := service.NewUser(userStore, logger)

	authService := service.NewAuth(
		googleAuth,
		jwtGenerator,
		cookieGenerator,
		userService,
		logger,
	)

	// Setup handlers
	statusHandler := handler.NewStatusH()
	authHandler := handler.NewAuth(authService, oauth2)

	routes := restServer.SetupRoutes()

	// Setting up routes
	routes.GET("/status", statusHandler.CheckStatus)

	v1 := routes.Group("v1")

	v1.GET("/swagger/*", echoSwagger.WrapHandler)
	v1.GET("/logout", authHandler.Logout)
	v1.GET("/login", authHandler.GetCode)
	v1.GET("/exchange", authHandler.ExchangeInfos)

	auth := v1.Group("/authenticate")
	auth.POST("/google", authHandler.Google)
	auth.GET("/auth/google/callback", authHandler.Callback)

	restServer.Run()
}

func oidc(c *Config, o2cfg *oauth2.Config, log *log.Logger) (*google.Oidc, error) {
	if !c.Google.AuthEnabled {
		log.Warn("Google authentication is disabled")
		return nil, nil
	}

	return google.NewOidc2(o2cfg)
}

func migrateTables(userStore *postgres.User) error {
	if err := userStore.MigrateTables(); err != nil {
		return err
	}

	return nil
}
