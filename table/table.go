package table

import (
	"sort"
	"strconv"

	"github.com/pkg/errors"
	"github.com/stepupdream/golang-support-tool/array"
	"github.com/stepupdream/golang-support-tool/delimited"
	supportFile "github.com/stepupdream/golang-support-tool/file"
)

// Key Make keys into structures to achieve multidimensional arrays.
type Key struct {
	id  int
	key string
}

// LoadMap Load the specified file and convert it to a map.
// If the file does not exist, return an empty map.
//
//goland:noinspection GoUnusedExportedFunction
func LoadMap(filePath string, filterNames []string) (map[Key]string, error) {
	if !supportFile.Exists(filePath) {
		return make(map[Key]string), nil
	}

	rows, err := delimited.Load(filePath, true, true)
	if err != nil {
		return nil, err
	}

	numbers := filterColumnNumbers(rows[0], filterNames)

	return convertMap(rows, numbers, filePath)
}

// convertMap
// Replacing separated value data (two-dimensional array of height and width) into a multidimensional associative array in a format
// that facilitates direct value specification by key.
func convertMap(rows [][]string, filterColumnNumbers []int, filepath string) (map[Key]string, error) {
	convertedData := make(map[Key]string)
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

			if len(filterColumnNumbers) != 0 && !array.Contains(filterColumnNumbers, columnNumber) {
				continue
			}

			id, _ := strconv.Atoi(row[idColumnNumber])
			if _, flg := convertedData[Key{id, keyName[columnNumber]}]; flg {
				return nil, errors.New("Duplicate key : " + filepath + " rowNumber : " + strconv.Itoa(rowNumber))

			}
			if value == "" {
				return nil, errors.New("Empty value : " + filepath + " rowNumber : " + strconv.Itoa(rowNumber))
			}
			convertedData[Key{id, keyName[columnNumber]}] = value
		}
	}

	if !findIdColumn {
		return nil, errors.New("Not found id column : " + filepath)
	}

	return convertedData, nil
}

// PluckId Pluck the ID from the map.
//
//goland:noinspection GoUnusedExportedFunction
func PluckId(valueMap map[Key]string) (r []int) {
	for mapKey := range valueMap {
		if mapKey.key == "id" {
			r = append(r, mapKey.id)
		}
	}

	sort.Ints(r)

	return r
}

// PluckKey Pluck the value of the specified key from the map.
//
//goland:noinspection GoUnusedExportedFunction
func PluckKey(valueMap map[Key]string, key string) (r []string) {
	for mapKey, value := range valueMap {
		if mapKey.key == key {
			r = append(r, value)
		}
	}

	return r
}

// filterColumnNumbers Get the column number of the column to filter.
func filterColumnNumbers(filterRows []string, filterColumnNames []string) (r []int) {
	for columnNumber, columnName := range filterRows {
		if array.Contains(filterColumnNames, columnName) {
			r = append(r, columnNumber)
		}
	}

	return r
}
