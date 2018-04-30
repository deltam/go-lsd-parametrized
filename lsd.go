package lsd_parametrized

import (
	"encoding/csv"
	"io"
	"os"
)

type LevenshteinParam struct {
	Insert  float64
	Delete  float64
	Replace float64
}

var NormalLSD = LevenshteinParam{Insert: 1, Delete: 1, Replace: 1}

func (p LevenshteinParam) cost(aRune, bRune rune, diagonal, above, left float64) float64 {
	cost := diagonal
	if aRune != bRune {
		cost += p.Replace
	}
	if c := above + p.Insert; c < cost {
		cost = c
	}
	if c := left + p.Delete; c < cost {
		cost = c
	}
	return cost
}

func (p LevenshteinParam) Distance(a, b string) float64 {
	ar, br := []rune(a), []rune(b)
	costRow := make([]float64, len(ar)+1)
	for i, _ := range costRow {
		costRow[i] = float64(i)
	}

	next := make([]float64, len(costRow))
	for bc := 1; bc < len(br)+1; bc++ {
		next[0] = float64(bc)
		for i := 1; i < len(next); i++ {
			next[i] = p.cost(ar[i-1], br[bc-1], costRow[i-1], costRow[i], next[i-1])
		}
		costRow, next = next, costRow
	}

	return costRow[len(costRow)-1]
}

func Lsd(a, b string) float64 {
	return NormalLSD.Distance(a, b)
}

type lsdResult struct {
	Str      string
	Distance float64
}

func (p LevenshteinParam) FindNearest(raw string, subjects []string) (nearest string, distance float64) {
	ch := make(chan lsdResult)
	for _, sub := range subjects {
		go func(s string) {
			d := p.Distance(raw, s)
			ch <- lsdResult{Str: s, Distance: d}
		}(sub)
	}

	distance = 1000000000
	for i := 0; i < len(subjects); i++ {
		result := <-ch
		if result.Distance < distance {
			distance = result.Distance
			nearest = result.Str
		}
	}
	return
}

/////////////////////////////////////////////////////////////////////////////////////
// evaluate parameters

type FailedReport struct {
	Raw        string
	FailedStr  string
	SucceedStr string
	Dist       float64
}

func (p LevenshteinParam) Evaluate(findStrs []string, collectCases map[string]string) (succeedRate float64, reports []FailedReport) {
	for s, succeedStr := range collectCases {
		ans, dist := p.FindNearest(s, findStrs)
		if ans != succeedStr {
			reports = append(reports, FailedReport{Raw: s, FailedStr: ans, SucceedStr: succeedStr, Dist: dist})
		}
	}
	succeedRate = 1.0 - float64(len(reports))/float64(len(collectCases))
	return
}

func csv2Records(filename string) (records [][]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		records = append(records, rec)
	}
	return records, nil
}

// patternCsvFilename:
// "some string1","pattern1"
// "some string2","pattern2"
// ...
//
// findStrCsvFilename:
// "pattern1"
// "pattern2"
// ...
func (p LevenshteinParam) EvaluateByCSV(patternCsvFilename string, findStrCsvFilename string) (float64, []FailedReport, error) {
	patternDict := make(map[string]string)
	records, err := csv2Records(patternCsvFilename)
	if err != nil {
		return 0.0, nil, err
	}
	for _, rec := range records {
		patternDict[rec[0]] = rec[1]
	}

	var findStrs []string
	if findStrCsvFilename != "" {
		records, err := csv2Records(findStrCsvFilename)
		if err != nil {
			return 0.0, nil, err
		}
		for _, rec := range records {
			findStrs = append(findStrs, rec[0])
		}
	} else {
		for _, s := range patternDict {
			findStrs = append(findStrs, s)
		}
	}

	rate, reports := p.Evaluate(findStrs, patternDict)
	return rate, reports, nil
}
