// Package lsdp is a Levenshtein distance and its extended interface
package lsdp

// DistanceMeasurer provides measurement of the distance between 2 strings
type DistanceMeasurer interface {
	Distance(string, string) float64
}

// Lsd returns standard Levenshtein distance
func Lsd(a, b string) int {
	wd := &Weights{1, 1, 1}
	return int(wd.Distance(a, b))
}

// Weights represents cost parameters for weighted Levenshtein distance
type Weights struct {
	Insert  float64
	Delete  float64
	Replace float64
}

// Distance returns weighted Levenshtein distance
func (w Weights) Distance(a, b string) float64 {
	result := accumulateCost(a, b, func(ar, br rune, diagonal, above, left editCell) (editCell, editCell, editCell) {
		if ar != br {
			diagonal.cost += w.Replace
		}
		above.cost += w.Insert
		left.cost += w.Delete
		return diagonal, above, left
	}, lessCost)
	return result.cost
}

// ByRune returns weighted levenshtein distance by rune
func ByRune(w *Weights) *WeightsByRune {
	return &WeightsByRune{
		w:       w,
		insRune: make(map[rune]float64),
		delRune: make(map[rune]float64),
		repRune: make(map[[2]rune]float64),
	}
}

// WeightsByRune represents weighted levenshtein distance by rune
type WeightsByRune struct {
	w       *Weights
	insRune map[rune]float64
	delRune map[rune]float64
	repRune map[[2]rune]float64
}

// Distance returns weighted levenshtein distance by rune
func (wr *WeightsByRune) Distance(a, b string) float64 {
	ret := accumulateCost(a, b, func(ar, br rune, diagonal, above, left editCell) (editCell, editCell, editCell) {
		if rw, ok := wr.repRune[[2]rune{ar, br}]; ok {
			diagonal.cost += rw
		} else if ar != br {
			diagonal.cost += wr.w.Replace
		}
		if rw, ok := wr.insRune[br]; ok {
			above.cost += rw
		} else {
			above.cost += wr.w.Insert
		}
		if rw, ok := wr.delRune[ar]; ok {
			left.cost += rw
		} else {
			left.cost += wr.w.Delete
		}
		return diagonal, above, left
	}, lessCost)
	return ret.cost
}

// Insert specify cost by insert rune
func (wr *WeightsByRune) Insert(runeGroup string, insCost float64) *WeightsByRune {
	for _, r := range runeGroup {
		wr.insRune[r] = insCost
	}
	return wr
}

// Delete specify cost by delete rune
func (wr *WeightsByRune) Delete(runeGroup string, delCost float64) *WeightsByRune {
	for _, r := range runeGroup {
		wr.delRune[r] = delCost
	}
	return wr
}

// Replace specify cost by replace rune
func (wr *WeightsByRune) Replace(runeGroupSrc, runeGroupDest string, repCost float64) *WeightsByRune {
	for _, rs := range runeGroupSrc {
		for _, rd := range runeGroupDest {
			wr.repRune[[2]rune{rs, rd}] = repCost
		}
	}
	return wr
}

// Normalized returns what wrapped the DistanceMeasurer with nomalize by string length
func Normalized(dm DistanceMeasurer) DistanceMeasurer {
	return normalizedParam{wrapped: dm}
}

type normalizedParam struct {
	wrapped DistanceMeasurer
}

func (p normalizedParam) Distance(a, b string) float64 {
	d := p.wrapped.Distance(a, b)
	l := len([]rune(a))
	if lb := len([]rune(b)); l < lb {
		l = lb
	}
	if l == 0 {
		return d
	}
	return d / float64(l)
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
	result := accumulateCost(a, b, func(aRune, bRune rune, diagonal, above, left editCell) (editCell, editCell, editCell) {
		if aRune != bRune {
			diagonal.inc(REPLACE)
		} else {
			diagonal.count[NONE]++
		}
		above.inc(INSERT)
		left.inc(DELETE)
		return diagonal, above, left
	}, func(ec1, ec2 editCell) bool {
		return ec1.cost-float64(ec1.count.Get(NONE)) < ec2.cost-float64(ec2.count.Get(NONE))
	})
	return int(result.cost), result.count
}

type costFunc func(ar, br rune, diagonal, above, left editCell) (rep, ins, del editCell)
type lessFunc func(a, b editCell) bool

func accumulateCost(a, b string, costf costFunc, less lessFunc) editCell {
	ar, br := []rune(a), []rune(b)
	costRow := make([]editCell, len(ar)+1)
	dummy := editCell{}
	for i := 1; i < len(costRow); i++ {
		_, _, costRow[i] = costf(ar[i-1], 0, dummy, dummy, costRow[i-1])
	}

	next := make([]editCell, len(costRow))
	for bc := 1; bc < len(br)+1; bc++ {
		_, next[0], _ = costf(0, br[bc-1], dummy, costRow[0], dummy)
		for i := 1; i < len(next); i++ {
			min, ins, del := costf(ar[i-1], br[bc-1], costRow[i-1], costRow[i], next[i-1])
			if less(ins, min) {
				min = ins
			}
			if less(del, min) {
				min = del
			}
			next[i] = min
		}
		costRow, next = next, costRow
	}

	return costRow[len(costRow)-1]
}

// editCell represents cost & number of edits
type editCell struct {
	cost  float64
	count EditCounts
}

func (c *editCell) inc(t EditType) {
	if t != NONE {
		c.cost++
	}
	c.count[t]++
}

func lessCost(a, b editCell) bool {
	return a.cost < b.cost
}
