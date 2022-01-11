package goose

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `goose:"primary" bson:"_id,omitempty""`
	Name        string             `goose:"-" bson:"name,omitempty"`
	Email       string             `goose:"-" bson:"email,omitempty"`
	CreatedTime time.Time          `goose:"createdAt" bson:"createdTime,omitempty"`
}

type Post struct {
	ID          primitive.ObjectID `goose:"primary" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `goose:"populate=User" bson:"userId,omitempty" ref:"TestUsers" forignKey:"_id"`
	Title       string             `goose:"-" bson:"title,omitempty"`
	Description string             `goose:"default='No description.'"  bson:"description,omitempty"`
	CreatedTime time.Time          `goose:"index,createdAt" bson:"createdTime,omitempty"`
	UpdatedTime time.Time          `goose:"updatedAt" bson:"updatedTime,omitempty"`
	ViewCount   int64              `goose:"default=0" bson:"viewCount"`
	Rate        float64            `goose:"default=0" bson:"rate"`
	IsPublished bool               `goose:"default=false" bson:"isPublished"`
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
// 	err = userModel.Save()
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
	db, err := NewMongoDatabase(&DatabaseOptions{
		UsingEnv: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	userID := primitive.NewObjectID()
	user := User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example",
	}
	userModel := NewModel("TestUsers", &user)
	postModel := NewModel("TestPosts", &Post{
		UserID: userID,
		Title:  "test post",
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

	postModel = NewModel("TestPosts", &Post{})
	result, err := postModel.Populate("User").Find(bson.M{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
