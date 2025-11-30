package utils

import "errors"

// --- First-class & Higher-order examples ---

// Predicate adalah tipe fungsi yang mengembalikan bool untuk sebuah elemen Bibit-like
type Predicate[T any] func(T) bool

// Mapper memetakan T -> R
type Mapper[T any, R any] func(T) R

// Reducer mengurangi slice T menjadi R
type Reducer[T any, R any] func(R, T) R

// Filter generic
// func Filter[T any](items []T, pred Predicate[T]) []T {
// 	var out []T
// 	for _, it := range items {
// 		if pred(it) {
// 			out = append(out, it)
// 		}
// 	}
// 	return out
// }

// // Map generic
// func Map[T any, R any](items []T, mapper Mapper[T, R]) []R {
// 	out := make([]R, 0, len(items))
// 	for _, it := range items {
// 		out = append(out, mapper(it))
// 	}
// 	return out
// }

// // Reduce generic
// func Reduce[T any, R any](items []T, initial R, reducer Reducer[T, R]) R {
// 	acc := initial
// 	for _, it := range items {
// 		acc = reducer(acc, it)
// 	}
// 	return acc
// }

// Compose: komposisi dua fungsi R->S dan T->R => T->S
func Compose[T any, R any, S any](f func(R) S, g func(T) R) func(T) S {
	return func(t T) S {
		return f(g(t))
	}
}

// Closure contoh: pembuat counter
func NewCounter(start int) func() int {
	cnt := start
	return func() int {
		cnt++
		return cnt
	}
}

// Pure function example:
// Rekommendasi rule sederhana sebagai pure function (tidak akses luar, deterministic)
type SeedCandidate struct {
	Name     string
	Soil     string
	MaxRain  int // recommend if curah <= MaxRain
	MinStock int // require stok >= MinStock
}

// RecommendPure memilih kandidat berdasarkan parameter input secara murni
func RecommendPure(candidates []SeedCandidate, soil string, rain int, availableStock int) (SeedCandidate, error) {
	// gunakan filter & reduce dari atas
	filtered := Filter(candidates, func(c SeedCandidate) bool {
		return c.Soil == soil && rain <= c.MaxRain && availableStock >= c.MinStock
	})
	if len(filtered) == 0 {
		return SeedCandidate{}, errors.New("no candidate")
	}
	// pilih yang punya MinStock terbesar (contoh reduce)
	chosen := Reduce(filtered, filtered[0], func(acc SeedCandidate, next SeedCandidate) SeedCandidate {
		if next.MinStock > acc.MinStock {
			return next
		}
		return acc
	})
	return chosen, nil
}

// Contoh rekursi: traverse kategori (dummy recursion)
func DepthSum(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	return nums[0] + DepthSum(nums[1:])
}
