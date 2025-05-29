package app

import (
	"context"
	"medods-test-task/config"
	"medods-test-task/internal/repository"
	"medods-test-task/internal/server"
	"medods-test-task/internal/service"
	"medods-test-task/internal/transport/http"
	"medods-test-task/internal/transport/http/routes"
	"medods-test-task/pkg/email/smtp"
	"medods-test-task/pkg/logger"
	"medods-test-task/pkg/migrator"
	"medods-test-task/pkg/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"medods-test-task/internal/db/postgres"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	serviceName     = "medods_auth_service"
	shutdownTimeout = 5 * time.Second
)

func Run() {
	logs, err := logger.New(serviceName)
	if err != nil {
		panic(err)
	}

	ctx := logger.SetToCtx(context.Background(), logs)

	err = godotenv.Load(".env")
	if err != nil {
		logs.Fatal(ctx, "Error loading .env file", zap.Error(err))
	}

	cfg := config.NewSettings()
	if err != nil {
		logs.Fatal(ctx, "Error config load", zap.Error(err))
	}

	db := postgres.New(ctx, &cfg.Postgres)

	err = migrator.Start(cfg)
	if err != nil {
		logs.Fatal(ctx, "failed to migrate", zap.Error(err))
	}

	sender, err := smtp.NewSMTPSender(cfg.SMTP.Mail, cfg.SMTP.Password, cfg.SMTP.Host, cfg.SMTP.Domain, cfg.SMTP.Port)
	if err != nil {
		logs.Fatal(ctx, "failed to create smtp sender", zap.Error(err))
	}

	tokenMananger := utils.NewManager(cfg)
	authRepo := repository.NewAuthRepo(db)
	emailService := service.NewEmailService(sender, logs, &cfg.Email)
	service := service.NewAuthService(authRepo, tokenMananger, emailService)

	handler := http.NewAppController(service, logs)

	app := gin.New()

	routes.RegistrationRoutes(app, tokenMananger, handler)

	// HTTP server
	srv := server.NewServer(cfg, app)

	go func() {
		if err := srv.Run(ctx); err != nil {
			logs.Error(ctx, "error occurred while running http server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	ctx, shutdown := context.WithTimeout(ctx, shutdownTimeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logs.Error(ctx, "failed shutting down the server", zap.Error(err))
	}

	if err := db.Close(); err != nil {
		logs.Error(ctx, "failed to close database connection", zap.Error(err))
	}

	logs.Info(ctx, "Server gracefully stopped")

	if err := logs.Stop(); err != nil {
		logs.Error(ctx, "failed to sync logger", zap.Error(err))
	}
}
