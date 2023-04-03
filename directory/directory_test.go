package directory

import (
	"reflect"
	"testing"
)

func TestExistMulti(t *testing.T) {
	type args struct {
		parentPaths []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ExistMulti",
			args: args{
				parentPaths: []string{"../directory", "../excel"},
			},
			want: true,
		},
		{
			name: "ExistMult2",
			args: args{
				parentPaths: []string{"../directory", "../blank"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExistMulti(tt.args.parentPaths); got != tt.want {
				t.Errorf("ExistMulti() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxFileName(t *testing.T) {
	type args struct {
		directoryPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "MaxFileName",
			args: args{
				directoryPath: "../directory/testdata",
			},
			want: "1_0_1_0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxFileName(tt.args.directoryPath); got != tt.want {
				t.Errorf("MaxFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNames(t *testing.T) {
	type args struct {
		path           string
		exclusionTexts []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "GetNames1",
			args: args{
				path:           "../directory/testdata",
				exclusionTexts: []string{},
			},
			want: []string{"1_0_0_0", "1_0_1_0"},
		},
		{
			name: "GetNames2",
			args: args{
				path:           "../directory/testdata",
				exclusionTexts: []string{"1_0_1_0"},
			},
			want: []string{"1_0_0_0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := GetNames(tt.args.path, tt.args.exclusionTexts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExist(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Exist1",
			args: args{
				path: "../directory/testdata",
			},
			want: true,
		},
		{
			name: "Exist2",
			args: args{
				path: "../directory/test2",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Exist(tt.args.path); got != tt.want {
				t.Errorf("Exist() = %v, want %v", got, tt.want)
			}
		})
	}
}
