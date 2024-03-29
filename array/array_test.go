package array

import (
	"reflect"
	"testing"
)

func TestSliceString(t *testing.T) {
	type args struct {
		all   []string
		start string
		end   string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "max",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "1_0_0_0",
				end:   "max",
			},
			want: []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
		},
		{
			name: "max2",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "1_0_1_0",
				end:   "max",
			},
			want: []string{"1_0_1_0", "1_0_2_0"},
		},
		{
			name: "target",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "1_0_0_0",
				end:   "1_0_1_0",
			},
			want: []string{"1_0_0_0", "1_0_1_0"},
		},
		{
			name: "target2",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "1_0_0_0",
				end:   "1_0_2_0",
			},
			want: []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
		},
		{
			name: "startEmpty",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "",
				end:   "1_0_0_0",
			},
			want: []string{"1_0_0_0"},
		},
		{
			name: "next",
			args: args{
				all:   []string{"1_0_0_0", "1_0_1_0", "1_0_2_0"},
				start: "1_0_1_0",
				end:   "next",
			},
			want: []string{"1_0_1_0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := SliceString(tt.args.all, tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextArrayValue(t *testing.T) {
	type args struct {
		allValues []string
		nowValue  string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				allValues: []string{"a", "b", "c", "d"},
				nowValue:  "c",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test2",
			args: args{
				allValues: []string{"1_0_0_0", "1_0_1_0_0", "1_0_0_2", "1_0_1_3", "1_0_0_0_0_1"},
				nowValue:  "1_0_0_0_0_1",
			},
			want:    "1_0_0_2",
			wantErr: false,
		},
		{
			name: "test3",
			args: args{
				allValues: []string{"11_0_0_0", "21_0_11_0_0", "21_0_0_12", "11_0_11_13", "11_0_0_0_0_11"},
				nowValue:  "11_0_11_13",
			},
			want:    "21_0_0_12",
			wantErr: false,
		},
		{
			name: "test4",
			args: args{
				allValues: []string{"1_0_0_0", "1_0_1_0_0", "1_0_0_2"},
				nowValue:  "1_0_1_0_0",
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NextArrayValue(tt.args.allValues, tt.args.nowValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("NextArrayValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NextArrayValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringUnique(t *testing.T) {
	type args struct {
		values []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Unique",
			args: args{[]string{"a", "e", "b", "e", "d"}},
			want: []string{"a", "e", "b", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntUnique(t *testing.T) {
	type args struct {
		values []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Unique",
			args: args{[]int{1, 1, 2, 2, 3, 5}},
			want: []int{1, 2, 3, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPluckStringByIndex(t *testing.T) {
	type args struct {
		rows  [][]string
		index int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "PluckStringByIndex",
			args: args{
				rows:  [][]string{{"a", "aaa"}, {"b", "bbb"}, {"c", "ccc"}},
				index: 0,
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PluckStringByIndex(tt.args.rows, tt.args.index); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluckStringByIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntContains(t *testing.T) {
	type args struct {
		slice  []int
		target int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "IntContains1",
			args: args{
				slice:  []int{1, 2, 3, 4, 5},
				target: 10,
			},
			want: false,
		},
		{
			name: "IntContains2",
			args: args{
				slice:  []int{1, 2, 3, 4, 5},
				target: 1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.slice, tt.args.target); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIntArrayUnique(t *testing.T) {
	type args struct {
		args []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "IsIntArrayUnique1",
			args: args{
				args: []int{1, 2, 3, 4, 5},
			},
			want: true,
		},
		{
			name: "IsIntArrayUnique2",
			args: args{
				args: []int{1, 2, 3, 5, 5},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnique(tt.args.args); got != tt.want {
				t.Errorf("IsUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsStringArrayUnique(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "IsStringArrayUnique1",
			args: args{
				args: []string{"a", "b", "c", "d", "e"},
			},
			want: true,
		},
		{
			name: "IsStringArrayUnique2",
			args: args{
				args: []string{"a", "b", "c", "e", "e"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnique(tt.args.args); got != tt.want {
				t.Errorf("IsUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeMap(t *testing.T) {
	type args struct {
		m1 map[string]any
		m2 map[string]any
	}
	tests := []struct {
		name string
		args args
		want map[string]any
	}{
		{
			name: "MergeMap",
			args: args{
				m1: map[string]any{"a": "aaa", "b": "bbb"},
				m2: map[string]any{"c": "ccc"},
			},
			want: map[string]any{"a": "aaa", "b": "bbb", "c": "ccc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeMap(tt.args.m1, tt.args.m2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrContains(t *testing.T) {
	type args struct {
		slice  []string
		target string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "StrContains1",
			args: args{
				slice:  []string{"a", "b", "c", "d", "e"},
				target: "g",
			},
			want: false,
		},
		{
			name: "StrContains2",
			args: args{
				slice:  []string{"a", "b", "c", "d", "e"},
				target: "a",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.slice, tt.args.target); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
