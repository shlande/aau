package classify

import (
	"reflect"
	"testing"
)

func Test_diff(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "contain", args: args{a: []string{"aa", "bb", "cc"}, b: []string{"aa", "bb"}}, want: nil},
		{name: "same", args: args{a: []string{"aa", "bb"}, b: []string{"aa", "bb"}}, want: nil},
		{name: "diff", args: args{a: []string{"aa", "bb"}, b: []string{"aa", "bb", "cc"}}, want: []string{"cc"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := diff(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("diff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Sort(t *testing.T) {

}
