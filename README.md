# goose

some database interface:

- `/pkg/mongo`: more human known interface mongo package, inspired by Node.JS package `mongoose`

## Docs

```shell
go get -v  golang.org/x/tools/cmd/godoc
godoc -http=:6060

# visit http://localhost:6060/pkg/github.com/pascallin/goose/pkg/mongo/
```

## Development

### Test

```shell script
go test -v ./test
```

## Strict rules

### Pagination

using pagination struct as below

``` golang
type Pagination struct {
  Page     int64 `json:"page"`
  PageSize int64 `json:"pageSize"`
}
```

you can validate pagination by using `ValidatePagination` method

### Soft delete

using `deletedAt` as soft delete specific field.

then you can using `SoftDeleteOne` and `SoftDeleteMany` to soft delete records.

### Query list data

the result of `FindAndCount` will be like as below:

```golang
type FindAndCountResult struct {
  Total int64
  Data  []bson.Raw
}
```
