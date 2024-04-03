package utils

import (
	"errors"
	"strconv"
)

func StringsToInts(strings []string) ([]int, error) {
	var ints []int
	for _, s := range strings[1:] {
		i, err := strconv.Atoi(s)
		if err != nil {
			return []int{}, errors.Join(errors.New("cmd .stringsToInt(): can't parse int: "), err)
		}
		ints = append(ints, i)
	}
	return ints, nil
}
