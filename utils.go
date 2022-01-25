package wordle

func equal[T comparable](s1 []T, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i, v := range s1 {
		if v != s2[i] {
			return false
		}
	}

	return true
}

func containsSlice[T comparable](s [][]T, e []T) bool {
	for _, a := range s {
		if equal(a, e) {
			return true
		}
	}

	return false
}

func contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
