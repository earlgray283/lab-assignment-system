package lib

func PointerOfValue[T any](v T) *T {
	return &v
}
