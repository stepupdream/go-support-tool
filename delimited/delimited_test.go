package delimited

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Use Test Main if you want to perform processing before and after the test.
func TestMain(m *testing.M) {
	// Create a test directory.
	currentDir, _ := os.Getwd()
	dirPath := filepath.Join(currentDir, "test")
	_ = os.Mkdir(dirPath, 0777)
	fileNames := []string{"sample.csv"}

	for _, fileName := range fileNames {
		f, _ := os.Create(filepath.Join(dirPath, fileName))
		_, _ = f.WriteString("id,sample,#,level\n#1,aaa,2,13\n2,bbb,3,43\n")
		_ = f.Close()
	}

	fileNames = []string{"sample.tsv"}
	for _, fileName := range fileNames {
		f, _ := os.Create(filepath.Join(dirPath, fileName))
		_, _ = f.WriteString("id	sample	#	level\n#1	ccc	2	13\n2	ddd	3	43")
		_ = f.Close()
	}

	// Run the test.
	code := m.Run()

	// Remove the test directory.
	_ = os.RemoveAll(dirPath)

	os.Exit(code)
}

func TestCreateNewFile(t *testing.T) {
	type args struct {
		path string
		rows [][]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "CreateNewFile1",
			args: args{
				path: "./test/sample10.csv",
				rows: [][]string{
					{"id", "sample", "#", "level"},
					{"#1", "aaa", "2", "13"},
					{"2", "bbb", "3", "43"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateNewFile(tt.args.path, tt.args.rows); (err != nil) != tt.wantErr {
				t.Errorf("CreateNewFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			expected, err := os.ReadFile("./test/sample.csv")
			if err != nil {
				t.Fatal(err)
			}
			actual, err := os.ReadFile("./test/sample10.csv")
			if err != nil {
				t.Fatal(err)
			}

			// Remove BOM and compare.
			actual = actual[3:]
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected file does not match actual file")
			}
		})
	}
}

func TestLoad(t *testing.T) {
	type args struct {
		targetPath        string
		isRowExclusion    bool
		isColumnExclusion bool
	}
	tests := []struct {
		name     string
		args     args
		wantRows [][]string
		wantErr  bool
	}{
		{
			name: "Load1",
			args: args{
				targetPath:        "./test/sample.csv",
				isRowExclusion:    false,
				isColumnExclusion: false,
			},
			wantRows: [][]string{
				{"id", "sample", "#", "level"},
				{"#1", "aaa", "2", "13"},
				{"2", "bbb", "3", "43"},
			},
			wantErr: false,
		},
		{
			name: "Load2",
			args: args{
				targetPath:        "./test/sample.csv",
				isRowExclusion:    true,
				isColumnExclusion: false,
			},
			wantRows: [][]string{
				{"id", "sample", "#", "level"},
				{"2", "bbb", "3", "43"},
			},
			wantErr: false,
		},
		{
			name: "Load3",
			args: args{
				targetPath:        "./test/sample.csv",
				isRowExclusion:    true,
				isColumnExclusion: true,
			},
			wantRows: [][]string{
				{"id", "sample", "level"},
				{"2", "bbb", "43"},
			},
			wantErr: false,
		},
		{
			name: "Load4",
			args: args{
				targetPath:        "./test/sample.tsv",
				isRowExclusion:    true,
				isColumnExclusion: true,
			},
			wantRows: [][]string{
				{"id", "sample", "level"},
				{"2", "ddd", "43"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRows, err := Load(tt.args.targetPath, tt.args.isRowExclusion, tt.args.isColumnExclusion)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRows, tt.wantRows) {
				t.Errorf("Load() gotRows = %v, want %v", gotRows, tt.wantRows)
			}
		})
	}
}
