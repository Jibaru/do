package utils

// Ptr returns a pointer to the value passed as argument.
func Ptr[T any](v T) *T {
	return &v
}
