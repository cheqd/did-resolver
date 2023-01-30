package types

func AddElemToSet(set []string, newElement string) []string {
	if set == nil {
		set = []string{}
	}
	for _, c := range set {
		if c == newElement {
			return set
		}
	}
	return append(set, newElement)
}
