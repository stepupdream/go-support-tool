package table

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetFilePathRecursive(t *testing.T) {
	type fields struct {
		name        string
		filterNames []string
		extension   string
		rows        map[Key]string
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
				path: "./testdata/sample",
			},
			want:    []string{"testdata/sample/sample.csv", "testdata/sample/sample1.csv", "testdata/sample/sample2.csv", "testdata/sample/sample3.csv", "testdata/sample/sub/sample1.csv"},
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
		rows        map[Key]string
	}
	type args struct {
		directoryPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[Key]string
		wantErr bool
	}{
		{
			name: "LoadByDirectoryPath1-1",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows:        map[Key]string{},
			},
			args: args{
				directoryPath: "./testdata/pattern1",
			},
			want: map[Key]string{
				{id: 2, key: "id"}:       "2",
				{id: 2, key: "sample"}:   "bbb",
				{id: 2, key: "level"}:    "43",
				{id: 100, key: "id"}:     "100",
				{id: 100, key: "sample"}: "AAA",
				{id: 100, key: "level"}:  "1000",
			},
			wantErr: false,
		},
		{
			name: "LoadByDirectoryPath1-2",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[Key]string{
					{id: 1, key: "id"}:     "1",
					{id: 1, key: "sample"}: "aaa",
					{id: 1, key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./testdata/pattern1",
			},
			want: map[Key]string{
				{id: 1, key: "id"}:       "1",
				{id: 1, key: "sample"}:   "aaa",
				{id: 1, key: "level"}:    "5",
				{id: 2, key: "id"}:       "2",
				{id: 2, key: "sample"}:   "bbb",
				{id: 2, key: "level"}:    "43",
				{id: 100, key: "id"}:     "100",
				{id: 100, key: "sample"}: "AAA",
				{id: 100, key: "level"}:  "1000",
			},
			wantErr: false,
		},
		{
			name: "LoadByDirectoryPath2",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[Key]string{
					{id: 1, key: "id"}:       "1",
					{id: 1, key: "sample"}:   "aaa",
					{id: 1, key: "level"}:    "5",
					{id: 2, key: "id"}:       "2",
					{id: 2, key: "sample"}:   "bbb",
					{id: 2, key: "level"}:    "43",
					{id: 100, key: "id"}:     "100",
					{id: 100, key: "sample"}: "AAA",
					{id: 100, key: "level"}:  "1000",
				},
			},
			args: args{
				directoryPath: "./testdata/pattern2",
			},
			want: map[Key]string{
				{id: 1, key: "id"}:       "1",
				{id: 1, key: "sample"}:   "aaa",
				{id: 1, key: "level"}:    "5",
				{id: 2, key: "id"}:       "2",
				{id: 2, key: "sample"}:   "ccc",
				{id: 2, key: "level"}:    "700",
				{id: 100, key: "id"}:     "100",
				{id: 100, key: "sample"}: "AAA",
				{id: 100, key: "level"}:  "1000",
			},
			wantErr: false,
		},
		{
			name: "LoadByDirectoryPath3",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[Key]string{
					{id: 1, key: "id"}:       "1",
					{id: 1, key: "sample"}:   "aaa",
					{id: 1, key: "level"}:    "5",
					{id: 2, key: "id"}:       "2",
					{id: 2, key: "sample"}:   "bbb",
					{id: 2, key: "level"}:    "43",
					{id: 100, key: "id"}:     "100",
					{id: 100, key: "sample"}: "AAA",
					{id: 100, key: "level"}:  "1000",
				},
			},
			args: args{
				directoryPath: "./testdata/pattern3",
			},
			want: map[Key]string{
				{id: 1, key: "id"}:       "1",
				{id: 1, key: "sample"}:   "aaa",
				{id: 1, key: "level"}:    "5",
				{id: 100, key: "id"}:     "100",
				{id: 100, key: "sample"}: "AAA",
				{id: 100, key: "level"}:  "1000",
			},
			wantErr: false,
		},
		{
			name: "LoadByDirectoryPath4",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[Key]string{
					{id: 1, key: "id"}:       "1",
					{id: 1, key: "sample"}:   "aaa",
					{id: 1, key: "level"}:    "5",
					{id: 2, key: "id"}:       "2",
					{id: 2, key: "sample"}:   "bbb",
					{id: 2, key: "level"}:    "43",
					{id: 100, key: "id"}:     "100",
					{id: 100, key: "sample"}: "AAA",
					{id: 100, key: "level"}:  "1000",
				},
			},
			args: args{
				directoryPath: "./testdata/pattern4",
			},
			want: map[Key]string{
				{id: 1, key: "id"}:       "1",
				{id: 1, key: "sample"}:   "nnn",
				{id: 1, key: "level"}:    "43",
				{id: 100, key: "id"}:     "100",
				{id: 100, key: "sample"}: "AAA",
				{id: 100, key: "level"}:  "1000",
				{id: 200, key: "id"}:     "200",
				{id: 200, key: "sample"}: "ttt",
				{id: 200, key: "level"}:  "9",
			},
			wantErr: false,
		},
		{
			name: "LoadByDirectoryPath5",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows:        map[Key]string{},
			},
			args: args{
				directoryPath: "./testdata/pattern5",
			},
			want:    map[Key]string{},
			wantErr: true,
		},
		{
			name: "LoadByDirectoryPath6",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[Key]string{
					{id: 1, key: "id"}:     "1",
					{id: 1, key: "sample"}: "aaa",
					{id: 1, key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./testdata/pattern6",
			},
			want:    map[Key]string{},
			wantErr: true,
		},
		{
			name: "LoadByDirectoryPath7",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[Key]string{
					{id: 1, key: "id"}:     "1",
					{id: 1, key: "sample"}: "aaa",
					{id: 1, key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./testdata/pattern7",
			},
			want:    map[Key]string{},
			wantErr: true,
		},
		{
			name: "LoadByDirectoryPath8",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[Key]string{
					{id: 1, key: "id"}:     "1",
					{id: 1, key: "sample"}: "aaa",
					{id: 1, key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./testdata/pattern8",
			},
			want:    map[Key]string{},
			wantErr: true,
		},
		{
			name: "LoadByDirectoryPath9",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[Key]string{
					{id: 1, key: "id"}:     "1",
					{id: 1, key: "sample"}: "aaa",
					{id: 1, key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./testdata/pattern9",
			},
			want:    map[Key]string{},
			wantErr: true,
		},
		{
			name: "LoadByDirectoryPath10",
			fields: fields{
				name:        "samples",
				filterNames: nil,
				extension:   ".csv",
				rows: map[Key]string{
					{id: 1, key: "id"}:     "1",
					{id: 1, key: "sample"}: "aaa",
					{id: 1, key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./testdata/pattern10",
			},
			want:    map[Key]string{},
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
