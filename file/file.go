package file

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/stepupdream/golang-support-tool/directory"
)

// Exists checks if the specified file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// RecursiveFilePathInParent returns the path of the specified file in the parent directory.
// If the specified file is not found, it will be searched recursively in the parent directory.
func RecursiveFilePathInParent(filename string) (string, error) {
	pathSeparator := string(os.PathSeparator)
	dirPath, _ := os.Getwd()
	for i := 0; i < 10; i++ {
		findPath := dirPath + pathSeparator + filename
		if Exists(findPath) {
			return findPath, nil
		}

		dirPath = filepath.Dir(dirPath)
	}

	return "", errors.New("The specified file could not be found : " + filename)
}

// RemoveFileExtension removes the file extension.
func RemoveFileExtension(path string) string {
	return path[:len(path)-len(filepath.Ext(path))]
}

// BaseFileName returns the file name of the specified path.
func BaseFileName(path string) string {
	return BaseFileNames([]string{path}, false)[0]
}

// BaseFileNames returns the file name of the specified path.
func BaseFileNames(paths []string, withExtension bool) (r []string) {
	for _, name := range paths {
		if withExtension {
			r = append(r, filepath.Base(name))
		} else {
			r = append(r, filepath.Base(RemoveFileExtension(name)))
		}
	}

	return r
}

// Copy copies the specified file.
//
//goland:noinspection GoUnusedExportedFunction
func Copy(basedPath string, targetPath string) error {
	if !directory.Exist(filepath.Dir(targetPath)) {
		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			return err
		}
	}

	newFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	oldFile, err := os.Open(basedPath)
	if err != nil {
		return err
	}

	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		return err
	}

	return nil
}
