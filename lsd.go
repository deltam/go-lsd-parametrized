/*
Package lsdp provides Weighted Levenshtein distance and its extended interface
*/
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
	result := accumulateCost(a, b, func(_, _ int, ar, br rune, diagonal, above, left float64) (float64, float64, float64) {
		if ar != br {
			diagonal += w.Replace
		}
		above += w.Insert
		left += w.Delete
		return diagonal, above, left
	}, minCost)
	return result
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
	ret := accumulateCost(a, b, func(_, _ int, ar, br rune, diagonal, above, left float64) (float64, float64, float64) {
		if rw, ok := wr.repRune[[2]rune{ar, br}]; ok {
			diagonal += rw
		} else if ar != br {
			diagonal += wr.w.Replace
		}
		if rw, ok := wr.insRune[br]; ok {
			above += rw
		} else {
			above += wr.w.Insert
		}
		if rw, ok := wr.delRune[ar]; ok {
			left += rw
		} else {
			left += wr.w.Delete
		}
		return diagonal, above, left
	}, minCost)
	return ret
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

type costFunc func(ai, bi int, ar, br rune, diagonal, above, left float64) (rep, ins, del float64)
type minFunc func(a, b, c float64) (min float64)

func accumulateCost(a, b string, costf costFunc, min minFunc) float64 {
	ar, br := []rune(a), []rune(b)
	costRow := make([]float64, len(ar)+1)
	for i := 1; i < len(costRow); i++ {
		_, _, costRow[i] = costf(i, 0, ar[i-1], 0, 0, 0, costRow[i-1])
	}

	var left float64
	for bc := 1; bc < len(br)+1; bc++ {
		_, left, _ = costf(0, bc, 0, br[bc-1], 0, costRow[0], 0)
		for i := 1; i < len(costRow); i++ {
			rep, ins, del := costf(i, bc, ar[i-1], br[bc-1], costRow[i-1], costRow[i], left)
			costRow[i-1] = left
			left = min(rep, ins, del)
		}
		costRow[len(costRow)-1] = left
	}

	return costRow[len(costRow)-1]
}

func minCost(a, b, c float64) (min float64) {
	min = a
	if b < min {
		min = b
	}
	if c < min {
		min = c
	}
	return
}
