package array

import (
	"github.com/stepupdream/golang-support-tool/log"
)

// StrContains checks if the specified string exists in the array.
func StrContains(slice []string, target string) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}

	return false
}

// IntContains checks if the specified integer exists in the array.
func IntContains(slice []int, target int) bool {
	for _, value := range slice {
		if value == target {
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

// IsIntArrayUnique checks if the specified array contains duplicate values.
func IsIntArrayUnique(args []int) bool {
	encountered := map[int]bool{}
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

// IsStringArrayUnique checks if the specified array contains duplicate values.
func IsStringArrayUnique(args []string) bool {
	encountered := map[string]bool{}
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
func NextArrayValue(allValues []string, nowValue string) string {
	if !StrContains(allValues, nowValue) {
		log.Fatal("Incorrect value specified. The specified value does not exist in the array : " + nowValue)
	}

	var nowKey int
	for key, value := range allValues {
		if value == nowValue {
			nowKey = key
		}
	}

	if len(allValues) < nowKey+2 {
		return ""
	}

	return allValues[nowKey+1]
}

// SliceString returns a slice of the specified array.
// If the start value is not specified, the first value of the array is used.
// If the end value is not specified, the last value of the array is used.
// If the end value is "next", the next value of the start value is used.
// If the end value is "max", the last value of the array is used.
func SliceString(all []string, start string, end string) []string {
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

	var result []string
	isEnd := false
	for _, value := range tmp {
		switch end {
		case "next":
			return []string{value}
		case "max":
			result = append(result, value)
		default:
			if !StrContains(all, end) {
				log.Fatal("The specified value could not be found : " + end)
			}
			if !isEnd {
				result = append(result, value)
			}

			if value == end {
				isEnd = true
			}
		}
	}

	return result
}

// StringUnique returns an array with duplicate values removed.
func StringUnique(values []string) []string {
	tmp := make(map[string]bool)
	var result []string

	for _, value := range values {
		if !tmp[value] {
			tmp[value] = true
			result = append(result, value)
		}
	}

	return result
}

// IntUnique returns an array with duplicate values removed.
func IntUnique(values []int) []int {
	tmp := make(map[int]bool)
	var result []int

	for _, value := range values {
		if !tmp[value] {
			tmp[value] = true
			result = append(result, value)
		}
	}

	return result
}

// PluckStringByIndex returns an array of the specified index of the specified array.
func PluckStringByIndex(rows [][]string, index int) []string {
	var result []string

	for _, row := range rows {
		result = append(result, row[index])
	}

	return result
}
