package types

func Map2Array[T comparable, R any](themap map[T]R) []R {
	arr := make([]R, 0, len(themap))

	for _, item := range themap {
		arr = append(arr, item)
	}

	return arr
}
