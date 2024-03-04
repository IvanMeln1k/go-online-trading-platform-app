package database

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

var (
	ErrConnPostgres = errors.New("error connect to postgres db")
)

func NewPostgresDB(cfg PostgresConfig) (*sqlx.DB, error) {
	strConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	conn, err := sqlx.Connect("postgres", strConn)
	if err != nil {
		return nil, ErrConnPostgres
	}
	return conn, nil
}
