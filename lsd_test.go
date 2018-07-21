package lsdp

import (
	"math"
	"math/rand"
	"testing"
)

func equals(a, b float64) bool {
	epsilon := 0.000000000000001
	return math.Abs(a-b) < epsilon
}

func TestCountEdit(t *testing.T) {
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
			t.Errorf(`%d: lsd("%s", "%s") = %d, want %d`, i, d.A, d.B, c, d.Cost)
		}
		if !eqEdit(mt, d.Edit) {
			t.Errorf(`%d: lsd_edit("%s", "%s") = %v, want %v`, i, d.A, d.B, mt, d.Edit)
		}
	}
}

func TestWeights_Distance(t *testing.T) {
	testdata := []struct {
		W    Weights
		A    string
		B    string
		Cost float64
	}{
		{Weights{Insert: 1, Delete: 1, Replace: 1}, "", "", 0},
		{Weights{Insert: 1, Delete: 1, Replace: 1}, "", "a", 1},
		{Weights{Insert: 1, Delete: 1, Replace: 1}, "a", "", 1},
		{Weights{Insert: 1, Delete: 1, Replace: 1}, "a", "a", 0},
		{Weights{Insert: 1, Delete: 1, Replace: 1}, "back", "books", 3},
		{Weights{Insert: 1, Delete: 1, Replace: 0}, "back", "books", 1},
		{Weights{Insert: 1, Delete: 0, Replace: 1}, "back", "books", 3},
		{Weights{Insert: 0, Delete: 1, Replace: 1}, "back", "books", 2},
		{Weights{Insert: 1, Delete: 0, Replace: 1}, "back", "boo", 2},
	}

	for i, d := range testdata {
		if c := d.W.Distance(d.A, d.B); !equals(c, d.Cost) {
			t.Errorf(`%d: weighted_lsd("%s", "%s") = %f, want %f`, i, d.A, d.B, c, d.Cost)
		}
	}
}

func TestWeightsByRune_Distance(t *testing.T) {
	std := Weights{1, 1, 1}
	wrIns := ByRune(&std).Insert("a", 0.1)
	wrDel := ByRune(&std).Delete("a", 0.01)
	wrRep := ByRune(&std).Replace("a", "b", 0.001)
	wrAll := ByRune(&std).Insert("a", 0.1).Delete("a", 0.01).Replace("a", "b", 0.001)
	testdata := []struct {
		WR   *WeightsByRune
		A    string
		B    string
		Dist float64
	}{
		{wrIns, "", "a", 0.1},
		{wrIns, "", "aa", 0.2},
		{wrDel, "a", "", 0.01},
		{wrDel, "aa", "", 0.02},
		{wrRep, "a", "b", 0.001},
		{wrRep, "aa", "bb", 0.002},
		{wrAll, "aabc", "bbcaa", 0.211},
	}
	for i, td := range testdata {
		if d := td.WR.Distance(td.A, td.B); !equals(d, td.Dist) {
			t.Errorf(`%d: wr.Distance("%s", "%s") is %f, want %f`, i, td.A, td.B, d, td.Dist)
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
		{LevenshteinParam{Insert: 1, Delete: 0, Replace: 0}, "", "a", 1},
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

func TestNormalized(t *testing.T) {
	testdata := []struct {
		A    string
		B    string
		Dist float64
	}{
		{"", "", 0},
		{"", "a", 1},
		{"a", "", 1},
		{"", "aa", 1},
		{"a", "b", 1},
		{"ab", "bb", 0.5},
		{"ab", "a", 0.5},
	}

	nd := Normalized(LevenshteinParam{Insert: 1, Delete: 1, Replace: 1})
	for _, td := range testdata {
		if d := nd.Distance(td.A, td.B); !equals(d, td.Dist) {
			t.Errorf(`dist("%s", "%s") = %f, want %f`, td.A, td.B, d, td.Dist)
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

func makeBenchInputLongInput() string {
	var alphaNum string
	for i := 0; i < 1<<10; i++ {
		alphaNum += "abcdefghijklmnopqrstuvwxyz0123456789"
	}
	runes := []rune(alphaNum)
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}

var benchLongInput = makeBenchInputLongInput()

func benchmarkLsd(b *testing.B, s string) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Lsd(s, benchLongInput)
		Lsd(benchLongInput, s)
	}
}

func BenchmarkLsd1(b *testing.B) { benchmarkLsd(b, "a") }
func BenchmarkLsd2(b *testing.B) { benchmarkLsd(b, "abababababababababababababababababababababab") }

func benchWeightsDist(b *testing.B, s string) {
	w := Weights{
		Insert:  rand.Float64(),
		Delete:  rand.Float64(),
		Replace: rand.Float64(),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Distance(s, benchLongInput)
		w.Distance(benchLongInput, s)
	}
}

func BenchmarkWeightsDistance1(b *testing.B) { benchWeightsDist(b, "a") }
func BenchmarkWeightsDistance2(b *testing.B) {
	benchWeightsDist(b, "abababababababababababababababababababababab")
}

func BenchmarkCountEdit1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CountEdit("a", benchLongInput)
	}
}
func BenchmarkCountEdit2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CountEdit("abababababababbaababababbababa", benchLongInput)
	}
}
