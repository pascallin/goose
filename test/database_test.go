package test

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"

	"github.com/pascallin/goose"
)

func TestConnectUsingEnv(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Error(err)
	}
	db, err := goose.NewMongoDatabase(&goose.DatabaseOptions{
		UsingEnv: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	t.Log("Mongo connected.")
}

func TestConnectUsingURL(t *testing.T) {
	mongoConnStringTemplate := "mongodb://%s:%s@%s:%s"
	connectionURI := fmt.Sprintf(mongoConnStringTemplate, "root", "example", "localhost", "27017")

	db, err := goose.NewMongoDatabase(&goose.DatabaseOptions{
		DatabaseName: "test",
		URL:          connectionURI,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	t.Log("Mongo connected.")
}
