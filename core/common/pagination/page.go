package pagination

type Page[T any] struct {
	PageNumber int64
	PageSize   int64
	TotalItems int64
	Contents   *[]T
}
