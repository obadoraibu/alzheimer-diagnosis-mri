package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/obadoraibu/go-auth/internal/config"
	"github.com/obadoraibu/go-auth/internal/repository"
	"github.com/obadoraibu/go-auth/internal/service"
	"github.com/obadoraibu/go-auth/internal/transport/rest"
	"github.com/obadoraibu/go-auth/internal/transport/rest/handler"
	"github.com/obadoraibu/go-auth/pkg/auth"
	"github.com/obadoraibu/go-auth/pkg/smtp"
	"github.com/sirupsen/logrus"
)

func Run(mainConfigPath, dbConfigPath string) error {
	cfg, err := config.NewConfig(mainConfigPath, dbConfigPath)
	if err != nil {
		logrus.Error("cannot create new config")
		return err
	}
	logrus.Println(cfg.DatabaseConfig.TokenRepositoryConfig.Host)
	logrus.Println(cfg.DatabaseConfig.TokenRepositoryConfig.Port)
	logrus.Println(cfg.DatabaseConfig.TokenRepositoryConfig.Password)

	logrus.Println(cfg.DatabaseConfig.UserRepositoryConfig.Host)
	logrus.Println(cfg.DatabaseConfig.UserRepositoryConfig.Port)
	logrus.Println(cfg.DatabaseConfig.UserRepositoryConfig.Password)
	logrus.Println(cfg.DatabaseConfig.UserRepositoryConfig.User)
	logrus.Println(cfg.DatabaseConfig.UserRepositoryConfig.Name)

	logrus.Println(cfg.HttpConfig.Port)
	logrus.Println(cfg.AuthConfig.SigningKey)

	logrus.Println(cfg.SmtpConfig.Host)
	logrus.Println(cfg.SmtpConfig.Port)
	logrus.Println(cfg.SmtpConfig.Password)
	logrus.Println(os.Getenv("SMTP_PASSWORD"))

	logrus.Println(cfg.DatabaseConfig.MinIORepositoryConfig.Endpoint)
	logrus.Println(cfg.DatabaseConfig.MinIORepositoryConfig.AccessKey)
	logrus.Println(cfg.DatabaseConfig.MinIORepositoryConfig.SecretKey)
	logrus.Println(cfg.DatabaseConfig.MinIORepositoryConfig.Bucket)
	logrus.Println(cfg.DatabaseConfig.MinIORepositoryConfig.UseSSL)

	repo, err := repository.NewRepository(cfg.DatabaseConfig)
	if err != nil {
		logrus.Fatalf("error creating repository, %s", err)
	}

	tokenManager := auth.NewTokenManager(cfg.AuthConfig)

	emailService := smtp.NewEmailSender(cfg.SmtpConfig)

	service := service.NewService(service.Dependencies{
		Repo:         repo,
		TokenManager: tokenManager,
		EmailService: emailService,
	})
	handler := handler.NewHandler(handler.Dependencies{
		Service:      service,
		TokenManager: tokenManager,
	})

	server := rest.NewServer()

	go func() {
		if err := server.Start(cfg.HttpConfig.Port, handler.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Printf("service started on port %s", cfg.HttpConfig.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("service shutting down")

	if err := server.Stop(context.Background()); err != nil {
		logrus.Errorf("error on server shutting down: %s", err.Error())
	}

	if err := repo.Close(); err != nil {
		logrus.Errorf("error on db connection close: %s", err.Error())
	}

	return nil
}
