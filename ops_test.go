package lsdp

import "testing"

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

func BenchmarkNearest(b *testing.B) {
	std := Weights{1, 1, 1}
	s := "aaaaaaaaaaaa000000000000000000"
	group := make([]string, 100)
	for i := 0; i < len(group); i++ {
		group[i] = "abcdefghijk0000000000000000000"
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Nearest(std, s, group)
	}
}

func BenchmarkDistanceAll(b *testing.B) {
	std := Weights{1, 1, 1}
	s := "aaaaaaaaaaaabbbbbbbbbbb00000000000000"
	group := make([]string, 1000)
	for i := 0; i < len(group); i++ {
		group[i] = "abcdefghijk0000000000000000000000000000"
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DistanceAll(std, s, group)
	}
}
