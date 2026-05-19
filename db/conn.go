package db

import (
	"database/sql"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type databaseConfig struct {
	host     string
	port     int
	user     string
	password string
	name     string
	sslMode  string
}

func ConnectDB() (*sql.DB, error) {
	_ = godotenv.Load()

	config, err := databaseConfigFromEnv()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", config.connectionString())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Println("Connected to " + config.name)

	return db, nil
}

func databaseConfigFromEnv() (databaseConfig, error) {
	port, err := requiredIntEnv("DB_PORT")
	if err != nil {
		return databaseConfig{}, err
	}

	config := databaseConfig{
		port:    port,
		sslMode: envOrDefault("DB_SSLMODE", "disable"),
	}

	if config.host, err = requiredEnv("DB_HOST"); err != nil {
		return databaseConfig{}, err
	}
	if config.user, err = requiredEnv("DB_USER"); err != nil {
		return databaseConfig{}, err
	}
	if config.password, err = requiredEnv("DB_PASSWORD"); err != nil {
		return databaseConfig{}, err
	}
	if config.name, err = requiredEnv("DB_NAME"); err != nil {
		return databaseConfig{}, err
	}

	return config, nil
}

func (c databaseConfig) connectionString() string {
	query := url.Values{}
	query.Set("sslmode", c.sslMode)

	databaseURL := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.user, c.password),
		Host:     net.JoinHostPort(c.host, strconv.Itoa(c.port)),
		Path:     "/" + c.name,
		RawQuery: query.Encode(),
	}

	return databaseURL.String()
}

func requiredEnv(key string) (string, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}

	return value, nil
}

func requiredIntEnv(key string) (int, error) {
	value, err := requiredEnv(key)
	if err != nil {
		return 0, err
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer: %w", key, err)
	}

	return number, nil
}

func envOrDefault(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}
