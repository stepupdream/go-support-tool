package table

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"github.com/stepupdream/golang-support-tool/array"
	"github.com/stepupdream/golang-support-tool/directory"
	"github.com/stepupdream/golang-support-tool/file"
)

// MasterData is a struct used to represent tabular data.
type MasterData struct {
	name        string
	filterNames []string
	extension   string
	rows        map[Key]string
}

// NewTabular Create a new MasterData.
//
//goland:noinspection GoUnusedExportedFunction
func NewTabular(extension string, rows map[Key]string) *MasterData {
	return &MasterData{extension: extension, rows: rows}
}

// LoadByDirectoryPath Load the specified directory path.
// The directory path must be the path to the directory containing the insert, update, and delete directories.
func (m *MasterData) LoadByDirectoryPath(directoryPath string) error {
	// Avoid immediately UPDATING an INSET record within the same version (since it is an unintended update).
	loadTypes := []string{"delete", "update", "insert"}
	pathSeparator := string(os.PathSeparator)
	if !directory.Exist(directoryPath+pathSeparator+loadTypes[0]+pathSeparator) &&
		!directory.Exist(directoryPath+pathSeparator+loadTypes[1]+pathSeparator) &&
		!directory.Exist(directoryPath+pathSeparator+loadTypes[2]+pathSeparator) {
		return errors.New("Neither insert/update/delete directories were found : " + directoryPath)
	}

	var editIdsAll []int
	for _, loadType := range loadTypes {
		loadTypePath := directoryPath + pathSeparator + loadType + pathSeparator
		if !directory.Exist(loadTypePath) {
			continue
		}

		filePaths, err := m.GetFilePathRecursive(loadTypePath)
		if err != nil {
			return err
		}

		for _, filePath := range filePaths {
			if m.name != file.BaseFileName(filePath) {
				continue
			}

			var editMap map[Key]string
			editMap, err = LoadMap(filePath, m.filterNames)
			if err != nil {
				return err
			}

			editIds := PluckId(editMap)
			editIdsAll = append(editIdsAll, editIds...)

			switch loadType {
			case "insert":
				if err = m.insert(editMap, filePath); err != nil {
					return err
				}
			case "update":
				if err = m.update(editMap, filePath); err != nil {
					return err
				}
			case "delete":
				if err = m.delete(editMap, filePath); err != nil {
					return err
				}
			}
		}
	}

	// Detect errors such as duplicate IDs for insert and update.
	// Logically, it's okay to have duplicate insert and update ids,
	// If it is duplicated, it is an error because it may be unintended input data.
	// [ex] when updating twice for the same id.
	if !array.IsIntArrayUnique(editIdsAll) {
		return errors.New("ID is not unique : " + directoryPath + " " + m.name)
	}

	return nil
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
	baseIds := PluckId(m.rows)

	for key := range editMap {
		if key.key == "id" {
			if !array.IntContains(baseIds, key.id) {
				return errors.New("Attempted to delete a non-existent ID : id " + strconv.Itoa(key.id) + " " + filePath)
			}
		}
		delete(m.rows, Key{id: key.id, key: key.key})
	}

	return nil
}

// insert the specified key into the map.
func (m *MasterData) insert(editMap map[Key]string, filePath string) error {
	baseIds := PluckId(m.rows)
	editIds := PluckId(editMap)

	for _, id := range editIds {
		if array.IntContains(baseIds, id) {
			return errors.New("Tried to do an insert on an existing ID : id " + strconv.Itoa(id) + " " + filePath)
		}
	}

	for mapKey, value := range editMap {
		m.rows[Key{id: mapKey.id, key: mapKey.key}] = value
	}

	return nil
}

// update the specified key in the map.
func (m *MasterData) update(editMap map[Key]string, filePath string) error {
	baseIds := PluckId(m.rows)
	editIds := PluckId(editMap)
	for _, id := range editIds {
		if !array.IntContains(baseIds, id) {
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
