package persistence

import (
	"context"
	"golang-jwt-example/pkg/config"
	"log"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	userRepo *UserRepository
)

func TestMain(m *testing.M) {
	log.Println("*** start persistence test ***")
	mongo := getDatabase()

	// init repo
	userRepo = NewUserRepository(mongo)

	res := m.Run()
	// after test

	log.Println("*** end persistence test ***")
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
