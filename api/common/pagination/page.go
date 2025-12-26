package pagination

import "dillmann.com.br/nginx-ignition/core/common/pagination"

type PageDTO[T any] struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	Contents   []T `json:"contents"`
}

func Convert[I, O any](
	page *pagination.Page[I],
	converter func(*I) O,
) PageDTO[O] {
	contents := make([]O, 0)
	for index, item := range page.Contents {
		contents[index] = converter(&item)
	}

	return PageDTO[O]{
		PageNumber: page.PageNumber,
		PageSize:   page.PageSize,
		TotalItems: page.TotalItems,
		Contents:   contents,
	}
}
