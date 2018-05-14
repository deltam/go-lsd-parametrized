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

type EditType int

const (
	INSERT EditType = iota
	DELETE
	REPLACE
	NONE
)

type EditCounts [4]int

func (ec EditCounts) Get(t EditType) int {
	return ec[t]
}

// normal Levenshtein distance
func Lsd(a, b string) float64 {
	d, _ := DistanceWithDetail(a, b)
	return d
}

// weighted Levenshtein distance
func (p LevenshteinParam) Distance(a, b string) float64 {
	_, cnt := DistanceWithDetail(a, b)
	return cnt.weighted(p)
}

func DistanceWithDetail(a, b string) (float64, EditCounts) {
	ar, br := []rune(a), []rune(b)
	costRow := make([]editCell, len(ar)+1)
	for i := 1; i < len(costRow); i++ {
		costRow[i] = costRow[i-1]
		costRow[i].incIns()
	}

	next := make([]editCell, len(costRow))
	for bc := 1; bc < len(br)+1; bc++ {
		next[0] = costRow[0]
		next[0].incIns()
		for i := 1; i < len(next); i++ {
			next[i] = cost(ar[i-1], br[bc-1], costRow[i-1], costRow[i], next[i-1])
		}
		costRow, next = next, costRow
	}

	return costRow[len(costRow)-1].Cost, costRow[len(costRow)-1].Counts
}

func (p LevenshteinParam) FindNearest(raw string, subjects []string) (nearest string, distance float64) {
	type lsdResult struct {
		Str      string
		Distance float64
	}

	ch := make(chan lsdResult)
	for _, sub := range subjects {
		go func(s string) {
			d := p.Distance(raw, s)
			ch <- lsdResult{Str: s, Distance: d}
		}(sub)
	}

	initFlag := true
	for i := 0; i < len(subjects); i++ {
		result := <-ch
		if initFlag || result.Distance < distance {
			distance = result.Distance
			nearest = result.Str
			initFlag = false
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

/////////////////////////////////////////////////////////////////////////////////////
// private

const (
	insertCost  = 1
	deleteCost  = 1
	replaceCost = 1
)

func (ec *EditCounts) inc(t EditType) {
	ec[t]++
}

func (ec EditCounts) weighted(p LevenshteinParam) float64 {
	return float64(ec.Get(INSERT))*p.Insert + float64(ec.Get(DELETE))*p.Delete + float64(ec.Get(REPLACE))*p.Replace
}

type editCell struct {
	Cost   float64
	Counts EditCounts
}

func (c *editCell) incIns() {
	c.Cost += insertCost
	c.Counts.inc(INSERT)
}

func cost(aRune, bRune rune, diagonal, above, left editCell) editCell {
	rep := int(diagonal.Cost) - diagonal.Counts[NONE]
	if aRune != bRune {
		rep += replaceCost
	}
	ins := int(above.Cost) + insertCost - above.Counts[NONE]
	del := int(left.Cost) + deleteCost - left.Counts[NONE]

	var minCell editCell
	minCell = diagonal
	if aRune != bRune {
		minCell.Cost += replaceCost
		minCell.Counts.inc(REPLACE)
	} else {
		minCell.Counts.inc(NONE)
	}
	if ins < rep {
		minCell = above
		minCell.Cost += insertCost
		minCell.Counts.inc(INSERT)
	}
	if del < ins {
		minCell = left
		minCell.Cost += deleteCost
		minCell.Counts.inc(DELETE)
	}

	return minCell
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
