package pagination

import "dillmann.com.br/nginx-ignition/core/common/pagination"

type PageDTO[T any] struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	Contents   []T `json:"contents"`
}

func Convert[I any, O any](
	page *pagination.Page[I],
	converter func(I) O,
) PageDTO[O] {

	var contents []O
	if page.Contents != nil {
		contents = make([]O, len(page.Contents))
		for index, item := range page.Contents {
			contents[index] = converter(item)
		}
	} else {
		contents = make([]O, 0)
	}

	return PageDTO[O]{
		PageNumber: page.PageNumber,
		PageSize:   page.PageSize,
		TotalItems: page.TotalItems,
		Contents:   contents,
	}
}
