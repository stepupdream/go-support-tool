package excel

import (
	"github.com/stepupdream/go-support-tool/directory"
)

// GetFilePathRecursive returns the path of the Excel file in the specified directory.
//
//goland:noinspection GoUnusedExportedFunction
func GetFilePathRecursive(path string) ([]string, error) {
	paths, err := directory.GetFilePathRecursive(path, []string{".xlsm", ".xlsx"})
	if err != nil {
		return nil, err
	}

	return paths, nil
}
