package lsd_parametrized

import "testing"

func TestGenerateLsdFunc(t *testing.T) {
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
		{LevenshteinParam{Insert: 1, Delete: 1, Replace: 1}, "こんにちは", "こんばんは", 2},
	}

	for _, d := range testdata {
		lsd := GenerateLsdFunc(d.Param)
		if c := lsd(d.A, d.B); c != d.Cost {
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
	}

	for _, d := range testdata {
		if c := Lsd(d.A, d.B); c != d.Cost {
			t.Errorf("lsd(\"%s\", \"%s\") = %f, want 2", d.A, d.B, d.Cost)
		}
	}
}
