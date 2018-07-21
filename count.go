package lsdp

// EditType represents authorized editing means in Levenshtein distance
type EditType int

// Authorized editing means: insert, delete, replace, none
const (
	INSERT EditType = iota
	DELETE
	REPLACE
	NONE
)

// EditCounts represents aggregating by editing types
type EditCounts [4]int

// Get the number of specified edit
func (ec EditCounts) Get(t EditType) int {
	return ec[t]
}

// CountEdit aggregates the minimum number of edits to change from a to b
func CountEdit(a, b string) (int, EditCounts) {
	cnts := make([]EditCounts, len([]rune(a))+1)
	var leftCnt, repCnt, insCnt, delCnt EditCounts
	var cur int
	result := accumulateCost(a, b, func(ai, bi int, ar, br rune, diagonal, above, left float64) (float64, float64, float64) {
		if ai == 0 {
			above++
			leftCnt = cnts[0]
			leftCnt[INSERT]++
			if len(cnts) == 1 {
				cnts[0] = leftCnt
			}
			return 0, above, 0
		} else if bi == 0 {
			left++
			cnts[ai] = cnts[ai-1]
			cnts[ai][DELETE]++
			return 0, 0, left
		}
		cur = ai
		repCnt = cnts[ai-1]
		insCnt = cnts[ai]
		delCnt = leftCnt
		if ar != br {
			diagonal++
			repCnt[REPLACE]++
		} else {
			repCnt[NONE]++
		}
		above++
		insCnt[INSERT]++
		left++
		delCnt[DELETE]++
		return diagonal, above, left
	}, func(rep, ins, del float64) (min float64) {
		min = rep
		minCnt := repCnt
		if ins-float64(insCnt.Get(NONE)) < min-float64(minCnt.Get(NONE)) {
			min = ins
			minCnt = insCnt
		}
		if del-float64(delCnt.Get(NONE)) < min-float64(minCnt.Get(NONE)) {
			min = del
			minCnt = delCnt
		}
		cnts[cur-1] = leftCnt
		leftCnt = minCnt
		if cur == len(cnts)-1 {
			cnts[cur] = leftCnt
		}
		return
	})
	return int(result), cnts[len(cnts)-1]
}

// LevenshteinParam represents Levenshtein distance parameters for weighted by edit counts
type LevenshteinParam struct {
	Insert  float64
	Delete  float64
	Replace float64
}

// Distance returns Levenshtein distance
func (p LevenshteinParam) Distance(a, b string) float64 {
	_, cnt := CountEdit(a, b)
	return float64(cnt.Get(INSERT))*p.Insert + float64(cnt.Get(DELETE))*p.Delete + float64(cnt.Get(REPLACE))*p.Replace
}
