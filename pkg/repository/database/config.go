package database

import "fmt"

// Config holds all configuration values for the DB setup.
type Config struct {
	// Address is in format host:port.
	Address  string
	User     string
	Password string
	Name     string
	LogMode  bool
}

// NewConfig creates a new instance of Config.
func NewConfig(address, user, password, name string) *Config {
	return &Config{
		Address:  address,
		User:     user,
		Password: password,
		Name:     name,
	}
}

// ConnectionURL returns a valid string for gorm.Open.
func (c *Config) ConnectionURL() string {
	//  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// "password=%s dbname=%s sslmode=disable",
	// host, port, user, password, dbname)
	// }
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Address,
		c.Name,
	)
}
