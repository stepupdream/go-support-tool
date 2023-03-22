package tabular

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stepupdream/golang-support-tool/file"
)

// Use Test Main if you want to perform processing before and after the test.
func TestMain(m *testing.M) {
	// Create a test directory.
	currentDir, _ := os.Getwd()
	dirPath := filepath.Join(currentDir, "test")
	_ = os.Mkdir(dirPath, 0777)
	fileNames := []string{"sample.csv", "sample2.csv", "sample3.csv"}

	for _, fileName := range fileNames {
		f, _ := os.Create(filepath.Join(dirPath, fileName))
		_, _ = f.WriteString("id,sample,#,level\n#1,aaa,2,13\n2,bbb,3,43")
		_ = f.Close()
	}

	// Run the test.
	code := m.Run()

	// Remove the test directory.
	_ = os.RemoveAll(dirPath)

	os.Exit(code)
}

func TestDeleteCSV(t *testing.T) {
	type args struct {
		baseCSV map[Key]string
		editCSV map[Key]string
	}
	tests := []struct {
		name string
		args args
		want map[Key]string
	}{
		{
			name: "DeleteBaseCSV",
			args: args{
				baseCSV: map[Key]string{
					{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
					{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "bbbb",
					{Id: 3, Key: "id"}: "3", {Id: 3, Key: "name"}: "cccc",
				},
				editCSV: map[Key]string{
					{Id: 3, Key: "id"}: "3", {Id: 3, Key: "name"}: "cccc",
				},
			},
			want: map[Key]string{
				{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
				{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "bbbb",
			},
		},
	}
	tabular := NewTabular("csv", ".csv")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tabular.delete(tt.args.baseCSV, tt.args.editCSV, ""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deleteCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertCSV(t *testing.T) {
	type args struct {
		baseCSV map[Key]string
		editCSV map[Key]string
	}
	tests := []struct {
		name string
		args args
		want map[Key]string
	}{
		{
			name: "InsertBaseCSV",
			args: args{
				baseCSV: map[Key]string{
					{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
					{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "bbbb",
					{Id: 3, Key: "id"}: "3", {Id: 3, Key: "name"}: "cccc",
				},
				editCSV: map[Key]string{
					{Id: 4, Key: "id"}: "3", {Id: 3, Key: "name"}: "dddd",
				},
			},
			want: map[Key]string{
				{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
				{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "bbbb",
				{Id: 3, Key: "id"}: "3", {Id: 3, Key: "name"}: "cccc",
				{Id: 4, Key: "id"}: "3", {Id: 3, Key: "name"}: "dddd",
			},
		},
	}
	separatedValue := NewTabular("csv", ".csv")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := separatedValue.insert(tt.args.baseCSV, tt.args.editCSV, ""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("insertCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateCSV(t *testing.T) {
	type args struct {
		baseCSV map[Key]string
		editCSV map[Key]string
	}
	tests := []struct {
		name string
		args args
		want map[Key]string
	}{
		{
			name: "InsertBaseCSV",
			args: args{
				baseCSV: map[Key]string{
					{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
					{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "bbbb",
				},
				editCSV: map[Key]string{
					{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "eeee",
				},
			},
			want: map[Key]string{
				{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
				{Id: 2, Key: "id"}: "2", {Id: 2, Key: "name"}: "eeee",
			},
		},
	}
	separatedValue := NewTabular("csv", ".csv")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := separatedValue.update(tt.args.baseCSV, tt.args.editCSV, ""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	type args struct {
		filepath          string
		isRowExclusion    bool
		isColumnExclusion bool
	}
	tests := []struct {
		name  string
		args  args
		want  [][]string
		want1 []string
	}{
		{
			name: "LoadCsv",
			args: args{
				filepath:          "test/sample2.csv",
				isRowExclusion:    false,
				isColumnExclusion: false,
			},
			want: [][]string{
				{"id", "sample", "#", "level"},
				{"#1", "aaa", "2", "13"},
				{"2", "bbb", "3", "43"},
			},
		},
		{
			name: "LoadCsv2",
			args: args{
				filepath:          "test/sample2.csv",
				isRowExclusion:    true,
				isColumnExclusion: true,
			},
			want: [][]string{
				{"id", "sample", "level"},
				{"2", "bbb", "43"},
			},
		},
	}
	separatedValue := NewTabular("csv", ".csv")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := separatedValue.Load(tt.args.filepath, tt.args.isRowExclusion, tt.args.isColumnExclusion)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadCsv() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluckId(t *testing.T) {
	type fields struct {
		separatedType string
		extension     string
	}
	type args struct {
		valueMap map[Key]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []int
	}{
		{
			name: "PluckId",
			fields: fields{
				separatedType: "csv",
				extension:     ".csv",
			},
			args: args{
				valueMap: map[Key]string{
					{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
					{Id: 2, Key: "id"}: "1", {Id: 2, Key: "name"}: "bbbb",
				},
			},
			want: []int{1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabular := &Tabular{
				separatedType: tt.fields.separatedType,
				extension:     tt.fields.extension,
			}
			if got := tabular.PluckId(tt.args.valueMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluckId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluckKey(t *testing.T) {
	type fields struct {
		separatedType string
		extension     string
	}
	type args struct {
		valueMap map[Key]string
		key      string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "PluckKey",
			fields: fields{
				separatedType: "csv",
				extension:     ".csv",
			},
			args: args{
				valueMap: map[Key]string{
					{Id: 1, Key: "id"}: "1", {Id: 1, Key: "name"}: "aaaa",
				},
				key: "name",
			},
			want: []string{
				"aaaa",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabular := &Tabular{
				separatedType: tt.fields.separatedType,
				extension:     tt.fields.extension,
			}
			if got := tabular.PluckKey(tt.args.valueMap, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluckKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFilePathRecursive(t *testing.T) {
	currentDir, _ := os.Getwd()
	dirPath := filepath.Join(currentDir, "test")
	file1 := filepath.Join(dirPath, "sample.csv")
	file2 := filepath.Join(dirPath, "sample2.csv")
	file3 := filepath.Join(dirPath, "sample3.csv")

	type fields struct {
		separatedType string
		extension     string
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
			name: "GetFilePathRecursive",
			fields: fields{
				separatedType: "csv",
				extension:     ".csv",
			},
			args: args{
				path: dirPath,
			},
			want:    []string{file1, file2, file3},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabular := &Tabular{
				separatedType: tt.fields.separatedType,
				extension:     tt.fields.extension,
			}
			got, err := tabular.GetFilePathRecursive(tt.args.path)
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

func TestCreateNewFile(t *testing.T) {
	currentDir, _ := os.Getwd()
	dirPath := filepath.Join(currentDir, "test")
	f := filepath.Join(dirPath, "sample3.csv")

	type fields struct {
		separatedType string
		extension     string
	}
	type args struct {
		path string
		rows [][]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "CreateNewFile",
			fields: fields{
				separatedType: "csv",
				extension:     ".csv",
			},
			args: args{
				path: f,
				rows: [][]string{
					{"id", "name"},
					{"1", "aaaa"},
					{"2", "bbbb"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tabular := &Tabular{
				separatedType: tt.fields.separatedType,
				extension:     tt.fields.extension,
			}
			tabular.CreateNewFile(tt.args.path, tt.args.rows)
		})
	}
	if !file.Exists(f) {
		t.Errorf("File creation failure. : CreateNewFile")
	}
}
