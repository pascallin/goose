package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	DB      *mongo.Database
	Client  *mongo.Client
	Context context.Context
}

type databaseEnvSetting struct {
	username string
	password string
	endpoint string
	port string
}

type DatabaseOptions struct {
	ConnectTimeout time.Duration `validate:"isdefault=5"` // Timeout operations after N seconds
	URL string `validate:"required_without=UsingEnv"` // connection URL string
	UsingEnv bool `validate:"required_without=URL"` // using env
	DatabaseName string `validate:"required"` // databaseName
	databaseEnvSetting
}

var DB *mongo.Database

func getMongoConnURLFromEnv(ops *DatabaseOptions) (string, error) {
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	clusterEndpoint := os.Getenv("MONGODB_ENDPOINT")
	port := os.Getenv("MONGODB_PORT")
	if username == "" || password == "" || clusterEndpoint == "" || port == "" {
		return "", errors.New("missing env")
	}
	log.WithFields(log.Fields{
		"username": username,
		"password": password,
		"clusterEndpoint": clusterEndpoint,
	})
	mongoConnStringTemplate := "mongodb://%s:%s@%s:%s"
	connectionURI := fmt.Sprintf(mongoConnStringTemplate, username, password, clusterEndpoint, port)
	return connectionURI, nil
}

func validateOptions(ops *DatabaseOptions) error {
	var validate = validator.New()
	err := validate.Struct(ops)
	if err != nil {
		return err
	}
	log.WithField("ops", ops).Info("mongo database options.")
	return nil
}

func NewMongoDatabase(ops *DatabaseOptions) (*Database, error) {
	err := godotenv.Load()

	err = validateOptions(ops)
	if err != nil {
		log.Error("database options not valid.")
		return nil, err
	}
	var connectionURL string
	if ops.URL != "" {
		log.Info("using options URL.")
		connectionURL = ops.URL
	} else {
		log.Info("using env URL.")
		connectionURL, err = getMongoConnURLFromEnv(ops)
		if err != nil {
			return nil, err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), ops.ConnectTimeout * time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURL))
	if err != nil {
		log.WithField("URL", connectionURL).Error("could not parse connection URL.")
		return nil, err
	}
	ctxPING, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	err = client.Ping(ctxPING, readpref.Primary())
	if err != nil {
		return nil, err
	}
	db := client.Database(ops.DatabaseName)
	mongoClient := &Database{DB: db, Client: client, Context: ctx}
	log.Info("mongodb has been connected")
	DB = mongoClient.DB
	return mongoClient, nil
}

func (d *Database) Close() error {
	err := d.Client.Disconnect(d.Context)
	if err != nil {
		log.WithField("message", err.Error()).Error("mongo database close error")
		return err
	}
	return nil
}
