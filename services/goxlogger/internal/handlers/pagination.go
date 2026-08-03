package handlers

import (
	"errors"
	"net/http"
	"strconv"
)

type PaginationResponseTemplate[T []any] struct {
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
	Total   int `json:"total"`
	Pages   int `json:"pages"`
	Results T `json:"results"`
}

// GetPaginationResponse constructs a PaginationResponseTemplate with the provided data, limit, and offset.
func GetPaginationResponse[T []any](data T, limit, offset int) PaginationResponseTemplate[T] {
	var totalPages int

	if len(data) > 0 && limit > 0 {
		totalPages = (len(data) + limit - 1) / limit // Calculate total pages
	} else {
		totalPages = 0
	}

	return PaginationResponseTemplate[T]{
		Limit:   limit,
		Offset:  offset,
		Total:   len(data),
		Pages:   totalPages,
		Results: data,
	}
}

// PaginateData paginates the given data slice based on the "limit" and "offset" query parameters in the request.
func PaginateData[T any](r *http.Request, data []T) ([]T, error) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	if limit == "" {
		limit = "100"
	}

	if offset == "" {
		offset = "0"
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return data, err
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return data, err
	}

	if offsetInt < 0 || limitInt < 0 {
		return data, errors.New("invalid limit or offset")
	}

	endIndex := min(offsetInt+limitInt, len(data))

	return data[offsetInt:endIndex], nil
}
