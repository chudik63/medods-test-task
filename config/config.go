package config

import (
	"time"

	"github.com/spf13/viper"
)

type AuthJWT struct {
	Secret          string
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

type ServerConfig struct {
	MigrationsPath string
}

type EmailConfig struct {
	IPWarningSubject  string
	IPWarningTemplate string
}

type SMTPConfig struct {
	Mail     string
	Host     string
	Port     int
	Password string
	Domain   string
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
	Server   ServerConfig
	SMTP     SMTPConfig
	Email    EmailConfig
}

func NewSettings() *Config {
	viper.AutomaticEnv()

	return &Config{
		AuthJWT: AuthJWT{
			Secret:          viper.GetString("JWT_SECRET"),
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
		Server: ServerConfig{
			MigrationsPath: viper.GetString("MIGRATIONS_PATH"),
		},
		SMTP: SMTPConfig{
			Mail:     viper.GetString("SMTP_MAIL"),
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetInt("SMTP_PORT"),
			Password: viper.GetString("SMTP_PASSWORD"),
			Domain:   viper.GetString("DOMAIN"),
		},
		Email: EmailConfig{
			IPWarningSubject:  viper.GetString("IP_WARNING_SUBJECT"),
			IPWarningTemplate: viper.GetString("IP_WARNING_TEMPLATE"),
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
