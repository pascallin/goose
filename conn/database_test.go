package goose

import (
	"fmt"
	"testing"
)

func TestConnectUsingURL(t *testing.T) {
	mongoConnStringTemplate := "mongodb://@%s:%s"
	connectionURI := fmt.Sprintf(mongoConnStringTemplate, "localhost", "27017")

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
