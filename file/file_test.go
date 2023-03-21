package file

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBaseFileNames(t *testing.T) {
	type args struct {
		paths         []string
		withExtension bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "BaseFileNames1",
			args: args{
				paths:         []string{"C:/develop/aaa.csv", "C:/develop/bbb.csv"},
				withExtension: false,
			},
			want: []string{"aaa", "bbb"},
		},
		{
			name: "BaseFileNames2",
			args: args{
				paths:         []string{"C:/develop/aaa.csv", "C:/develop/bbb.csv"},
				withExtension: true,
			},
			want: []string{"aaa.csv", "bbb.csv"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BaseFileNames(tt.args.paths, tt.args.withExtension); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BaseFileNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Exists1",
			args: args{
				path: "../test/test.txt",
			},
			want: false,
		},
		{
			name: "Exists2",
			args: args{
				path: "../test/test.txt",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exists(tt.args.path); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecursiveFilePathInParent(t *testing.T) {
	type args struct {
		filename string
	}
	dirPath, _ := os.Getwd()
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "RecursiveFilePathInParent",
			args: args{
				filename: "README.md",
			},
			want: filepath.Dir(dirPath) + "/README.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RecursiveFilePathInParent(tt.args.filename); got != tt.want {
				t.Errorf("RecursiveFilePathInParent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveFileExtension(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "RemoveFileExtension",
			args: args{
				path: "C:/sample/aaa.csv",
			},
			want: "C:/sample/aaa",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveFileExtension(tt.args.path); got != tt.want {
				t.Errorf("RemoveFileExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}
