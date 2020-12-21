package test

import (
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/pascallin/goose"
)

type User struct {
	ID    primitive.ObjectID `goose:"objectID,primary,required" bson="_id"`
	Name  string             `goose:"required" bson="name"`
	Email string             `goose:"required" bson="email"`
}

type Post struct {
	ID     primitive.ObjectID `goose:"objectID,primary" bson="_id"`
	UserId int                `goose:"index=1,populate=User" bson="userId"`
	Title  string             `goose:"-" bson="title"`
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

	user := &User{
		ID:    primitive.NewObjectID(),
		Name:  "John Doe",
		Email: "john@example",
	}
	userModel := goose.NewModel("TestUsers", user)
	// postModel := goose.NewModel("TestPosts", &Post{
	// 	ID:     primitive.NewObjectID(),
	// 	UserId: 1,
	// 	Title:  "test post",
	// })

	user.Name = "Pascal Lin"

	err = userModel.Save()
	if err != nil {
		t.Fatal(err)
	}
}
