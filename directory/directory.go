package directory

import (
	"os"

	"github.com/stepupdream/golang-support-tool/array"
)

// Exist checks if the specified directory exists.
func Exist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// GetNames returns the file names in the specified directory.
// Specified exclusion texts are excluded.
func GetNames(path string, exclusionTexts []string) ([]string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := dir.Close()
		if err == nil {
			err = closeErr
		}
	}()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, name := range names {
		if !array.Contains(exclusionTexts, name) {
			result = append(result, name)
		}
	}

	return result, nil
}

// ExistMulti checks if any of the specified directories exist.
func ExistMulti(parentPaths []string) bool {
	isExist := false

	for _, path := range parentPaths {
		if Exist(path) {
			isExist = true
		}
	}

	return isExist
}

// MaxFileName returns the file name with the largest value in the specified directory.
func MaxFileName(directoryPath string) string {
	maxName := ""
	dirEntries, _ := os.ReadDir(directoryPath)
	for _, dirEntry := range dirEntries {
		if maxName == "" {
			maxName = dirEntry.Name()
			continue
		}

		if maxName < dirEntry.Name() {
			maxName = dirEntry.Name()
		}
	}

	return maxName
}
