package postgres

import "fmt"

// Config contains all properties required to connect to a Postgres instance.
type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Sslmode  string `json:"sslmode"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// ConnectionString generates a connection string from the configuration.
func (c Config) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Sslmode,
	)
}
