package pointers

func Dereference[T any](values []*T) []T {
	result := make([]T, len(values))
	for index, value := range values {
		result[index] = *value
	}
	return result
}

func Reference[T any](values []T) []*T {
	result := make([]*T, len(values))
	for index, value := range values {
		result[index] = &value
	}
	return result
}
