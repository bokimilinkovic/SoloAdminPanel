package main

import (
	"time"

	"github.com/kolosek/pkg/app"
	"github.com/kolosek/pkg/repository/database"
	"github.com/kolosek/pkg/server"
	"github.com/kolosek/pkg/service/log"
	"github.com/spf13/viper"
)

// @title SoLo Admin
// @version 1.0
// @description This is a solo admin.
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	loadConfig()

	app.Run(&app.Config{
		Logger: log.Config{
			Type:  viper.GetString("logger.type"),
			Level: viper.GetString("logger.level"),
		},
		Server: server.Config{
			Address: viper.GetString("server.address"),
			//Address:      os.Getenv("PORT"),
			ReadTimeout:  viper.GetDuration("server.read_timeout"),
			WriteTimeout: viper.GetDuration("server.write_timeout"),
			Debug:        viper.GetBool("server.debug"),
			CORS: server.CORSConfig{
				AllowCredentials: viper.GetBool("server.cors.allow_credentials"),
				Headers:          viper.GetStringSlice("server.cors.headers"),
				Methods:          viper.GetStringSlice("server.cors.methods"),
				Origins:          viper.GetStringSlice("server.cors.origins"),
			},
		},
		Google: app.GoogleConfig{
			AuthEnabled:  viper.GetBool("google.auth.enabled"),
			Endpoint:     viper.GetString("google.auth.endpoint"),
			ClientID:     viper.GetString("google.auth.client_id"),
			ClientSecret: viper.GetString("google.auth.client_secret"),
			RedirectURI:  viper.GetString("google.auth.redirect_url"),
			GoogleScopes: []string{"openid", "email", "profile"},
		},
		CookieDomain: viper.GetString("cookie_domain"),
		JWTLifetime:  parseEnvDuration(viper.GetString("JWT.lifetime")),
		JWTSecret:    viper.GetString("JWT.secret"),
		DatabaseConfig: database.Config{
			Address:  viper.GetString("database.address"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			Name:     viper.GetString("database.name"),
			LogMode:  true,
		},
	})
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic("cannot read config: " + err.Error())
	}
}

func parseEnvDuration(env string) time.Duration {
	drt, err := time.ParseDuration(env)
	if err != nil {
		return time.Hour
	}

	return drt
}
