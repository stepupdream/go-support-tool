package tabular

import (
	"bufio"
	"encoding/csv"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"github.com/stepupdream/golang-support-tool/array"
	"github.com/stepupdream/golang-support-tool/directory"
	supportFile "github.com/stepupdream/golang-support-tool/file"
	"github.com/stepupdream/golang-support-tool/logger"
)

// Key Make keys into structures to achieve multidimensional arrays.
type Key struct {
	Id  int
	Key string
}

// Tabular is a struct used to represent tabular data.
type Tabular struct {
	separatedType string
	extension     string
}

// NewTabular Create a new Tabular instance.
func NewTabular(separatedType string, extension string) *Tabular {
	return &Tabular{separatedType: separatedType, extension: extension}
}

// GetExtension Get the file extension.
func (tabular *Tabular) GetExtension() string {
	return tabular.extension
}

// Load the specified file.
func (tabular *Tabular) Load(filepath string, isRowExclusion bool, isColumnExclusion bool) [][]string {
	file, err := os.Open(filepath)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			logger.Fatal(err)
		}
	}()

	// If BOM is included, delete the BOM.
	// https://pinzolo.github.io/2017/03/29/utf8-csv-with-bom-on-golang.html
	ioReader := bufio.NewReader(file)
	if hasBOM(ioReader) {
		if _, err = ioReader.Discard(3); err != nil {
			logger.Fatal(err)
		}
	}

	csvReader := csv.NewReader(ioReader)
	if isRowExclusion {
		csvReader.Comment = '#'
	}
	if tabular.separatedType == "tsv" {
		csvReader.Comma = '\t'
		csvReader.LazyQuotes = true
	}

	rows, err := csvReader.ReadAll()
	if err != nil {
		logger.Fatal(err)
	}

	if isColumnExclusion {
		return tabular.exclusionColumn(rows, isColumnExclusion)
	}

	return rows
}

// hasBOM Check if the file has a BOM.
func hasBOM(reader *bufio.Reader) bool {
	bytes, err := reader.Peek(3)
	if err != nil || bytes[0] != 0xEF || bytes[1] != 0xBB || bytes[2] != 0xBF {
		return false
	}

	return true
}

// LoadMap Load the specified file and convert it to a map.
// If the file does not exist, return an empty map.
func (tabular *Tabular) LoadMap(filePath string, filterNames []string, isColumnExclusion bool) map[Key]string {
	if !supportFile.Exists(filePath) {
		return make(map[Key]string)
	}

	var filterColumnNumbers []int
	rows := tabular.Load(filePath, true, isColumnExclusion)
	if len(filterNames) != 0 {
		filterColumnNumbers = tabular.filterColumnNumbers(filePath, filterNames)
	}

	return tabular.convertMap(rows, filterColumnNumbers, filePath)
}

// exclusionColumn Exclude the column with the exclusion mark.
// If the exclusion mark is not set, return the original data.
func (tabular *Tabular) exclusionColumn(rows [][]string, isExclusion bool) [][]string {
	var disableColumnIndexes []int
	for index, value := range rows[0] {
		if isExclusion && value == "#" {
			disableColumnIndexes = append(disableColumnIndexes, index)
		}
	}

	var newRows [][]string
	for _, row := range rows {
		var newRow []string
		for index, value := range row {
			if !array.IntContains(disableColumnIndexes, index) {
				newRow = append(newRow, value)
			}
		}
		newRows = append(newRows, newRow)
	}

	return newRows
}

// convertMap
// Replacing separated value data (two-dimensional array of height and width) into a multidimensional associative array in a format
// that facilitates direct value specification by key.
func (tabular *Tabular) convertMap(rows [][]string, filterColumnNumbers []int, filepath string) map[Key]string {
	result := make(map[Key]string)
	keyName := map[int]string{}
	findIdColumn := false
	idColumnNumber := 0

	for rowNumber, row := range rows {
		for columnNumber, value := range row {
			// The first line is the key.
			if rowNumber == 0 {
				if value == "id" {
					findIdColumn = true
					idColumnNumber = columnNumber
				}
				keyName[columnNumber] = value
				continue
			}

			if len(filterColumnNumbers) != 0 && !array.IntContains(filterColumnNumbers, columnNumber) {
				continue
			}

			id, _ := strconv.Atoi(row[idColumnNumber])
			if _, flg := result[Key{id, keyName[columnNumber]}]; flg {
				logger.Fatal("ID is not unique : " + filepath)
			}
			if value == "" {
				logger.Fatal("Blank space is prohibited because it is impossible to determine if you forgot to enter the information. : " + filepath + " rowNumber : " + strconv.Itoa(rowNumber))
			}
			result[Key{id, keyName[columnNumber]}] = value
		}
	}

	if !findIdColumn {
		logger.Fatal("Separated value without ID column cannot be read : " + filepath)
	}

	return result
}

// PluckId Pluck the ID from the map.
func (tabular *Tabular) PluckId(valueMap map[Key]string) []int {
	var ids []int

	for mapKey := range valueMap {
		if mapKey.Key == "id" {
			ids = append(ids, mapKey.Id)
		}
	}

	sort.Ints(ids)

	return ids
}

// PluckKey Pluck the value of the specified key from the map.
func (tabular *Tabular) PluckKey(valueMap map[Key]string, key string) []string {
	var values []string

	for mapKey, value := range valueMap {
		if mapKey.Key == key {
			values = append(values, value)
		}
	}

	return values
}

// GetFilePathRecursive Get the file path of the specified extension recursively.
func (tabular *Tabular) GetFilePathRecursive(path string) ([]string, error) {
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
		if extension != tabular.extension {
			return nil
		}

		paths = append(paths, path)

		return nil
	})

	return paths, err
}

// CreateNewFile creates a new file at the specified path and writes the specified rows to it.
// If the file already exists, it will be overwritten.
func (tabular *Tabular) CreateNewFile(path string, rows [][]string) {
	// create allows you to create a new file and overwrite a new file.
	separatedFile, err := os.Create(path)
	if err != nil {
		logger.Fatal(err)
	}

	// Make it with BOM to avoid garbled characters.
	_, err = separatedFile.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		logger.Fatal(err)
	}

	writer := csv.NewWriter(separatedFile)
	defer writer.Flush()

	if tabular.separatedType == "tsv" {
		writer.Comma = '\t'
	}

	if err = writer.WriteAll(rows); err != nil {
		logger.Fatal(err)
	}
}

// delete Delete the specified key from the map.
func (tabular *Tabular) delete(baseMap map[Key]string, editMap map[Key]string, filePath string) map[Key]string {
	baseIds := tabular.PluckId(baseMap)

	for key := range editMap {
		if key.Key == "id" {
			if !array.IntContains(baseIds, key.Id) {
				logger.Fatal("Attempted to delete a non-existent ID : id " + strconv.Itoa(key.Id) + " " + filePath)
			}
		}
		delete(baseMap, Key{Id: key.Id, Key: key.Key})
	}

	return baseMap
}

// insert Insert the specified key into the map.
func (tabular *Tabular) insert(baseMap map[Key]string, editMap map[Key]string, filePath string) map[Key]string {
	baseIds := tabular.PluckId(baseMap)
	editIds := tabular.PluckId(editMap)

	for _, id := range editIds {
		if array.IntContains(baseIds, id) {
			logger.Fatal("Tried to do an insert on an existing ID : id " + strconv.Itoa(id) + " " + filePath)
		}
	}

	result := make(map[Key]string)

	for mapKey, value := range baseMap {
		result[Key{Id: mapKey.Id, Key: mapKey.Key}] = value
	}
	for mapKey, value := range editMap {
		result[Key{Id: mapKey.Id, Key: mapKey.Key}] = value
	}

	return result
}

// update Update the specified key in the map.
func (tabular *Tabular) update(baseMap map[Key]string, editMap map[Key]string, filePath string) map[Key]string {
	baseIds := tabular.PluckId(baseMap)
	editIds := tabular.PluckId(editMap)
	for _, id := range editIds {
		if !array.IntContains(baseIds, id) {
			logger.Fatal("Tried to update a non-existent ID : id " + strconv.Itoa(id) + " " + filePath)
		}
	}

	baseMap = tabular.delete(baseMap, editMap, filePath)
	baseMap = tabular.insert(baseMap, editMap, filePath)

	return baseMap
}

// filterColumnNumbers Get the column number of the column to filter.
func (tabular *Tabular) filterColumnNumbers(filepath string, filterColumnNames []string) []int {
	rows := tabular.Load(filepath, true, false)

	var columnNumbers []int
	for columnNumber, columnName := range rows[0] {
		if array.StrContains(filterColumnNames, columnName) {
			columnNumbers = append(columnNumbers, columnNumber)
		}
	}

	return columnNumbers
}

// LoadByDirectoryPath Load the specified directory path.
// The directory path must be the path to the directory containing the insert, update, and delete directories.
func (tabular *Tabular) LoadByDirectoryPath(directoryPath string, fileName string, baseMap map[Key]string, filterNames []string) map[Key]string {
	// Avoid immediately UPDATING an INSET record within the same version (since it is an unintended update).
	loadTypes := []string{"delete", "update", "insert"}
	if !directory.Exist(directoryPath+"/"+loadTypes[0]+"/") &&
		!directory.Exist(directoryPath+"/"+loadTypes[1]+"/") &&
		!directory.Exist(directoryPath+"/"+loadTypes[2]+"/") {
		logger.Fatal("Neither insert/update/delete directories were found : " + directoryPath)
	}

	var editIdsAll []int

	for _, loadType := range loadTypes {
		loadTypePath := directoryPath + "/" + loadType + "/"
		if !directory.Exist(loadTypePath) {
			continue
		}

		separatedValueFilePaths, err := tabular.GetFilePathRecursive(loadTypePath)
		if err != nil {
			logger.Fatal(err)
		}

		for _, filePath := range separatedValueFilePaths {
			if fileName != filepath.Base(filePath) {
				continue
			}

			var editSeparatedValueMap map[Key]string
			if len(filterNames) != 0 {
				editSeparatedValueMap = tabular.LoadMap(filePath, filterNames, false)
			} else {
				editSeparatedValueMap = tabular.LoadMap(filePath, filterNames, true)
			}

			editIds := tabular.PluckId(editSeparatedValueMap)
			editIdsAll = append(editIdsAll, editIds...)

			switch loadType {
			case "insert":
				baseMap = tabular.insert(baseMap, editSeparatedValueMap, filePath)
			case "update":
				baseMap = tabular.update(baseMap, editSeparatedValueMap, filePath)
			case "delete":
				baseMap = tabular.delete(baseMap, editSeparatedValueMap, filePath)
			}
		}
	}

	if !array.IsIntArrayUnique(editIdsAll) {
		logger.Fatal("ID is not unique : " + directoryPath + " " + fileName)
	}

	return baseMap
}
