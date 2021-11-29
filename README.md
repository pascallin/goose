# goose

A golang struct tags(schema) style mongo package hyper interface base on `mongo-driver`, inspired by Node.JS package [mongoose](https://github.com/Automattic/mongoose)

## Examples

You can check the `/examples/` folder.

```shell
cd ./examples && go run ./simple.go
```

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
| primary | `goose:"primary"` | define a primary key for you collection model, default will set model primary key `_id` |
| index | `goose:"index"` | add field indexes to collection |
| default |  `goose:"default='test'"` or `goose:"default=1"` or `goose:"default=1.1"` or `goose:"default=false"` | set default value for model field, `string` should be quote by `'` and not including `,`; int and float will convert to 64 bit, you should not add `bson:omitempty` if `default=0` |
| populate | `goose:"populate=Users"` or `goose:"populate=User" ref="Users" foreignKey="_id"` | populate data from other collection, if not setting `ref` and `foreignKey`, populate should be `populate=[COLLECTION_NAME]` and default foreignKey is `_id`  |
| createdAt | `goose:"createdAt"` | set field as created time
| updatedAt | `goose:"updatedAt"` | set field as updated time
| deletedAt | `goose:"deletedAt"` |  set field as soft delete time
| - | `goose:"-"` | do nothing

A whole example:

```go
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
```

### Collection

You can still using collection from `mongo-driver` database as usual. such as,

```go
singleResult := goose.DB.Collection("TestUsers").FindOne(context.Background(), bson.M{"_id": userID})
var userResult User
singleResult.Decode(&userResult)
fmt.Println("user: ", userResult)
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
