package name

import (
	"github.com/pkg/errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// CompareByNumericSegments The function compareByNumericSegments receives a slice of strings and sorts them by the numeric values within the strings.
func CompareByNumericSegments(versionNames []string) ([]string, error) {
	for _, versionName := range versionNames {
		if err := checkVersion(versionName); err != nil {
			return []string{}, err
		}
	}

	sort.Slice(versionNames, func(i, j int) bool {
		segments1 := strings.Split(versionNames[i], "_")
		segments2 := strings.Split(versionNames[j], "_")

		maxLen := len(segments1)
		if len(segments2) > maxLen {
			maxLen = len(segments2)
		}

		for k := 0; k < maxLen; k++ {
			var num1, num2 int

			if k >= len(segments1) {
				num1 = 0
			} else {
				num1, _ = strconv.Atoi(segments1[k])
			}

			if k >= len(segments2) {
				num2 = 0
			} else {
				num2, _ = strconv.Atoi(segments2[k])
			}

			if num1 != num2 {
				return num1 < num2
			}
		}

		return len(segments1) < len(segments2)
	})

	return versionNames, nil
}

func checkVersion(version string) error {
	re := regexp.MustCompile(`^(\d+(_\d+)*$)`)
	if !re.MatchString(version) {
		return errors.New("Invalid version format")
	}

	return nil
}
