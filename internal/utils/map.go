package utils

func Map[From interface{}, To interface{}](from []From, f func(From) To) []To {
	var to []To
	for _, i := range from {
		to = append(to, f(i))
	}
	return to
}

func MapWithError[From interface{}, To interface{}](from []From, f func(From) (To, error)) ([]To, error) {
	var to []To
	for _, i := range from {
		t, err := f(i)
		if err != nil {
			return []To{}, err
		}
		to = append(to, t)
	}
	return to, nil
}

func Fold[From interface{}, To interface{}](from []From, start To, foldFunc func(From, To) To) To {
	var result To = start
	for _, v := range from {
		result = foldFunc(v, result)
	}
	return result
}
