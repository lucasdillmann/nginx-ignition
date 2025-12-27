package pagination

type Page[T any] struct {
	Contents   []T
	PageNumber int
	PageSize   int
	TotalItems int
}

func New[T any](pageNumber, pageSize, totalItems int, contents []T) *Page[T] {
	return &Page[T]{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalItems: totalItems,
		Contents:   contents,
	}
}
