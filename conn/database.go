package goose

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Database is a base struct for goose mongo
type Database struct {
	DataBase *mongo.Database
	Client   *mongo.Client
	Context  context.Context
}

// DatabaseOptions init database options
type DatabaseOptions struct {
	ConnectTimeout time.Duration `validate:"isdefault=5"`               // Timeout operations after N seconds
	URL            string        `validate:"required_without=UsingEnv"` // connection URL string
	UsingEnv       bool          `validate:"required_without=URL"`      // using env
	DatabaseName   string        `validate:"required_without=UsingEnv"` // databaseName
	username       string
	password       string
	endpoint       string
	port           string
}

// DB a mongo database instance after init
var DB *mongo.Database

func validateOptions(ops *DatabaseOptions) error {
	var validate = validator.New()
	err := validate.Struct(ops)
	if err != nil {
		return err
	}
	log.WithField("ops", ops).Info("mongo database options.")
	return nil
}

// NewMongoDatabase new a goose mongo database
func NewMongoDatabase(ops *DatabaseOptions) (*Database, error) {
	err := validateOptions(ops)
	if err != nil {
		log.Error("database options not valid.")
		return nil, err
	}
	var connectionURL string
	var databaseName string
	connectionURL = ops.URL
	databaseName = ops.DatabaseName
	ctx, cancel := context.WithTimeout(context.Background(), ops.ConnectTimeout*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURL))
	if err != nil {
		log.WithField("URL", connectionURL).Error(err)
		return nil, err
	}
	ctxPING, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctxPING, readpref.Primary())
	if err != nil {
		return nil, err
	}
	db := client.Database(databaseName)
	mongoClient := &Database{DataBase: db, Client: client, Context: ctx}
	DB = mongoClient.DataBase
	return mongoClient, nil
}

// Close close database connection
func (d *Database) Close() error {
	err := d.Client.Disconnect(d.Context)
	if err != nil {
		log.WithField("message", err.Error()).Error("mongo database close error")
		return err
	}
	return nil
}
