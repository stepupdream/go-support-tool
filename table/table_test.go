package table

import (
	"reflect"
	"testing"
)

func TestLoadMap(t *testing.T) {
	type args struct {
		filePath    string
		filterNames []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[Key]string
		wantErr bool
	}{
		{
			name: "LoadMap1",
			args: args{
				filePath:    "./testdata/sample.csv",
				filterNames: nil,
			},
			want: map[Key]string{
				{id: 2, key: "id"}:     "2",
				{id: 2, key: "sample"}: "bbb",
				{id: 2, key: "level"}:  "43",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadMap(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluckId(t *testing.T) {
	type args struct {
		valueMap map[Key]string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "PluckId1",
			args: args{
				valueMap: map[Key]string{
					{id: 1, key: "id"}:     "1",
					{id: 1, key: "sample"}: "aaa",
					{id: 1, key: "level"}:  "50",
				},
			},
			want: []int{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PluckId(tt.args.valueMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluckId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluckKey(t *testing.T) {
	type args struct {
		valueMap map[Key]string
		key      string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "PluckKey1",
			args: args{
				valueMap: map[Key]string{
					{id: 1, key: "id"}:     "1",
					{id: 1, key: "sample"}: "aaa",
					{id: 1, key: "level"}:  "50",
				},
				key: "sample",
			},
			want: []string{"aaa"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PluckKey(tt.args.valueMap, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluckKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
