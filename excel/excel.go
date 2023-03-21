package excel

import (
	"io/fs"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/stepupdream/golang-support-tool/log"
)

// GetFilePathRecursive returns the path of the Excel file in the specified directory.
//
//goland:noinspection GoUnusedExportedFunction
func GetFilePathRecursive(path string) []string {
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
		if extension != ".xlsx" && extension != ".xlsm" {
			return nil
		}

		paths = append(paths, path)

		return nil
	})

	if err != nil {
		log.Fatal("Failed to get the path to the Excel file", err)
	}

	return paths
}
