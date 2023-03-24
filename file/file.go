package file

import (
	"io"
	"os"
	"path/filepath"

	"github.com/stepupdream/golang-support-tool/directory"
	"github.com/stepupdream/golang-support-tool/logger"
)

// Exists checks if the specified file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// RecursiveFilePathInParent returns the path of the specified file in the parent directory.
// If the specified file is not found, it will be searched recursively in the parent directory.
func RecursiveFilePathInParent(filename string) string {
	pathSeparator := string(os.PathSeparator)
	dirPath, _ := os.Getwd()
	for i := 0; i < 10; i++ {
		findPath := dirPath + pathSeparator + filename
		if Exists(findPath) {
			return findPath
		}

		dirPath = filepath.Dir(dirPath)
	}
	logger.Fatal("The specified file could not be found : " + filename)

	return ""
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
func BaseFileNames(paths []string, withExtension bool) []string {
	var names []string
	for _, name := range paths {
		if withExtension {
			names = append(names, filepath.Base(name))
		} else {
			names = append(names, filepath.Base(RemoveFileExtension(name)))
		}
	}

	return names
}

// Copy copies the specified file.
//
//goland:noinspection GoUnusedExportedFunction
func Copy(basedPath string, targetPath string) {
	if !directory.Exist(filepath.Dir(targetPath)) {
		err := os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			logger.Fatal(err)
		}
	}

	newFile, err := os.Create(targetPath)
	if err != nil {
		logger.Fatal(err)
	}

	oldFile, err := os.Open(basedPath)
	if err != nil {
		logger.Fatal(err)
	}

	_, err = io.Copy(newFile, oldFile)
	if err != nil {
		logger.Fatal(err)
	}
}
