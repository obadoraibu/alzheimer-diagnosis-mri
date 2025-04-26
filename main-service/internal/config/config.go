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
	PostgresRepositoryConfig *PostgresRepositoryConfig `yaml:"postgres-db"`
	RedisRepositoryConfig    *RedisRepositoryConfig    `yaml:"redis-db"`
	MinIORepositoryConfig    *MinIOConfig              `yaml:"minio"`
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

type PostgresRepositoryConfig struct {
	Port     string `yaml:"port" env:"POSTGRES_DB_PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"POSTGRES_DB_HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"POSTGRES_DB_NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"POSTGRES_DB_USER" env-default:"user"`
	Password string `env:"POSTGRES_DB_PASSWORD"`
}

type RedisRepositoryConfig struct {
	Port     string `yaml:"port" env:"REDIS_DB_PORT" env-default:"6379"`
	Host     string `yaml:"host" env:"REDIS_DB_HOST" env-default:"localhost"`
	Password string `yaml:"password" env:"REDIS_DB_PASSWORD"`
}

type MinIOConfig struct {
    InternalEndpoint string `yaml:"internal_endpoint" env:"MINIO_INTERNAL_ENDPOINT" env-default:"minio:9000"`
    PublicEndpoint   string `yaml:"public_endpoint"  env:"MINIO_PUBLIC_ENDPOINT"  env-default:"localhost:9000"`
    AccessKey        string `yaml:"access_key"       env:"MINIO_ACCESS_KEY"       env-default:"minioadmin"`
    SecretKey        string `yaml:"secret_key"       env:"MINIO_SECRET_KEY"       env-default:"minioadmin"`
    Bucket           string `yaml:"bucket"           env:"MINIO_BUCKET"           env-default:"mri-scans"`
    UseSSL           bool   `yaml:"use_ssl"          env:"MINIO_SSL"              env-default:"false"`
}

func NewConfig(mainConfigPath, dbConfigPath string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(mainConfigPath, &cfg)
	if err != nil {
		logrus.Error("cannot read the config")
		return nil, err
	}

	var postgresCfg PostgresRepositoryConfig
	var redisCfg RedisRepositoryConfig
	var minIOCfg MinIOConfig

	err = cleanenv.ReadConfig(dbConfigPath, &postgresCfg)
	if err != nil {
		logrus.Error("cannot read the config")
		return nil, err
	}

	err = cleanenv.ReadConfig(dbConfigPath, &redisCfg)
	if err != nil {
		logrus.Error("cannot read the config")
		return nil, err
	}

	err = cleanenv.ReadConfig(dbConfigPath, &minIOCfg)
	if err != nil {
		logrus.Error("cannot read the config")
		return nil, err
	}

	cfg.DatabaseConfig = &DatabaseConfig{
		PostgresRepositoryConfig: &postgresCfg,
		RedisRepositoryConfig:    &redisCfg,
		MinIORepositoryConfig:    &minIOCfg,
	}

	return &cfg, nil
}
