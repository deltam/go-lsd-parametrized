package lsd_parametrized

import (
	"math"
	"testing"
)

func equals(a, b float64) bool {
	epsilon := 0.000000000000001
	return math.Abs(a-b) < epsilon
}

func TestLevesteinParam_Distance(t *testing.T) {
	testdata := []struct {
		Param LevenshteinParam
		A     string
		B     string
		Cost  float64
	}{
		{LevenshteinParam{Insert: 1, Delete: 1, Replace: 1}, "book", "back", 2},
		{LevenshteinParam{Insert: 1, Delete: 1, Replace: 0.1}, "book", "back", 0.2},
		{LevenshteinParam{Insert: 0.1, Delete: 1, Replace: 1}, "book", "back", 2},
		{LevenshteinParam{Insert: 1, Delete: 0.1, Replace: 1}, "book", "back", 2},
		{LevenshteinParam{Insert: 0.01, Delete: 1, Replace: 0.1}, "book", "backs", 0.21},
		{LevenshteinParam{Insert: 1, Delete: 1, Replace: 1}, "こんにちは", "こんばんは", 2},
		{LevenshteinParam{Insert: 1, Delete: 1, Replace: 0.1}, "こんにちは", "こんばんは", 0.2},
		{LevenshteinParam{Insert: 1, Delete: 1, Replace: 1}, "book", "board", 3},
		{LevenshteinParam{Insert: 0.1, Delete: 1, Replace: 1}, "book", "board", 2.1},
		{LevenshteinParam{Insert: 1, Delete: 1, Replace: 0.1}, "book", "board", 1.2},
	}

	for _, d := range testdata {
		if c := d.Param.Distance(d.A, d.B); !equals(c, d.Cost) {
			t.Errorf("lsd(\"%s\", \"%s\") = %f, want %f", d.A, d.B, c, d.Cost)
		}
	}
}

func TestLsd(t *testing.T) {
	testdata := []struct {
		A    string
		B    string
		Cost float64
	}{
		{"book", "back", 2},
		{"こんにちは", "こんばんは", 2},
		{"kitten", "sitting", 3},
	}

	for _, d := range testdata {
		if c := Lsd(d.A, d.B); !equals(c, d.Cost) {
			t.Errorf("lsd(\"%s\", \"%s\") = %f, want 2", d.A, d.B, d.Cost)
		}
	}
}

func TestLevesteinParam_FindNearest(t *testing.T) {
	answers := []string{
		"book",
		"back",
		"cook",
	}
	testdata := []struct {
		Raw    string
		Answer string
	}{
		{"book", "book"},
		{"pack", "back"},
		{"sick", "back"},
		{"cop", "cook"},
	}

	param := LevenshteinParam{1, 1, 1}

	for i, td := range testdata {
		ans, score := param.FindNearest(td.Raw, answers)
		if ans != td.Answer {
			t.Errorf(`FindNearest error: i=%d found "%s"(%f), want "%s"`, i, ans, score, td.Answer)
		}
	}
}
