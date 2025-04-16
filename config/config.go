package config

import (
	"time"

	"github.com/spf13/viper"
)

type AuthJWT struct {
	Secret          string
	Algorithm       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type PostgresConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

type HttpConfig struct {
	Host               string
	Port               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxHeaderMegabytes int
}

type Config struct {
	AuthJWT  AuthJWT
	Postgres PostgresConfig
	HTTP     HttpConfig
}

func NewSettings() *Config {
	viper.AutomaticEnv()

	return &Config{
		AuthJWT: AuthJWT{
			Secret:          viper.GetString("JWT_SECRET"),
			Algorithm:       viper.GetString("JWT_ALGORITHM"),
			AccessTokenTTL:  viper.GetDuration("ACCESS_TOKEN_TTL"),
			RefreshTokenTTL: viper.GetDuration("REFRESH_TOKEN_TTL"),
		},
		Postgres: PostgresConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			Name:     viper.GetString("DB_NAME"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			SSLMode:  viper.GetString("DB_SSL"),
		},
		HTTP: HttpConfig{
			Host:               viper.GetString("HTTP_HOST"),
			Port:               viper.GetString("HTTP_PORT"),
			ReadTimeout:        viper.GetDuration("READ_TIMEOUT"),
			WriteTimeout:       viper.GetDuration("WRITE_TIMEOUT"),
			MaxHeaderMegabytes: viper.GetInt("MAX_HEADER_MBYTES"),
		},
	}
}

func (cfg *Config) GetAuthJWTSecret() string {
	return cfg.AuthJWT.Secret
}

func (cfg *Config) GetAccessTokenExpiration() time.Duration {
	return time.Duration(cfg.AuthJWT.AccessTokenTTL)
}

func (cfg *Config) GetRefreshTokenExpiration() time.Duration {
	return time.Duration(cfg.AuthJWT.RefreshTokenTTL)
}
