package directory

import (
	"os"

	"github.com/stepupdream/golang-support-tool/array"
)

// Exist checks if the specified directory exists.
func Exist(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

// GetNames returns the file names in the specified directory.
// Specified exclusion texts are excluded.
func GetNames(path string, exclusionTexts []string) (result []string, err error) {
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

	for _, name := range names {
		if !array.Contains(exclusionTexts, name) {
			result = append(result, name)
		}
	}

	return result, nil
}

// ExistMulti checks if any of the specified directories exist.
func ExistMulti(parentPaths []string) (r bool) {
	for _, path := range parentPaths {
		if Exist(path) {
			r = true
		}
	}

	return r
}

// MaxFileName returns the file name with the largest value in the specified directory.
func MaxFileName(directoryPath string) (r string) {
	dirEntries, _ := os.ReadDir(directoryPath)
	for _, dirEntry := range dirEntries {
		if r == "" {
			r = dirEntry.Name()
			continue
		}

		if r < dirEntry.Name() {
			r = dirEntry.Name()
		}
	}

	return r
}
