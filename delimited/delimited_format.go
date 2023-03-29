package delimited

import (
	"bufio"
	"encoding/csv"
	"os"
	"path/filepath"

	"github.com/stepupdream/golang-support-tool/array"
)

// Load the specified file.
func Load(targetPath string, isRowExclusion bool, isColumnExclusion bool) (rows [][]string, err error) {
	extension := filepath.Ext(targetPath)
	f, err := os.Open(targetPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	// If BOM is included, deleteData the BOM.
	// https://pinzolo.github.io/2017/03/29/utf8-csv-with-bom-on-golang.html
	ioReader := bufio.NewReader(f)
	if hasBOM(ioReader) {
		if _, err = ioReader.Discard(3); err != nil {
			return nil, err
		}
	}

	csvReader := csv.NewReader(ioReader)
	if isRowExclusion {
		csvReader.Comment = '#'
	}
	if extension == ".tsv" {
		csvReader.Comma = '\t'
		csvReader.LazyQuotes = true
	}

	rows, err = csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if isColumnExclusion {
		rows = exclusionColumn(rows, isColumnExclusion)
	}

	return rows, nil
}

// hasBOM Check if the file has a BOM.
func hasBOM(reader *bufio.Reader) bool {
	bytes, err := reader.Peek(3)
	if err != nil || bytes[0] != 0xEF || bytes[1] != 0xBB || bytes[2] != 0xBF {
		return false
	}

	return true
}

// exclusionColumn Exclude the column with the exclusion mark.
// If the exclusion mark is not set, return the original data.
func exclusionColumn(rows [][]string, isExclusion bool) [][]string {
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
			if !array.Contains(disableColumnIndexes, index) {
				newRow = append(newRow, value)
			}
		}
		newRows = append(newRows, newRow)
	}

	return newRows
}

// CreateNewFile creates a new file at the specified path and writes the specified rows to it.
// If the file already exists, it will be overwritten.
//
//goland:noinspection GoUnusedExportedFunction
func CreateNewFile(path string, rows [][]string) (err error) {
	separatedFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := separatedFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	// Make it with BOM to avoid garbled characters.
	_, err = separatedFile.Write([]byte{0xEF, 0xBB, 0xBF})
	if err != nil {
		return err
	}

	writer := csv.NewWriter(separatedFile)
	defer writer.Flush()

	if filepath.Ext(path) == "tsv" {
		writer.Comma = '\t'
	}

	if err = writer.WriteAll(rows); err != nil {
		return err
	}

	return nil
}
