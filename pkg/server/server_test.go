package server

import (
	"context"
	"golang-jwt-example/pkg/config"
	"golang-jwt-example/pkg/handler"
	"golang-jwt-example/pkg/infrastructure/persistence"
	"golang-jwt-example/pkg/io"
	"golang-jwt-example/pkg/middleware"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

var (
	testServer *httptest.Server
)

func TestMain(m *testing.M) {
	// before test
	// create logger
	var err error
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("failed to setup loggerr: %s\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// load config
	ctx := context.Background()
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Printf("failed to load config: %s\n", err)
		os.Exit(1)
	}

	// init db
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(cfg.DB.URI))
	if err != nil {
		log.Printf("failed to create mongo db client: %s\n", err)
		os.Exit(1)
	}
	mongoCtx, mongoCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoCancel()

	if err := mongoClient.Connect(mongoCtx); err != nil {
		log.Printf("failed to connect to mongo db: %s\n", err)
		os.Exit(1)
	}

	if err := mongoClient.Ping(mongoCtx, readpref.Primary()); err != nil {
		log.Printf("failed to ping mongo db: %s\n", err)
		os.Exit(1)
	}

	mongoDB := mongoClient.Database(cfg.DB.Database)

	repositories, err := persistence.NewRepositories(mongoDB)
	if err != nil {
		log.Printf("failed to connect to mongo db: %s\n", err)
		os.Exit(1)
	}

	// redis
	redisClient := io.NewRedisClient(cfg)

	// start server
	handlerConfig := &handler.Config{
		AccessTokenSecret:          cfg.Auth.AccessTokenSecret,
		AccessTokenExpiredDuration: time.Duration(cfg.Auth.AccessTokenExpiredDuration),
	}
	middlewareConfig := &middleware.Config{
		AccessTokenSecret: cfg.Auth.AccessTokenSecret,
		// RefreshTokenSecret:          cfg.Auth.RefreshTokenSecret,
		AccessTokenExpiredDuration: time.Duration(cfg.Auth.AccessTokenExpiredDuration),
		// RefreshTokenExpiredDuration: time.Duration(cfg.Auth.RefreshTokenExpiredDuration),
	}
	registry := handler.NewHandler(logger, repositories, handlerConfig, redisClient)
	s := NewServer(
		registry,
		middleware.NewMiddleware(logger, repositories, middlewareConfig),
		&Config{Log: logger},
	)
	testServer = httptest.NewServer(s.Mux)
	defer testServer.Close()

	res := m.Run()
	// after test

	os.Exit(res)
}

func getDatabase() *mongo.Database {
	cfg, _ := config.LoadConfig(context.Background())
	opts := &options.ClientOptions{}
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(cfg.DB.URI), opts)
	if err != nil {
		log.Fatal()
	}
	mongoCtx, mongoCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer mongoCancel()
	if err := mongoClient.Connect(mongoCtx); err != nil {
		log.Fatal()
	}

	if err := mongoClient.Ping(mongoCtx, readpref.Primary()); err != nil {
		log.Fatal()
	}

	return mongoClient.Database(cfg.DB.Database)

}
