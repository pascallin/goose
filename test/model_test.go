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
	ID    primitive.ObjectID `goose:"primary" bson:"_id,omitempty""`
	Name  string             `goose:"-" bson:"name,omitempty"`
	Email string             `goose:"-" bson:"email,omitempty"`
}

type Post struct {
	ID        primitive.ObjectID `goose:"primary" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `goose:"populate=User" bson:"userId,omitempty" ref:"TestUsers" forignKey:"_id"`
	Title     string             `goose:"-" bson:"title,omitempty"`
	CreatedAt time.Time          `goose:"index" bson:"createdAt,omitempty"`
}

// func TestModel(t *testing.T) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	db, err := goose.NewMongoDatabase(&goose.DatabaseOptions{
// 		UsingEnv: true,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer db.Close()
// 	userID := primitive.NewObjectID()
// 	user := &User{
// 		ID:    userID,
// 		Name:  "John Doe",
// 		Email: "john@example",
// 	}
// 	userModel := goose.NewModel("TestUsers", user)
// 	postModel := goose.NewModel("TestPosts", &Post{
// 		ID:        primitive.NewObjectID(),
// 		UserID:    userID,
// 		Title:     "test post",
// 		CreatedAt: time.Now(),
// 	})
// 	user.Name = "Pascal Lin"
// 	err  userModel.Save()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = postModel.Save()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	result, err := userModel.FindOne(bson.M{"_id": userID})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(result)
// }

func TestPopulate(t *testing.T) {
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
	postModel := goose.NewModel("TestPosts", &Post{})
	result, err := postModel.Populate("User").Find(bson.M{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
