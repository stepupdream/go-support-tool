package delimited

import (
	"os"
	"reflect"
	"testing"
)

// Use Test Main if you want to perform processing before and after the test.
func TestMain(m *testing.M) {
	// Run the test.
	code := m.Run()

	// Remove the test file.
	currentDir, _ := os.Getwd()
	_ = os.RemoveAll(currentDir + "/testdata/sample10.csv")

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
				path: "./testdata/sample10.csv",
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
			expected, err := os.ReadFile("./testdata/sample.csv")
			if err != nil {
				t.Fatal(err)
			}
			actual, err := os.ReadFile("./testdata/sample10.csv")
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
				targetPath:        "./testdata/sample.csv",
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
				targetPath:        "./testdata/sample.csv",
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
				targetPath:        "./testdata/sample.csv",
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
				targetPath:        "./testdata/sample.tsv",
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
