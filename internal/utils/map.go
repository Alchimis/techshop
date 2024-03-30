package utils

func Map[From interface{}, To interface{}](from []From, f func(From) To) []To {
	var to []To
	for _, i := range from {
		to = append(to, f(i))
	}
	return to
}
