package data

import (
	"testing"
)

func TestParseDateData(t *testing.T) {
	type data struct {
		date string
		data int
	}

	type input struct {
		rawHTML string
		year    string
	}

	tests := []struct {
		name string
		args input
		want []data
	}{
		{name: "simple match", args: input{"data-count=\"10\" data-date=\"2020-11-22\"", "2020"}, want: []data{data{"2020-11-22", 10}}},
		{name: "match html string", args: input{"<foobar>data-count=\"10\" data-date=\"2020-11-22\"</foobar>", "2020"}, want: []data{data{"2020-11-22", 10}}},
		{name: "no matches", args: input{"<foobar></foobar>", "2020"}, want: []data{}},
	}

	for _, test := range tests {
		got := parseDateData(test.args.rawHTML, test.args.year)
		for i, gotData := range got {
			if gotData.Date != test.want[i].date || gotData.Data != test.want[i].data {
				t.Errorf("Got and want were incorrect, got: %+v, want: %+v", gotData, test.want[i])
			}
		}
	}
}

func TestParseTotal(t *testing.T) {
	tests := []struct {
		name string
		args string
		want int
	}{
		{name: "simple match", args: "1337 contributions this year", want: 1337},
		{name: "no match", args: "1337 contribution this year", want: 0},
		{name: "empty", args: "", want: 0},
	}

	for _, test := range tests {
		got := parseTotal(test.args)
		if got != test.want {
			t.Errorf("Got and want were incorrect, got: %d, want: %d", got, test.want)
		}
	}
}
