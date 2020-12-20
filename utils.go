package goose

import (
	"errors"
)

// ValidatePagination validate pagination
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
