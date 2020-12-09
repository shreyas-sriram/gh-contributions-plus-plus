package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDateData(t *testing.T) {
	type input struct {
		rawHTML string
		year    string
	}

	tests := []struct {
		name string
		args input
		want []ContributionEntry
	}{
		{name: "simple match", args: input{"data-count=\"10\" data-date=\"2020-11-22\"", "2020"}, want: []ContributionEntry{ContributionEntry{"2020-11-22", 10}}},
		{name: "match html string", args: input{"<foobar>data-count=\"10\" data-date=\"2020-11-22\"</foobar>", "2020"}, want: []ContributionEntry{ContributionEntry{"2020-11-22", 10}}},
		{name: "no matches", args: input{"<foobar></foobar>", "2020"}, want: []ContributionEntry{}},
	}

	for _, test := range tests {
		got := parseDateData(test.args.rawHTML, test.args.year)
		assert.Equal(t, got, test.want, "got: %+v, want %+v", got, test.want)
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
		assert.Equal(t, got, test.want, "got: %+v, want %+v", got, test.want)
	}
}
