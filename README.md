# goose

some database interface:

- `/pkg/mongo`: inspired by Node.JS package `mongoose`

## test

```shell script
go test -v ./test
```

## strict rules

### pagination

using pagination struct as below

``` golang
type Pagination struct {
  Page     int64 `json:"page"`
  PageSize int64 `json:"pageSize"`
}
```

you can validate pagination by using `ValidatePagination` method

### soft delete

using `deletedAt` as soft delete specific field.

then you can using `SoftDeleteOne` and `SoftDeleteMany` to soft delete records.

### query list data

the result of `FindAndCount` will be like as below:

```golang
type FindAndCountResult struct {
  Total int64
  Data  []bson.Raw
}
```
