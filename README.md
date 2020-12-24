# goose

A golang struct tags(schema) style mongo package hyper interface base on `mongo-driver`, inspired by Node.JS package `mongoose`

## Document

### Database

Once you connect database in main function, you can new a model in everywhere.

#### Connect to Mongo by `.env` file

```go
func main() {
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
  
  // do something else
}
```

#### Connect to Mongo

```go
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

  // do something else
```

### Model

#### Init a model

```go
// define struct as a schema
type User struct {
  ID    primitive.ObjectID `goose:"objectID,primary,required" bson="_id"`
  Name  string             `goose:"required" bson="name"`
  Email string             `goose:"required" bson="email"`
}

func main() {
  // init model with some data
  user := &User{
    ID:    primitive.NewObjectID(),
    Name:  "Pascal",
    Email: "pascal@example",
  }

  // new a model
  userModel := goose.NewModel("TestUsers", user)

  // update model data
  user.Name = "Pascal Lin"

  // save model, create or update
  err = userModel.Save()
  if err != nil {
    t.Fatal(err)
  }
}
```

### Tags

Using `goose`, you can using tags to specific some data relationship and normal business logic, there is the tag list below:

|TagName | Usage | Description|
|--- | --- | ---|
| primary | `goose="primary"` | define a primary key for you collection model |
| index | `goose="index"` | add field indexes to collection
| populate | `goose="populate=User"` or `goose="populate=User" ref="Users" foreignKey="userId"` | populate data from other collection, `ref` and `foreignKey` is optional |
| createdAt | `goose:"createdAt"` | set a created time
| updatedAt | `goose:"createdAt"` | set a updated time
| deletedAt | `goose:"createdAt"` | set a soft delete time and as soft delete signal
| - | `goose:"-"` | do nothing

A whole example:

```go
type Post struct {
  ID          primitive.ObjectID `goose:"primary" bson:"_id,omitempty"`
  UserID      primitive.ObjectID `goose:"populate=User" bson:"userId,omitempty" ref:"TestUsers" forignKey:"_id"`
  Title       string             `goose:"-" bson:"title,omitempty"`
  CreatedTime time.Time          `goose:"index,createdAt" bson:"createdTime,omitempty"`
  UpdatedTime time.Time          `goose:"updatedAt" bson:"updatedTime,omitempty"`
}
```

## Development

## Run godoc documents

```shell
git clone https://github.com/pascallin/goose.git

go get -v  golang.org/x/tools/cmd/godoc

godoc -http=:6060

# visit http://localhost:6060/pkg/github.com/pascallin/goose/pkg/mongo/
```

### Test

```shell script
go test -v ./test
```