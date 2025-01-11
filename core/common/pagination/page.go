package pagination

type Page[T any] struct {
	PageNumber int
	PageSize   int
	TotalItems int
	Contents   *[]T
}

func New[T any](pageNumber, pageSize, totalItems int, contents *[]T) *Page[T] {
	return &Page[T]{
		PageNumber: pageNumber,
		PageSize:   pageSize,
		TotalItems: totalItems,
		Contents:   contents,
	}
}
