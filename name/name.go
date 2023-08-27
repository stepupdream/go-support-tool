package name

import (
	"github.com/pkg/errors"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func SortByNumericSegments(versionNames []string) ([]string, error) {
	if err := checkVersions(versionNames); err != nil {
		return []string{}, err
	}

	sort.Slice(versionNames, func(i, j int) bool {
		nums1 := getNumericSegments(versionNames[i])
		nums2 := getNumericSegments(versionNames[j])

		for k := 0; k < len(nums1) || k < len(nums2); k++ {
			num1, num2 := getSegmentNumbers(nums1, nums2, k)
			if num1 != num2 {
				return num1 < num2
			}
		}
		return len(nums1) < len(nums2)
	})

	return versionNames, nil
}

func IsGreaterVersion(referenceVersion string, comparisonVersion string) (bool, error) {
	if err := checkVersions([]string{referenceVersion, comparisonVersion}); err != nil {
		return false, err
	}

	v1Segments := getNumericSegments(referenceVersion)
	v2Segments := getNumericSegments(comparisonVersion)

	for i := 0; i < len(v1Segments) || i < len(v2Segments); i++ {
		num1, num2 := getSegmentNumbers(v1Segments, v2Segments, i)

		if num1 < num2 {
			return false, nil
		}
		if num1 > num2 {
			return true, nil
		}
	}

	return len(v1Segments) > len(v2Segments), nil
}

func checkVersions(versions []string) error {
	re := regexp.MustCompile(`^(\d+(_\d+)*$)`)
	for _, version := range versions {
		if !re.MatchString(version) {
			return errors.New("Invalid version format")
		}
	}
	return nil
}

func getNumericSegments(version string) []int {
	segments := strings.Split(version, "_")
	nums := make([]int, len(segments))
	for i, segment := range segments {
		nums[i], _ = strconv.Atoi(segment)
	}
	return nums
}

func getSegmentNumbers(nums1, nums2 []int, index int) (int, int) {
	num1, num2 := 0, 0
	if index < len(nums1) {
		num1 = nums1[index]
	}
	if index < len(nums2) {
		num2 = nums2[index]
	}
	return num1, num2
}
