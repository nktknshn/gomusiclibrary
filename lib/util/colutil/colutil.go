package colutil

func MapSlice[T any, R any](s []T, f func(T) R) []R {
	res := make([]R, len(s))
	for idx, a := range s {
		res[idx] = f(a)
	}
	return res
}
