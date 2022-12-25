package db

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db *mongo.Client
)

var (
	user         = os.Getenv("DB_USER")
	password     = os.Getenv("DB_PASSWORD")
	host         = os.Getenv("DB_HOST")
	port         = os.Getenv("DB_PORT")
	databaseName = os.Getenv("DB_NAME")
)

type MongoRepository struct {
	client *mongo.Client
}

func InitPostgresRepo() {

	dbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, password, host, port, databaseName)
	clientOptions := options.Client().ApplyURI(dbURI)

	var err error
	db, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	err = db.Ping(context.TODO(), nil)

	if err != nil {
		panic(err)
	}
	log.Info("database connected")

	Repo = MongoRepository{
		client:    db,
	}

}

func (p MongoRepository) Todo() error {
	return nil
}