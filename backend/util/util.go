package util

func Filter[T any](vals []T, test func(T) bool) (ret []T) {
	for _, val := range vals {
		if test(val) {
			ret = append(ret, val)
		}
	}
	return ret
}

func Map[T any, U any](vals []T, f func(T) U) (ret []U) {
	for _, val := range vals {
		ret = append(ret, f(val))
	}
	return ret
}
