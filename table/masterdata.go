package table

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/stepupdream/go-support-tool/array"
	"github.com/stepupdream/go-support-tool/directory"
	"github.com/stepupdream/go-support-tool/file"
)

// MasterData is a struct used to represent tabular data.
type MasterData struct {
	name           string
	isPartialMatch bool
	extension      string
	Rows           map[Key]string
}

// NewTabular Create a new MasterData.
//
//goland:noinspection GoUnusedExportedFunction
func NewTabular(name string, extensionName string, rows map[Key]string, isPartialMatch bool) *MasterData {
	extension := "." + extensionName

	return &MasterData{
		name:           name,
		isPartialMatch: isPartialMatch,
		extension:      extension,
		Rows:           rows,
	}
}

// LoadByDirectoryPath Load the specified directory path.
// The directory path must be the path to the directory containing the insert, update, and delete directories.
func (m *MasterData) LoadByDirectoryPath(directoryPath string) error {
	// Avoid immediately UPDATING an INSET record within the same version (since it is an unintended update).
	loadTypes := []string{"delete", "update", "insert"}
	pathSeparator := string(os.PathSeparator)

	if !directoryExists(directoryPath, loadTypes) {
		return errors.New("Neither insert/update/delete directories were found : " + directoryPath)
	}

	var editIdsAll []int
	for _, loadType := range loadTypes {
		loadTypePath := directoryPath + pathSeparator + loadType + pathSeparator
		if !directory.Exist(loadTypePath) {
			continue
		}

		filePaths, err := directory.GetFilePathRecursive(loadTypePath, []string{m.extension})
		if err != nil {
			return err
		}

		for _, filePath := range filePaths {
			if !((m.isPartialMatch && strings.HasPrefix(file.BaseFileName(filePath), m.name)) || (m.name == file.BaseFileName(filePath))) {
				continue
			}

			var editMap map[Key]string
			editMap, err = LoadMap(filePath)
			if err != nil {
				return err
			}

			editIds := PluckId(editMap)
			editIdsAll = append(editIdsAll, editIds...)

			switch loadType {
			case "insert":
				err = m.insert(editMap, filePath)
			case "update":
				err = m.update(editMap, filePath)
			case "delete":
				err = m.delete(editMap, filePath)
			}

			if err != nil {
				return err
			}
		}
	}

	// Detect errors such as duplicate IDs for insert and update.
	// Logically, it's okay to have duplicate insert and update ids,
	// If it is duplicated, it is an error because it may be unintended input data.
	// [ex] when updating twice for the same id.
	if !array.IsUnique(editIdsAll) {
		return errors.New("ID is not unique : " + directoryPath + " " + m.name)
	}

	return nil
}

func directoryExists(directoryPath string, loadTypes []string) bool {
	pathSeparator := string(os.PathSeparator)
	for _, loadType := range loadTypes {
		if directory.Exist(directoryPath + pathSeparator + loadType + pathSeparator) {
			return true
		}
	}
	return false
}

// GetFilePathRecursive Get the file path of the specified extension recursively.
func (m *MasterData) GetFilePathRecursive(path string) ([]string, error) {
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
		if extension != m.extension {
			return nil
		}

		paths = append(paths, path)

		return nil
	})

	return paths, err
}

// delete the specified key from the map.
func (m *MasterData) delete(editMap map[Key]string, filePath string) error {
	baseIds := PluckId(m.Rows)

	for key := range editMap {
		if key.Key == "id" {
			if !array.Contains(baseIds, key.Id) {
				return errors.New("Attempted to delete a non-existent ID : id " + strconv.Itoa(key.Id) + " " + filePath)
			}
		}
		delete(m.Rows, Key{Id: key.Id, Key: key.Key})
	}

	return nil
}

// insert the specified key into the map.
func (m *MasterData) insert(editMap map[Key]string, filePath string) error {
	baseIds := PluckId(m.Rows)
	editIds := PluckId(editMap)

	for _, id := range editIds {
		if array.Contains(baseIds, id) {
			return errors.New("Tried to do an insert on an existing ID : id " + strconv.Itoa(id) + " " + filePath)
		}
	}

	for mapKey, value := range editMap {
		m.Rows[Key{Id: mapKey.Id, Key: mapKey.Key}] = value
	}

	return nil
}

// update the specified key in the map.
func (m *MasterData) update(editMap map[Key]string, filePath string) error {
	baseIds := PluckId(m.Rows)
	editIds := PluckId(editMap)
	for _, id := range editIds {
		if !array.Contains(baseIds, id) {
			return errors.New("Tried to update a non-existent ID : id " + strconv.Itoa(id) + " " + filePath)
		}
	}

	if err := m.delete(editMap, filePath); err != nil {
		return err
	}

	if err := m.insert(editMap, filePath); err != nil {
		return err
	}

	return nil
}
