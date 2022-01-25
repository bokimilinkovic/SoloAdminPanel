package app

import (
	"time"

	"github.com/kolosek/pkg/repository/database"
	"github.com/kolosek/pkg/server"
	"github.com/kolosek/pkg/service/log"
)

//Config structure for project
type Config struct {
	Logger         log.Config
	Server         server.Config
	Google         GoogleConfig
	CookieDomain   string
	JWTLifetime    time.Duration
	JWTSecret      string
	DatabaseConfig database.Config
}

//GoogleConfig structure for OAuth 2
type GoogleConfig struct {
	AuthEnabled  bool
	Endpoint     string
	ClientID     string
	ClientSecret string
	RedirectURI  string
	GoogleScopes []string
}
