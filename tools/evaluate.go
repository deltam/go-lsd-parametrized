package tools

import (
	"encoding/csv"
	"io"
	"os"

	. "github.com/deltam/go-lsd-parametrized"
)

type FailedReport struct {
	Raw        string
	FailedStr  string
	SucceedStr string
	Dist       float64
}

func Evaluate(dm DistanceMeasurer, findStrs []string, collectCases map[string]string) (succeedRate float64, reports []FailedReport) {
	for s, succeedStr := range collectCases {
		ans, dist := FindNearest(dm, s, findStrs)
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
func EvaluateByCSV(dm DistanceMeasurer, patternCsvFilename string, findStrCsvFilename string) (float64, []FailedReport, error) {
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

	rate, reports := Evaluate(dm, findStrs, patternDict)
	return rate, reports, nil
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
