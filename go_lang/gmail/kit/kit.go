package kit

func Split(origin []string, limit int) [][]string {
	var result [][]string

	n := len(origin)
	for i := 0; i < n; i += limit {
		result = append(result, origin[i:min(i+limit, n)])
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
