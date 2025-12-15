package utils

func FromPtr[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}
