package configs

import (
	"database/sql"
	"fmt"
	"net/url"
	"task_tracker/internal/helpers"
)

type PostgresConfig struct {
	host     string
	Port     string
	user     string
	password string
	name     string
	sslMode  string
}

func LoadPostgres() (*PostgresConfig, error) {
	host := helpers.GetEnv("DB_HOST", "localhost")
	port := helpers.GetEnv("DB_PORT", "5432")
	user := helpers.GetEnv("DB_USER", "")
	password := helpers.GetEnv("DB_PASSWORD", "")
	name := helpers.GetEnv("DB_NAME", "postgres")
	sslMode := helpers.GetEnv("DB_SSLMODE", "disable")

	return &PostgresConfig{
		host:     host,
		Port:     port,
		user:     user,
		password: password,
		name:     name,
		sslMode:  sslMode,
	}, nil
}
func (c *PostgresConfig) dsn() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.user, c.password),
		Host:   fmt.Sprintf("%s:%s", c.host, c.Port),
		Path:   c.name,
	}

	q := u.Query()
	q.Set("sslmode", c.sslMode)
	u.RawQuery = q.Encode()

	return u.String()
}
func New(cfg PostgresConfig) (*sql.DB, error) {
	var dsn string = cfg.dsn()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
