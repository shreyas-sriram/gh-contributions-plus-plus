package draw

import (
	"testing"
)

func TestFindMax(t *testing.T) {
	tests := []struct {
		name string
		args []int
		want float32
	}{
		{name: "different numbers", args: []int{0, 3, 0, 10, 2, 7, 9}, want: 10.0},
		{name: "same numbers", args: []int{0, 0, 0}, want: 0.0},
	}

	for _, test := range tests {
		got := findMax(test.args)
		if got != test.want {
			t.Errorf("Got and want were incorrect, got: %+v, want: %+v", got, test.want)
		}
	}
}

func TestFindIntensities(t *testing.T) {
	tests := []struct {
		name string
		args []int
		want []intensity
	}{
		{name: "different numbers with 0", args: []int{3, 0, 7}, want: []intensity{level2, level0, level4}},
		{name: "different numbers without 0", args: []int{3, 2, 7}, want: []intensity{level2, level2, level4}},
		{name: "same numbers", args: []int{3, 3, 3}, want: []intensity{level4, level4, level4}},
	}

	for _, test := range tests {
		got := findIntensities(test.args)
		for i := range test.want {
			if got[i] != test.want[i] {
				t.Errorf("Got and want were incorrect, got: %+v, want: %+v", got[i], test.want[i])
			}
		}
	}
}
