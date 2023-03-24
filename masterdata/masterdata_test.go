package masterdata

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stepupdream/golang-support-tool/table"
)

// Use Test Main if you want to perform processing before and after the test.
func TestMain(m *testing.M) {
	// Create a test directory.
	currentDir, _ := os.Getwd()
	dirPath := filepath.Join(currentDir, "test")
	directoryNames := []string{dirPath, "test/sub"}
	for _, name := range directoryNames {
		_ = os.Mkdir(name, 0777)
	}

	fileNames := []string{"sample1.csv", "sample2.csv", "sample3.csv", "sub/sample1.csv"}
	for _, fileName := range fileNames {
		f, _ := os.Create(filepath.Join(dirPath, fileName))
		_, _ = f.WriteString("id,sample,#,level\n#1,aaa,2,13\n2,bbb,3,43\n")
		_ = f.Close()
	}

	// Run the test.
	code := m.Run()

	// Remove the test directory.
	_ = os.RemoveAll(dirPath)

	os.Exit(code)
}

func TestGetFilePathRecursive(t *testing.T) {
	type fields struct {
		name        string
		filterNames []string
		extension   string
		rows        map[table.Key]string
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "GetFilePathRecursive1",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows:        nil,
			},
			args: args{
				path: "./test",
			},
			want:    []string{"test/sample1.csv", "test/sample2.csv", "test/sample3.csv", "test/sub/sample1.csv"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MasterData{
				name:        tt.fields.name,
				filterNames: tt.fields.filterNames,
				extension:   tt.fields.extension,
				rows:        tt.fields.rows,
			}
			got, err := m.GetFilePathRecursive(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilePathRecursive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFilePathRecursive() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadByDirectoryPath(t *testing.T) {
	type fields struct {
		name        string
		filterNames []string
		extension   string
		rows        map[table.Key]string
	}
	type args struct {
		directoryPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[table.Key]string
		wantErr bool
	}{
		{
			name: "LoadByDirectoryPath1-1",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows:        map[table.Key]string{},
			},
			args: args{
				directoryPath: "./sample/pattern1",
			},
			want: map[table.Key]string{
				{Id: 2, Key: "id"}:       "2",
				{Id: 2, Key: "sample"}:   "bbb",
				{Id: 2, Key: "level"}:    "43",
				{Id: 100, Key: "id"}:     "100",
				{Id: 100, Key: "sample"}: "AAA",
				{Id: 100, Key: "level"}:  "1000",
			},
			wantErr: false,
		},
		{
			name: "LoadByDirectoryPath1-2",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[table.Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./sample/pattern1",
			},
			want: map[table.Key]string{
				{Id: 1, Key: "id"}:       "1",
				{Id: 1, Key: "sample"}:   "aaa",
				{Id: 1, Key: "level"}:    "5",
				{Id: 2, Key: "id"}:       "2",
				{Id: 2, Key: "sample"}:   "bbb",
				{Id: 2, Key: "level"}:    "43",
				{Id: 100, Key: "id"}:     "100",
				{Id: 100, Key: "sample"}: "AAA",
				{Id: 100, Key: "level"}:  "1000",
			},
			wantErr: false,
		},
		{
			name: "LoadByDirectoryPath2",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[table.Key]string{
					{Id: 1, Key: "id"}:       "1",
					{Id: 1, Key: "sample"}:   "aaa",
					{Id: 1, Key: "level"}:    "5",
					{Id: 2, Key: "id"}:       "2",
					{Id: 2, Key: "sample"}:   "bbb",
					{Id: 2, Key: "level"}:    "43",
					{Id: 100, Key: "id"}:     "100",
					{Id: 100, Key: "sample"}: "AAA",
					{Id: 100, Key: "level"}:  "1000",
				},
			},
			args: args{
				directoryPath: "./sample/pattern2",
			},
			want: map[table.Key]string{
				{Id: 1, Key: "id"}:       "1",
				{Id: 1, Key: "sample"}:   "aaa",
				{Id: 1, Key: "level"}:    "5",
				{Id: 2, Key: "id"}:       "2",
				{Id: 2, Key: "sample"}:   "ccc",
				{Id: 2, Key: "level"}:    "700",
				{Id: 100, Key: "id"}:     "100",
				{Id: 100, Key: "sample"}: "AAA",
				{Id: 100, Key: "level"}:  "1000",
			},
			wantErr: false,
		},
		{
			name: "LoadByDirectoryPath3",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[table.Key]string{
					{Id: 1, Key: "id"}:       "1",
					{Id: 1, Key: "sample"}:   "aaa",
					{Id: 1, Key: "level"}:    "5",
					{Id: 2, Key: "id"}:       "2",
					{Id: 2, Key: "sample"}:   "bbb",
					{Id: 2, Key: "level"}:    "43",
					{Id: 100, Key: "id"}:     "100",
					{Id: 100, Key: "sample"}: "AAA",
					{Id: 100, Key: "level"}:  "1000",
				},
			},
			args: args{
				directoryPath: "./sample/pattern3",
			},
			want: map[table.Key]string{
				{Id: 1, Key: "id"}:       "1",
				{Id: 1, Key: "sample"}:   "aaa",
				{Id: 1, Key: "level"}:    "5",
				{Id: 100, Key: "id"}:     "100",
				{Id: 100, Key: "sample"}: "AAA",
				{Id: 100, Key: "level"}:  "1000",
			},
			wantErr: false,
		},
		{
			name: "LoadByDirectoryPath4",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[table.Key]string{
					{Id: 1, Key: "id"}:       "1",
					{Id: 1, Key: "sample"}:   "aaa",
					{Id: 1, Key: "level"}:    "5",
					{Id: 2, Key: "id"}:       "2",
					{Id: 2, Key: "sample"}:   "bbb",
					{Id: 2, Key: "level"}:    "43",
					{Id: 100, Key: "id"}:     "100",
					{Id: 100, Key: "sample"}: "AAA",
					{Id: 100, Key: "level"}:  "1000",
				},
			},
			args: args{
				directoryPath: "./sample/pattern4",
			},
			want: map[table.Key]string{
				{Id: 1, Key: "id"}:       "1",
				{Id: 1, Key: "sample"}:   "nnn",
				{Id: 1, Key: "level"}:    "43",
				{Id: 100, Key: "id"}:     "100",
				{Id: 100, Key: "sample"}: "AAA",
				{Id: 100, Key: "level"}:  "1000",
				{Id: 200, Key: "id"}:     "200",
				{Id: 200, Key: "sample"}: "ttt",
				{Id: 200, Key: "level"}:  "9",
			},
			wantErr: false,
		},
		{
			name: "LoadByDirectoryPath5",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows:        map[table.Key]string{},
			},
			args: args{
				directoryPath: "./sample/pattern5",
			},
			want:    map[table.Key]string{},
			wantErr: true,
		},
		{
			name: "LoadByDirectoryPath6",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[table.Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./sample/pattern6",
			},
			want:    map[table.Key]string{},
			wantErr: true,
		},
		{
			name: "LoadByDirectoryPath7",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[table.Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./sample/pattern7",
			},
			want:    map[table.Key]string{},
			wantErr: true,
		},
		{
			name: "LoadByDirectoryPath8",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[table.Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./sample/pattern8",
			},
			want:    map[table.Key]string{},
			wantErr: true,
		},
		{
			name: "LoadByDirectoryPath9",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[table.Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./sample/pattern9",
			},
			want:    map[table.Key]string{},
			wantErr: true,
		},
		{
			name: "LoadByDirectoryPath10",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[table.Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./sample/pattern10",
			},
			want:    map[table.Key]string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MasterData{
				name:        tt.fields.name,
				filterNames: tt.fields.filterNames,
				extension:   tt.fields.extension,
				rows:        tt.fields.rows,
			}
			err := m.LoadByDirectoryPath(tt.args.directoryPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilePathRecursive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println(err)
			} else if !reflect.DeepEqual(m.rows, tt.want) {
				t.Errorf("GetFilePathRecursive() got = %v, want %v", m.rows, tt.want)
			}
		})
	}
}
