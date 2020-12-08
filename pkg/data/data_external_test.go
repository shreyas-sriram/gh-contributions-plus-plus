package data_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shreyas-sriram/gh-contributions-aggregator/pkg/data"
)

func TestParseDateData(t *testing.T) {
	type contribution struct {
		date string
		data int
	}

	type want struct {
		total int
		data  []contribution
	}

	type args struct {
		rawHTML string
		year    string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "simple match", args: args{"10 contributions this year blah blah blah data-count=\"10\" data-date=\"2020-11-22\"", "2020"}, want: want{10, []contribution{contribution{"2020-11-22", 10}}}},
		{name: "multiple match", args: args{"14 contributions this year blah blah blah data-count=\"6\" data-date=\"2020-11-22\"\ndata-count=\"8\" data-date=\"2020-11-23\"", "2020"}, want: want{14, []contribution{contribution{"2020-11-22", 6}, contribution{"2020-11-23", 8}}}},
		{name: "no match", args: args{"blah blah blah", "2020"}, want: want{}},
	}

	for _, test := range tests {
		got := data.ParseContributionsData(test.args.rawHTML, test.args.year)
		if got.Total != test.want.total {
			t.Errorf("Got and want were incorrect, got: %+vs, want: %+v", got.Total, test.want.total)
		}

		for i, gotData := range got.ContributionData {
			if gotData.Date != test.want.data[i].date {
				t.Errorf("Got and want were incorrect, got: %+vs, want: %+v", gotData.Date, test.want.data[i].date)
			}
			if gotData.Data != test.want.data[i].data {
				t.Errorf("Got and want were incorrect, got: %+vs, want: %+v", gotData.Data, test.want.data[i].data)
			}
		}
	}
}

func TestGetRawPage(t *testing.T) {
	type args struct {
		username string
		year     string
	}

	type want struct {
		rawHTML string
		err     error
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "invalid user", args: args{"so-mnbvcxz-random", "2020"}, want: want{"", fmt.Errorf("user not found")}},
		{name: "random user", args: args{"random", "2020"}, want: want{"contributions", nil}},
	}

	for _, test := range tests {
		got, err := data.GetRawPage(test.args.username, test.args.year)
		if err != nil {
			if got != test.want.rawHTML {
				t.Errorf("Got and want were incorrect, got: %+vs, want: %+v", got, test.want.rawHTML)
			}
		}
		res := strings.Contains(got, test.want.rawHTML)
		if !res {
			t.Errorf("Got and want were incorrect, got: %+vs, want: %+v", got, test.want.rawHTML)
		}
	}
}
