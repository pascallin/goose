package test

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/pascallin/goose"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TestUser struct {
	ID          primitive.ObjectID `goose:"primary" bson:"_id,omitempty""`
	Name        string             `goose:"-" bson:"name,omitempty"`
	Email       string             `goose:"-" bson:"email,omitempty"`
	CreatedTime time.Time          `goose:"createdAt" bson:"createdTime,omitempty"`
	DeletedAt   time.Time          `goose:"deletedAt" bson:"deletedAt,omitempty"`
}

type TestPost struct {
	ID          primitive.ObjectID `goose:"primary,objectID" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `goose:"populate=User" bson:"userId,omitempty" ref:"TestUsers" forignKey:"_id"`
	Title       string             `goose:"-" bson:"title,omitempty"`
	Description string             `goose:"default='No description.'"  bson:"description,omitempty"`
	CreatedTime time.Time          `goose:"index,createdAt" bson:"createdTime,omitempty"`
	UpdatedTime time.Time          `goose:"updatedAt" bson:"updatedTime,omitempty"`
	ViewCount   int64              `goose:"default=0" bson:"viewCount"`
	Rate        float64            `goose:"default=0" bson:"rate"`
	Rate2       float64            `goose:"default=0.5" bson:"rate2,omitempty"`
	IsPublished bool               `goose:"default=false" bson:"isPublished"`
}

var (
	UserID = primitive.NewObjectID()
	PostID = primitive.NewObjectID()

	PostID1 = primitive.NewObjectID()
	PostID2 = primitive.NewObjectID()
)

func TestMain(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db, err := goose.NewMongoDatabase(&goose.DatabaseOptions{
		UsingEnv: true,
	})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	t.Run("Test Save", func(t *testing.T) {
		save(t)
	})
	t.Run("Test BulkWrite", func(t *testing.T) {
		batchWrite(t)
	})
	t.Run("Test Find Group", func(t *testing.T) {
		t.Run("Test Populate", func(t *testing.T) {
			populate(t)
		})
		t.Run("Test FindAndCount", func(t *testing.T) {
			populate(t)
		})
	})
	t.Run("Test Delete", func(t *testing.T) {
		delete(t)
	})
}

func checkDefault(t *testing.T) {
	postModel := goose.NewModel("TestPosts", &TestPost{})

	t.Log(postModel.CurValue)
}

func save(t *testing.T) {
	user := &TestUser{
		ID:    UserID,
		Name:  "John Doe",
		Email: "john@example",
	}
	userModel := goose.NewModel("TestUsers", user)
	postModel := goose.NewModel("TestPosts", &TestPost{
		UserID: UserID,
		Title:  "test post",
	})
	user.Name = "Pascal Lin"
	u, err := userModel.Save()
	if err != nil {
		t.Fatal(err)
	}
	p, err := postModel.Save()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(u)
	t.Log(p)
}

func batchWrite(t *testing.T) {
	postModel := goose.NewModel("TestPosts", &TestPost{})

	data := []interface{}{
		&TestPost{
			ID:     PostID,
			UserID: UserID,
			Title:  "This is Post",
		},
		&TestPost{
			ID:     PostID1,
			UserID: UserID,
			Title:  "This is Post 1",
		},
		&TestPost{
			ID:     PostID2,
			UserID: UserID,
			Title:  "This is Post 2",
		},
	}

	bulkResult, err := postModel.BulkInsert(data)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(bulkResult.InsertedCount)

	updateResult, err := postModel.UpdateMany(bson.M{
		"_id": bson.M{"$in": []primitive.ObjectID{PostID1, PostID2}},
	}, &TestPost{
		Description: "update description",
		IsPublished: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(updateResult.ModifiedCount, updateResult.UpsertedCount)

	bulkResult, err = postModel.BulkUpdate(bson.M{
		"isPublished": false,
	}, []interface{}{
		&TestPost{
			Description: "update not published description",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(bulkResult.ModifiedCount, bulkResult.UpsertedCount)
}

func findAndCount(t *testing.T) {
	postModel := goose.NewModel("TestPosts", &TestPost{})

	findResult, err := postModel.FindAndCount(bson.M{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(findResult)
}

func populate(t *testing.T) {
	postModel := goose.NewModel("TestPosts", &TestPost{})
	result, err := postModel.Populate("User").Find(bson.M{})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func delete(t *testing.T) {
	postModel := goose.NewModel("TestPosts", &TestPost{})
	postDeleted, err := postModel.DeleteOneByID(PostID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(postDeleted.DeletedCount)

	userModel := goose.NewModel("TestUsers", &TestUser{})
	userDeleted, err := userModel.SoftDeleteOne(bson.M{"_id": UserID})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(userDeleted.ModifiedCount)
}
