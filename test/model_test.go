package test

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/pascallin/goose"
)

type User struct {
	ID    primitive.ObjectID `goose:"primary" bson:"_id,omitempty" json:"_id,omitempty"`
	Name  string             `goose:"-" bson:"name" json:"name"`
	Email string             `goose:"-" bson:"email" json:"email"`
}

type Post struct {
	ID        primitive.ObjectID `goose:"primary" bson:"_id,omitempty" json:"_id,omitempty"`
	UserID    primitive.ObjectID `goose:"index=1,populate=User" bson:"userId" json:"userId"`
	Title     string             `goose:"-" bson:"title" json:"title"`
	CreatedAt time.Time          `goose:"-" bson:"createdAt" json:"createdAt"`
}

func TestDecode(t *testing.T) {
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

	userID := primitive.NewObjectID()

	user := &User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example",
	}
	userModel := goose.NewModel("TestUsers", user)
	postModel := goose.NewModel("TestPosts", &Post{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Title:     "test post",
		CreatedAt: time.Now(),
	})

	user.Name = "Pascal Lin"

	err = userModel.Save()
	if err != nil {
		t.Fatal(err)
	}

	err = postModel.Save()
	if err != nil {
		t.Fatal(err)
	}

	result, err := userModel.FindOne(bson.M{"_id": userID})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
