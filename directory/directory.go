package directory

import (
	"github.com/pkg/errors"
	"github.com/stepupdream/go-support-tool/array"
	"github.com/stepupdream/go-support-tool/name"
	"io"
	"io/fs"
	"os"
	"path/filepath"
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

	for _, n := range names {
		if !array.Contains(exclusionTexts, n) {
			result = append(result, n)
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
func MaxFileName(directoryPath string) (r string, err error) {
	dirEntries, err := os.ReadDir(directoryPath)

	for _, dirEntry := range dirEntries {
		if r == "" {
			r = dirEntry.Name()
			continue
		}

		if r < dirEntry.Name() {
			r = dirEntry.Name()
		}
	}

	return r, nil
}

func MaxVersion(directoryPath string) (r string, err error) {
	dirEntries, err := os.ReadDir(directoryPath)
	var names []string

	for _, dirEntry := range dirEntries {
		names = append(names, dirEntry.Name())
	}

	if len(names) == 0 {
		return r, nil
	}

	namesSorted, err := name.CompareByNumericSegments(names)

	return namesSorted[len(namesSorted)-1], nil
}

// GetFilePathRecursive returns the path of the file in the specified directory.
func GetFilePathRecursive(path string, extensions []string) ([]string, error) {
	var paths []string

	// Recursively retrieve directories and files. (use WalkDir since Walk is now deprecated)
	err := filepath.WalkDir(path, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed filepath.WalkDir")
		}
		if info.IsDir() {
			return nil
		}
		extension := filepath.Ext(path)
		if len(extensions) > 0 && !array.Contains(extensions, extension) {
			return nil
		}
		paths = append(paths, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

// GetFilePath returns the path of the file in the specified directory.
func GetFilePath(directoryPath string, targetPath string) (string, error) {
	var filePath string

	// Recursively retrieve directories and files. (use WalkDir since Walk is now deprecated)
	err := filepath.WalkDir(directoryPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed filepath.WalkDir")
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) == targetPath {
			filePath = path
			return io.EOF // return io.EOF to stop the WalkDir function
		}
		return nil
	})

	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	return filePath, nil
}

func Create(targetDirectoryPath string, isGitkeep bool) error {
	if !Exist(targetDirectoryPath) {
		err := os.MkdirAll(targetDirectoryPath, 0755)
		if err != nil {
			return err
		}
	}

	if !isGitkeep {
		return nil
	}

	pathSeparator := string(os.PathSeparator)
	newFile, err := os.Create(targetDirectoryPath + pathSeparator + ".gitkeep")
	if err != nil {
		return err
	}
	defer func() {
		closeErr := newFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	return nil
}

func FindPrevious(arr []string, value string) (string, error) {
	if array.Contains(arr, value) {
		return value, nil
	}

	arr, e := name.CompareByNumericSegments(append(arr, value))
	if e != nil {
		return "", e
	}

	for i, current := range arr {
		if current == value {
			if i == 0 {
				return "", errors.New("No smaller or equivalent value found in the array")
			} else {
				return arr[i-1], nil
			}
		}
	}

	return "", errors.New("Value not found in array")
}
