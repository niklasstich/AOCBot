package handlers

import (
	"github.com/niklasstich/AOCBot/aoc"
	"testing"
)

func TestFormatDays(t *testing.T) {
	testCases := []struct {
		start, end int
		want string
	}{
		{1,2, " 1  2"},
		{1,24, " 1  2  3  4  5  6  7  8  9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24"},
		{1, 1, " 1"},
	}

	for _, tc := range testCases {
		got := formatDays(tc.start, tc.end)

		if tc.want != got {
			t.Fatalf("Wanted \"%s\", got \"%s\"", tc.want, got)
		}
	}
}

func TestFormatStars(t *testing.T) {
	testCases := []struct {
		m aoc.Member
		start, end int
		want string
	}{
		{m: aoc.Member{
			Name:               "test",
			LocalScore:         0,
			GlobalScore:        0,
			Stars:              0,
			CompletionDayLevel: map[string]map[string]interface{}{
				"1": {"1": nil, "2": nil},
				"2": {"1": nil, "2": nil},
				"3": {"1": nil, "2": nil},
				"4": {"1": nil, "2": nil},
			},
		}, start: 1, end: 4, want:"[*][*][*][*]"},
		{m: aoc.Member{
			Name:               "test",
			LocalScore:         0,
			GlobalScore:        0,
			Stars:              0,
			CompletionDayLevel: map[string]map[string]interface{}{
				"1": {"1": nil},
				"2": {"1": nil, "2": nil},
				"3": {"1": nil},
				"4": {"1": nil, "2": nil},
			},
		}, start: 1, end: 4, want:"(*)[*](*)[*]"},
		{m: aoc.Member{
			Name:               "test",
			LocalScore:         0,
			GlobalScore:        0,
			Stars:              0,
			CompletionDayLevel: map[string]map[string]interface{}{
				"1": {"1": nil},
				"3": {"1": nil},
				"4": {"1": nil, "2": nil},
			},
		}, start: 1, end: 4, want:"(*)   (*)[*]"},
		{m: aoc.Member{
			Name:               "test",
			LocalScore:         0,
			GlobalScore:        0,
			Stars:              0,
			CompletionDayLevel: map[string]map[string]interface{}{
				"1": {"1": nil},
				"2": {},
				"3": {"1": nil},
				"4": {"1": nil, "2": nil},
			},
		}, start: 1, end: 4, want:"(*)   (*)[*]"},
		{m: aoc.Member{
			Name:               "test",
			LocalScore:         0,
			GlobalScore:        0,
			Stars:              0,
			CompletionDayLevel: map[string]map[string]interface{}{
				"1": {"1": nil},
				"2": {},
				"3": {"1": nil},
				"4": {"1": nil, "2": nil},
			},
		}, start: 1, end: 3, want:"(*)   (*)"},
		{m: aoc.Member{
			Name:               "test",
			LocalScore:         0,
			GlobalScore:        0,
			Stars:              0,
			CompletionDayLevel: map[string]map[string]interface{}{
				"1": {"1": nil},
				"2": {},
				"3": {"1": nil},
				"4": {"1": nil, "2": nil},
			},
		}, start: 3, end: 4, want:"(*)[*]"},
		{m: aoc.Member{
			Name:               "test",
			LocalScore:         0,
			GlobalScore:        0,
			Stars:              0,
			CompletionDayLevel: map[string]map[string]interface{}{
				"1": {"1": nil},
				"2": {},
				"3": {"1": nil},
				"4": {"1": nil, "2": nil},
			},
		}, start: 2, end: 2, want:"   "},
	}

	for _, tc := range testCases {
		got := formatMemberStars(tc.m,tc.start,tc.end)

		if tc.want != got {
			t.Fatalf("Wanted \"%s\", got \"%s\"", tc.want, got)
		}
	}
}