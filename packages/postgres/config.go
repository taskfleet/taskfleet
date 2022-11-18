package postgres

import "fmt"

// ConnectionConfig contains all properties required to connect to a Postgres instance. This type
// can be used with "kelseyhightower/envconfig". If it is not used with the library, all options
// must be set.
type ConnectionConfig struct {
	Host     string `required:"true"`
	Port     int    `default:"5432"`
	Sslmode  string `default:"require"`
	Database string `required:"true"`
	Username string `required:"true"`
	Password string `required:"true"`
}

// String generates a connection string from the configuration.
func (c ConnectionConfig) String() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Sslmode,
	)
}

func (c ConnectionConfig) sqlxString() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Database, c.Username, c.Password, c.Sslmode,
	)
}
