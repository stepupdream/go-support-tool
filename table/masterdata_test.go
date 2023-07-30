package table

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestGetFilePathRecursive(t *testing.T) {
	separator := string(os.PathSeparator)

	type fields struct {
		name      string
		extension string
		rows      map[Key]string
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
				name:      "samples",
				extension: ".csv",
				rows:      nil,
			},
			args: args{
				path: "." + separator + "testdata" + separator + "sample",
			},
			want: []string{
				"testdata" + separator + "sample" + separator + "sample.csv",
				"testdata" + separator + "sample" + separator + "sample1.csv",
				"testdata" + separator + "sample" + separator + "sample2.csv",
				"testdata" + separator + "sample" + separator + "sample3.csv",
				"testdata" + separator + "sample" + separator + "sub" + separator + "sample1.csv",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MasterData{
				name:      tt.fields.name,
				extension: tt.fields.extension,
				Rows:      tt.fields.rows,
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
		name           string
		isPartialMatch bool
		extension      string
		rows           map[Key]string
	}
	type args struct {
		directoryPath string
	}
	tests := []struct {
		name           string
		isPartialMatch bool
		fields         fields
		args           args
		want           map[Key]string
		wantErr        bool
	}{
		{
			name: "LoadByDirectoryPath1-1",
			fields: fields{
				name:           "samples",
				isPartialMatch: false,
				extension:      ".csv",
				rows:           map[Key]string{},
			},
			args: args{
				directoryPath: "./testdata/pattern1",
			},
			want: map[Key]string{
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
				name:           "samples",
				isPartialMatch: false,
				extension:      ".csv",
				rows: map[Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./testdata/pattern1",
			},
			want: map[Key]string{
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
				name:           "samples",
				isPartialMatch: false,
				extension:      ".csv",
				rows: map[Key]string{
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
				directoryPath: "./testdata/pattern2",
			},
			want: map[Key]string{
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
				name:           "samples",
				isPartialMatch: false,
				extension:      ".csv",
				rows: map[Key]string{
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
				directoryPath: "./testdata/pattern3",
			},
			want: map[Key]string{
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
				name:           "samples",
				isPartialMatch: false,
				extension:      ".csv",
				rows: map[Key]string{
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
				directoryPath: "./testdata/pattern4",
			},
			want: map[Key]string{
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
				name:           "samples",
				isPartialMatch: false,
				extension:      ".csv",
				rows:           map[Key]string{},
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
				name:      "samples",
				extension: ".csv",
				rows: map[Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
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
				name:           "samples",
				isPartialMatch: false,
				extension:      ".csv",
				rows: map[Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
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
				name:           "samples",
				isPartialMatch: false,
				extension:      ".csv",
				rows: map[Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
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
				name:           "samples",
				isPartialMatch: false,
				extension:      ".csv",
				rows: map[Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
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
				name:           "samples",
				isPartialMatch: false,
				extension:      ".csv",
				rows: map[Key]string{
					{Id: 1, Key: "id"}:     "1",
					{Id: 1, Key: "sample"}: "aaa",
					{Id: 1, Key: "level"}:  "5",
				},
			},
			args: args{
				directoryPath: "./testdata/pattern10",
			},
			want:    map[Key]string{},
			wantErr: true,
		},
		{
			name: "LoadByDirectoryPath11",
			fields: fields{
				name:           "samples",
				isPartialMatch: true,
				extension:      ".csv",
				rows:           map[Key]string{},
			},
			args: args{
				directoryPath: "./testdata/pattern11",
			},
			want: map[Key]string{
				{Id: 2, Key: "id"}:       "2",
				{Id: 2, Key: "sample"}:   "bbb",
				{Id: 2, Key: "level"}:    "43",
				{Id: 100, Key: "id"}:     "100",
				{Id: 100, Key: "sample"}: "AAA",
				{Id: 100, Key: "level"}:  "1000",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MasterData{
				name:           tt.fields.name,
				isPartialMatch: tt.fields.isPartialMatch,
				extension:      tt.fields.extension,
				Rows:           tt.fields.rows,
			}
			err := m.LoadByDirectoryPath(tt.args.directoryPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilePathRecursive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				fmt.Println(err)
			} else if !reflect.DeepEqual(m.Rows, tt.want) {
				t.Errorf("GetFilePathRecursive() got = %v, want %v", m.Rows, tt.want)
			}
		})
	}
}
