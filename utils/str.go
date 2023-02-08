package utils

func IndexOf(array []string, searchElement string, fromIndex int) int {
	for i, v := range array[fromIndex:] {
		if v == searchElement {
			return fromIndex + i
		}
	}

	return -1
}

func Contains(vs []string, t string) bool {
	return IndexOf(vs, t, 0) >= 0
}
