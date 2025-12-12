package utils

// MAP
func Map[T any, R any](src []T, f func(T) R) []R {
	dst := make([]R, len(src))
	for i, v := range src {
		dst[i] = f(v)
	}
	return dst
}

// FILTER
func Filter[T any](src []T, f func(T) bool) []T {
	var res []T
	for _, v := range src {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}

// REDUCE
func Reduce[T any, R any](items []T, initial R, fn func(R, T) R) R {
	acc := initial
	for _, v := range items {
		acc = fn(acc, v)
	}
	return acc
}

func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
