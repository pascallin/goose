package goose

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"
)

func TestConnectUsingEnv(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Error(err)
	}
	db, err := NewMongoDatabase(&DatabaseOptions{
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

	db, err := NewMongoDatabase(&DatabaseOptions{
		DatabaseName: "test",
		URL:          connectionURI,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	t.Log("Mongo connected.")
}
