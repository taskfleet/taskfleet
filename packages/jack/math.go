package jack

import "golang.org/x/exp/constraints"

// Min returns the minimum of the two arguments.
func Min[T constraints.Ordered](lhs, rhs T) T {
	if lhs < rhs {
		return lhs
	}
	return rhs
}
