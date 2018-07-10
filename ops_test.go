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
