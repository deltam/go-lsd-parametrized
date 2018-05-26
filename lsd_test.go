package lsd_parametrized

import (
	"math"
	"testing"
)

func equals(a, b float64) bool {
	epsilon := 0.000000000000001
	return math.Abs(a-b) < epsilon
}

func TestLevesteinParam_CountEdit(t *testing.T) {
	testdata := []struct {
		A    string
		B    string
		Cost int
		Edit [4]int
	}{
		{"a", "aaa", 2, EditCounts{2, 0, 0, 1}},
		{"aaa", "a", 2, EditCounts{0, 2, 0, 1}},
		{"aaaaa", "aaa", 2, EditCounts{0, 2, 0, 3}},
		{"book", "back", 2, EditCounts{0, 0, 2, 2}},
		{"book", "backs", 3, EditCounts{1, 0, 2, 2}},
		{"こんにちは", "こんばんは", 2, EditCounts{0, 0, 2, 3}},
		{"book", "board", 3, EditCounts{1, 0, 2, 2}},
		{"book", "boo", 1, EditCounts{0, 1, 0, 3}},
	}
	eqEdit := func(et EditCounts, edit EditCounts) bool {
		for i := 0; i < 4; i++ {
			if et[i] != edit[i] {
				return false
			}
		}
		return true
	}

	for i, d := range testdata {
		c, mt := CountEdit(d.A, d.B)
		if c != d.Cost {
			t.Errorf("%d: lsd(\"%s\", \"%s\") = %d, want %d", i, d.A, d.B, c, d.Cost)
		}
		if !eqEdit(mt, d.Edit) {
			t.Errorf("%d: lsd_edit(\"%s\", \"%s\") = %v, want %v", i, d.A, d.B, mt, d.Edit)
		}
	}
}

func TestLevesteinParam_Distance(t *testing.T) {
	testdata := []struct {
		Param LevenshteinParam
		A     string
		B     string
		Cost  float64
	}{
		{LevenshteinParam{Insert: 1, Delete: 0, Replace: 0}, "a", "aaa", 2},
		{LevenshteinParam{Insert: 0, Delete: 1, Replace: 0}, "aaa", "a", 2},
		{LevenshteinParam{Insert: 0, Delete: 0, Replace: 1}, "aaa", "abc", 2},
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

	for i, d := range testdata {
		if c := d.Param.Distance(d.A, d.B); !equals(c, d.Cost) {
			t.Errorf("%d: lsd(\"%s\", \"%s\") = %f, want %f", i, d.A, d.B, c, d.Cost)
		}
	}
}

func TestLsd(t *testing.T) {
	testdata := []struct {
		A    string
		B    string
		Cost int
	}{
		{"book", "back", 2},
		{"こんにちは", "こんばんは", 2},
		{"kitten", "sitting", 3},
	}

	for _, d := range testdata {
		if c := Lsd(d.A, d.B); c != d.Cost {
			t.Errorf("lsd(\"%s\", \"%s\") = %d, want 2", d.A, d.B, d.Cost)
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
		ans, score := FindNearest(param, td.Raw, answers)
		if ans != td.Answer {
			t.Errorf(`FindNearest error: i=%d found "%s"(%f), want "%s"`, i, ans, score, td.Answer)
		}
	}
}
