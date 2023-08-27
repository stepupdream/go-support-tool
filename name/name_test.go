package name

import (
	"reflect"
	"testing"
)

func TestCompareByNumericSegments(t *testing.T) {
	type args struct {
		versionNames []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "CompareByNumericSegments1",
			args: args{
				versionNames: []string{"1_0_1", "1_0_2", "1_0_0"},
			},
			want:    []string{"1_0_0", "1_0_1", "1_0_2"},
			wantErr: false,
		},
		{
			name: "CompareByNumericSegments2",
			args: args{
				versionNames: []string{"1_10_0", "1_2_2", "1_0_0"},
			},
			want:    []string{"1_0_0", "1_2_2", "1_10_0"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SortByNumericSegments(tt.args.versionNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("SortByNumericSegments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortByNumericSegments() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsGreater(t *testing.T) {
	type args struct {
		v1 string
		v2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"IsGreater1",
			args{"1_10_0", "1_2_2"},
			true,
		},
		{
			"IsGreater2",
			args{"1_2_2", "1_10_0"},
			false,
		},
		{
			"IsGreater3",
			args{"1_0_1", "1_0_1_0"},
			false,
		},
		{
			"IsGreater4",
			args{"1_0_0", "1_0_0_0_1"},
			false,
		},
		{
			"IsGreater5",
			args{"10_0_10", "1_100_0_0_1"},
			true,
		},
		{
			"IsGreater6",
			args{"10_0_10", "10_1_10"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := IsGreaterVersion(tt.args.v1, tt.args.v2); got != tt.want {
				t.Errorf("isGreater() = %v, want %v", got, tt.want)
			}
		})
	}
}
