package configs

import (
	"database/sql"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"task_tracker/internal/helpers"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConfig struct {
	host     string
	Port     string
	user     string
	password string
	name     string
	sslMode  string

	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

func LoadPostgres() (*PostgresConfig, error) {
	host := helpers.GetEnv("DB_HOST", "localhost")
	port := helpers.GetEnv("DB_PORT", "5432")
	user := helpers.GetEnv("DB_USER", "")
	password := helpers.GetEnv("DB_PASSWORD", "")
	name := helpers.GetEnv("DB_NAME", "postgres")
	sslMode := helpers.GetEnv("DB_SSLMODE", "disable")
	maxOpenConns := getEnvInt("DB_MAX_OPEN_CONNS", 10)
	maxIdleConns := getEnvInt("DB_MAX_IDLE_CONNS", 5)
	connMaxLifetime := getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute)
	connMaxIdleTime := getEnvDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute)

	return &PostgresConfig{
		host:            host,
		Port:            port,
		user:            user,
		password:        password,
		name:            name,
		sslMode:         sslMode,
		maxOpenConns:    maxOpenConns,
		maxIdleConns:    maxIdleConns,
		connMaxLifetime: connMaxLifetime,
		connMaxIdleTime: connMaxIdleTime,
	}, nil
}

func getEnvInt(name string, fallback int) int {
	raw := helpers.GetEnv(name, "")
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return fallback
	}
	return value
}

func getEnvDuration(name string, fallback time.Duration) time.Duration {
	raw := helpers.GetEnv(name, "")
	if raw == "" {
		return fallback
	}
	value, err := time.ParseDuration(raw)
	if err != nil || value < 0 {
		return fallback
	}
	return value
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

	db.SetMaxOpenConns(cfg.maxOpenConns)
	db.SetMaxIdleConns(cfg.maxIdleConns)
	db.SetConnMaxLifetime(cfg.connMaxLifetime)
	db.SetConnMaxIdleTime(cfg.connMaxIdleTime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
