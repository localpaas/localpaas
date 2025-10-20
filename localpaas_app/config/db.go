package config

import (
	"fmt"
	"time"
)

type DB struct {
	Host            string        `yaml:"host" env:"LP_DB_HOST"`
	Port            int           `yaml:"port" env:"LP_DB_PORT"`
	User            string        `yaml:"user" env:"LP_DB_USER"`
	Password        string        `yaml:"password" env:"LP_DB_PASSWORD"`
	DBName          string        `yaml:"db_name" env:"LP_DB_DB_NAME"`
	MaxOpenConns    int           `yaml:"max_open_conns" env:"LP_DB_MAX_OPEN_CONNS" default:"20"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"LP_DB_MAX_IDLE_CONNS" default:"20"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"LP_DB_MAX_LIFETIME" default:"60m"`
	SSLMode         string        `yaml:"ssl_mode" env:"LP_DB_SSL_MODE" default:"require"`
}

func (c *DB) GetDSN() string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.DBName,
		c.SSLMode,
	)
	return dsn
}
