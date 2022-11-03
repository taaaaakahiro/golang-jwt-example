package command

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang-jwt-example/pkg/config"
	"golang-jwt-example/pkg/handler"
	"golang-jwt-example/pkg/infrastructure/persistence"
	"golang-jwt-example/pkg/io"
	"golang-jwt-example/pkg/middleware"
	"golang-jwt-example/pkg/server"
	"golang-jwt-example/pkg/version"
	"golang.org/x/sync/errgroup"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	exitOk    = 0
	exitError = 1
)

func Run() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	// Logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup logger: %s\n", err)
		return exitError
	}
	defer logger.Sync()
	logger = logger.With(zap.String("version", version.Version))

	// Config
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		logger.Error("failed to load config", zap.Error(err))
		return exitError
	}

	// Listener
	listener, err := net.Listen("tcp", cfg.Address())
	if err != nil {
		logger.Error("failed to listen port", zap.Int("port", cfg.Server.Port), zap.Error(err))
		return exitError
	}
	logger.Info("server start listening", zap.Int("port", cfg.Server.Port))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// DB
	sqlSetting := &config.SQLDBSettings{
		SqlDsn:              cfg.DB.DSN,
		SqlMaxOpenConns:     cfg.DB.MaxOpenConns,
		SqlMaxIdleConns:     cfg.DB.MaxIdleConns,
		SqlConnsMaxLifetime: cfg.DB.ConnsMaxLifetime,
	}
	db, err := io.NewDatabase(sqlSetting)
	if err != nil {
		logger.Error("failed to connect db", zap.Error(err))
		return exitError
	} else {
		logger.Info("successed to connect db")
	}

	// Repository
	repositories, err := persistence.NewRepositories(db)
	if err != nil {
		logger.Error("failed to create repositories", zap.Error(err))
		return exitError
	}

	// Http server
	handlerConfig := &handler.Config{
		AccessTokenSecret:           cfg.Auth.AccessTokenSecret,
		RefreshTokenSecret:          cfg.Auth.RefreshTokenSecret,
		AccessTokenExpiredDuration:  time.Duration(cfg.Auth.AccessTokenExpiredDuration),
		RefreshTokenExpiredDuration: time.Duration(cfg.Auth.RefreshTokenExpiredDuration),
	}
	middlewareConfig := &middleware.Config{
		AccessTokenSecret:           cfg.Auth.AccessTokenSecret,
		RefreshTokenSecret:          cfg.Auth.RefreshTokenSecret,
		AccessTokenExpiredDuration:  time.Duration(cfg.Auth.AccessTokenExpiredDuration),
		RefreshTokenExpiredDuration: time.Duration(cfg.Auth.RefreshTokenExpiredDuration),
	}
	httpServer := server.NewServer(
		handler.NewHandler(logger, repositories, handlerConfig),
		middleware.NewMiddleware(logger, repositories, middlewareConfig),
		&server.Config{Log: logger},
	)
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error {
		return httpServer.Serve(listener)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	logger.Sugar().Infof("SIGNAL %d received, then shutting down...", <-quit)

	if err := httpServer.GracefulShutdown(ctx); err != nil {
		logger.Error("failed to gracefully shutdown", zap.Error(err))
		return exitError
	}

	logger.Info("server shutdown")
	return exitOk
}
