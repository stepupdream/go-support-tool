package array

import (
	"github.com/pkg/errors"
)

// Contains checks if the specified value exists in the slice.
func Contains[T comparable](args []T, target T) bool {
	for _, arg := range args {
		if arg == target {
			return true
		}
	}

	return false
}

// MergeMap merges two maps.
func MergeMap(m1, m2 map[string]any) map[string]any {
	ans := make(map[string]any)

	for k, v := range m1 {
		ans[k] = v
	}
	for k, v := range m2 {
		ans[k] = v
	}

	return ans
}

// IsUnique checks if the specified array contains duplicate values.
func IsUnique[T comparable](args []T) bool {
	encountered := map[T]bool{}
	count := len(args)
	for i := 0; i < count; i++ {
		if !encountered[args[i]] {
			encountered[args[i]] = true
		} else {
			return false
		}
	}

	return true
}

// NextArrayValue returns the next value of the specified value in the array.
func NextArrayValue(allValues []string, nowValue string) (string, error) {
	if !Contains(allValues, nowValue) {
		return "", errors.New("Incorrect value specified. The specified value does not exist in the array : " + nowValue)
	}

	var nowKey int
	for key, value := range allValues {
		if value == nowValue {
			nowKey = key
		}
	}

	if len(allValues) < nowKey+2 {
		return "", nil
	}

	return allValues[nowKey+1], nil
}

// SliceString returns a slice of the specified array.
// If the start value is not specified, the first value of the array is used.
// If the end value is not specified, the last value of the array is used.
// If the end value is "next", the next value of the start value is used.
// If the end value is "max", the last value of the array is used.
func SliceString(all []string, start string, end string) (r []string, err error) {
	var tmp []string
	if start == "" {
		start = all[0]
	}
	isStart := false
	for _, value := range all {
		if value == start {
			isStart = true
		}
		if isStart {
			tmp = append(tmp, value)
		}
	}

	isEnd := false
	for _, value := range tmp {
		switch end {
		case "next":
			return []string{value}, nil
		case "max":
			r = append(r, value)
		default:
			if !Contains(all, end) {
				return nil, errors.New("The specified value could not be found : " + end)
			}
			if !isEnd {
				r = append(r, value)
			}
			if value == end {
				isEnd = true
			}
		}
	}

	return r, nil
}

// Unique returns an array with duplicate values removed.
func Unique[T comparable](values []T) (r []T) {
	tmp := make(map[T]bool)

	for _, value := range values {
		if !tmp[value] {
			tmp[value] = true
			r = append(r, value)
		}
	}

	return r
}

// PluckStringByIndex returns an array of the specified index of the specified array.
func PluckStringByIndex(rows [][]string, index int) (r []string) {
	for _, row := range rows {
		r = append(r, row[index])
	}

	return r
}
