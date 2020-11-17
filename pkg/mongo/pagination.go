package mongo

import (
	"errors"
)

type Pagination struct {
	Page int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
}

func ValidatePagination(pagination *Pagination) (*Pagination, error) {
	if pagination == nil {
		return nil, errors.New("missing pagination")
	}
	if pagination.PageSize == 0 {
		pagination.PageSize = 10
	}
	if pagination.Page == 0 {
		pagination.Page = 1
	}
	return pagination, nil
}