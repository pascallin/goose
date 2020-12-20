package test

import (
	"testing"

	"github.com/joho/godotenv"

	"github.com/pascallin/goose"
)

type User struct {
	Id int `goose:"index=1,required"`
	Name string `goose:"required"`
	Email string `goose:"required"`
}

type Post struct {
	Id int `goose:"index=1"`
	UserId int `goose:"index=1,populate=User"`
	Title string `goose:"-"`
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

	userModel := goose.NewModel("TestUsers", User{
		Id: 1,
		Name: "John Doe",
		Email: "john@example",
	})
	postModel := goose.NewModel("TestPosts", Post{
		Id: 1,
		UserId: 1,
		Title: "test post",
	})

	userModel.PrintTags()
	postModel.GetFields()
}
