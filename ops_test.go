package lsdp

import (
	"math/rand"
	"testing"
)

func TestNearest(t *testing.T) {
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
		ans, score := Nearest(param, td.Raw, answers)
		if ans != td.Answer {
			t.Errorf(`Nearest error: i=%d found "%s"(%f), want "%s"`, i, ans, score, td.Answer)
		}
	}
}

func TestDistanceAll(t *testing.T) {
	std := Weights{1, 1, 1}
	testdata := []struct {
		orig  string
		strs  []string
		dists []float64
	}{
		{"", []string{"a", "aa", "aaa"}, []float64{1, 2, 3}},
		{"aa", []string{"aa", "a"}, []float64{0, 1}},
		{"a", []string{"", "a", "ab"}, []float64{1, 0, 1}},
	}
	for i, td := range testdata {
		ds := DistanceAll(std, td.orig, td.strs)
		for j, _ := range ds {
			if td.dists[j] != ds[j] {
				t.Errorf("%d: DistanceAll() is %v, want %v", i, ds, td.dists)
				break
			}
		}
	}
}

func makeBenchInputStrings() []string {
	strs := make([]string, 1<<6)
	var alphaNum string
	for i := 0; i < 10; i++ {
		alphaNum += "abcdefghijklmnopqrstuvwxyz0123456789"
	}
	runes := []rune(alphaNum)
	for i := 0; i < len(strs); i++ {
		rand.Shuffle(len(runes), func(i, j int) {
			runes[i], runes[j] = runes[j], runes[i]
		})
		strs[i] = string(runes)
	}
	return strs
}

var benchInputStrings = makeBenchInputStrings()

func BenchmarkNearest1(b *testing.B) { benchNearest(b, "a") }
func BenchmarkNearest2(b *testing.B) { benchNearest(b, "aaaaaaaaaaaa000000000000000000") }

func benchNearest(b *testing.B, s string) {
	std := Weights{1, 1, 1}
	input := makeBenchInputStrings()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Nearest(std, s, input)
	}
}

func BenchmarkDistanceAll1(b *testing.B) { benchDistanceAll(b, "a") }
func BenchmarkDistanceAll2(b *testing.B) { benchDistanceAll(b, "aaaaaaaaaaaa000000000000000000") }

func benchDistanceAll(b *testing.B, s string) {
	std := Weights{1, 1, 1}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DistanceAll(std, s, benchInputStrings)
	}
}
