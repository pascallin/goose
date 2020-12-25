package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/pascallin/goose"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := goose.NewMongoDatabase(&goose.DatabaseOptions{
		UsingEnv: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userID := primitive.NewObjectID()
	user := User{
		ID:    userID,
		Name:  "John Doe",
		Email: "john@example",
	}
	userModel := goose.NewModel("TestUsers", &user)
	postModel := goose.NewModel("TestPosts", &Post{
		UserID: userID,
		Title:  "test post",
	})
	user.Name = "Pascal Lin"
	err = userModel.Save()
	if err != nil {
		log.Fatal(err)
	}
	err = postModel.Save()
	if err != nil {
		log.Fatal(err)
	}

	postModel = goose.NewModel("TestPosts", &Post{})
	result, err := postModel.Populate("User").Find(bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("result: ", result)

	fmt.Println("======================")

	singleResult := goose.DB.Collection("TestUsers").FindOne(context.Background(), bson.M{"_id": userID})
	var userResult User
	singleResult.Decode(&userResult)
	fmt.Println("user: ", userResult)
}
