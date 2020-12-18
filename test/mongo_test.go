package test

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pascallin/goose/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func TestConnectUsingEnv(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Error(err)
	}
	db, err := mongo.NewMongoDatabase(&mongo.DatabaseOptions{
		UsingEnv: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	model := mongo.NewModel("test")

	id, err := model.InsertOne(bson.M{"name": "test"})

	t.Log(id)
}

func TestConnectUsingURL(t *testing.T) {
	mongoConnStringTemplate := "mongodb://%s:%s@%s:%s"
	connectionURI := fmt.Sprintf(mongoConnStringTemplate, "root", "example", "localhost", "27017")

	db, err := mongo.NewMongoDatabase(&mongo.DatabaseOptions{
		DatabaseName: "test",
		URL:          connectionURI,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	model := mongo.NewModel("test")

	id, err := model.InsertOne(bson.M{"name": "test 1"})

	t.Log(id)
}
