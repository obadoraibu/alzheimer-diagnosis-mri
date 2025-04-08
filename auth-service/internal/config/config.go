package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DatabaseConfig *DatabaseConfig
	HttpConfig     *HttpConfig `yaml:"http"`
	AuthConfig     *AuthConfig `yaml:"auth"`
	SmtpConfig     *SmtpConfig `yaml:"smtp"`
}

type DatabaseConfig struct {
	UserRepositoryConfig  *UserRepositoryConfig  `yaml:"user-db"`
	TokenRepositoryConfig *TokenRepositoryConfig `yaml:"token-db"`
}

type HttpConfig struct {
	Port string `yaml:"port" env:"PORT" env-default:"8080"`
}

type AuthConfig struct {
	SigningKey      string `yaml:"signing-key" env:"SIGNING_KEY"`
	AccessTokenTTL  string `yaml:"accessTokenTTL" env:"ACCESS_TOKEN_TTL" env-default:"15m"`
	RefreshTokenTTL string `yaml:"refreshTokenTTL" env:"REFRESH_TOKEN_TTL" env-default:"1440h"`
}

type SmtpConfig struct {
	Host     string `yaml:"host" env:"SMTP_HOST" env-default:"smtp.gmail.com"`
	Port     int    `yaml:"port" env:"SMTP_PORT" env-default:"587"`
	From     string `yaml:"from" env:"SMTP_FROM" env-default:"zes-amur@mail.ru"`
	Password string `yaml:"password" env:"SMTP_PASSWORD" env-default:"SuUh72xWTR6J3zu3zgi9"`
}

type UserRepositoryConfig struct {
	Port     string `yaml:"port" env:"USER_DB_PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"USER_DB_HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"USER_DB_NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"USER_DB_USER" env-default:"user"`
	Password string `env:"USER_DB_PASSWORD"`
}

type TokenRepositoryConfig struct {
	Port     string `yaml:"port" env:"TOKEN_DB_PORT" env-default:"6379"`
	Host     string `yaml:"host" env:"TOKEN_DB_HOST" env-default:"localhost"`
	Password string `yaml:"password" env:"TOKEN_DB_PASSWORD"`
}

func NewConfig(mainConfigPath, dbConfigPath string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(mainConfigPath, &cfg)
	if err != nil {
		logrus.Error("cannot read the config")
		return nil, err
	}

	var userCfg UserRepositoryConfig
	var tokenCfg TokenRepositoryConfig

	err = cleanenv.ReadConfig(dbConfigPath, &userCfg)
	if err != nil {
		logrus.Error("cannot read the config")
		return nil, err
	}

	err = cleanenv.ReadConfig(dbConfigPath, &tokenCfg)
	if err != nil {
		logrus.Error("cannot read the config")
		return nil, err
	}

	cfg.DatabaseConfig = &DatabaseConfig{
		UserRepositoryConfig:  &userCfg,
		TokenRepositoryConfig: &tokenCfg,
	}

	return &cfg, nil
}
