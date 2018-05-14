package tools

import (
	"testing"

	. "github.com/deltam/go-lsd-parametrized"
)

func TestLevesteinParam_Evaluate(t *testing.T) {
	param := LevenshteinParam{Insert: 1, Delete: 1, Replace: 1}

	findStrs := []string{
		"book",
		"back",
		"cook",
	}
	evalCases := make(map[string]string)
	evalCases["book"] = "book"
	evalCases["back"] = "back"
	evalCases["cook"] = "cook"
	evalCases["backs"] = "cook" // error case

	rate, reports := Evaluate(param, findStrs, evalCases)
	if rate != 3.0/4.0 {
		t.Errorf("rate == %f, want %f", rate, 3.0/4.0)
	}
	if len(reports) != 1 {
		t.Errorf("fail report is %d items, want 1 items\nreport = %v", len(reports), reports)
	}
	if reports[0].Raw != "backs" {
		t.Errorf("fail str == %s, want 'backs'", reports[0].Raw)
	}
}

func TestLevesteinParam_EvaluateByCSV(t *testing.T) {
	param := LevenshteinParam{Insert: 1, Delete: 1, Replace: 1}

	rate, reports, err := EvaluateByCSV(param, "testdata/pattern.csv", "")
	if err != nil {
		t.Errorf("err %v", err)
	}
	if rate != 3.0/4.0 {
		t.Errorf("rate == %f, want %f", rate, 3.0/4.0)
	}
	if len(reports) != 1 {
		t.Errorf("fail report is %d items, want 1 items\nreport = %v", len(reports), reports)
	}
	if reports[0].Raw != "backs" {
		t.Errorf("fail str == %s, want 'backs'", reports[0].Raw)
	}

	rate, reports, err = EvaluateByCSV(param, "testdata/pattern.csv", "testdata/findstrs.csv")
	if err != nil {
		t.Errorf("err %v", err)
	}
	if rate != 3.0/4.0 {
		t.Errorf("rate == %f, want %f", rate, 3.0/4.0)
	}
	if len(reports) != 1 {
		t.Errorf("fail report is %d items, want 1 items\nreport = %v", len(reports), reports)
	}
	if reports[0].Raw != "backs" {
		t.Errorf("fail str == %s, want 'backs'", reports[0].Raw)
	}
}
