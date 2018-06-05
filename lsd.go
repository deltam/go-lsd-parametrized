// Package lsdp is a Levenshtein distance and its extended interface
package lsdp

// DistanceMeasurer provides measurement of the distance between 2 strings
type DistanceMeasurer interface {
	Distance(string, string) float64
}

// Lsd returns standard Levenshtein distance
func Lsd(a, b string) int {
	d, _ := CountEdit(a, b)
	return d
}

// LevenshteinParam represents normal & weighted Levenshtein distance parameters
type LevenshteinParam struct {
	Insert  float64
	Delete  float64
	Replace float64
}

// Distance returns weighted Levenshtein distance
func (p LevenshteinParam) Distance(a, b string) float64 {
	_, cnt := CountEdit(a, b)
	return float64(cnt.Get(INSERT))*p.Insert + float64(cnt.Get(DELETE))*p.Delete + float64(cnt.Get(REPLACE))*p.Replace
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

// Nearest returns the nearest string in the specified distance measurer
func Nearest(dm DistanceMeasurer, raw string, subjects []string) (nearest string, distance float64) {
	type result struct {
		str  string
		dist float64
	}

	ch := make(chan result)
	for _, sub := range subjects {
		go func(s string) {
			d := dm.Distance(raw, s)
			ch <- result{s, d}
		}(sub)
	}

	initFlag := true
	for i := 0; i < len(subjects); i++ {
		r := <-ch
		if initFlag || r.dist < distance {
			distance = r.dist
			nearest = r.str
			initFlag = false
		}
	}
	return
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
	ar, br := []rune(a), []rune(b)
	costRow := make([]editCell, len(ar)+1)
	for i := 1; i < len(costRow); i++ {
		costRow[i] = costRow[i-1]
		costRow[i].inc(INSERT)
	}

	next := make([]editCell, len(costRow))
	for bc := 1; bc < len(br)+1; bc++ {
		next[0] = costRow[0]
		next[0].inc(INSERT)
		for i := 1; i < len(next); i++ {
			next[i] = cost(ar[i-1], br[bc-1], costRow[i-1], costRow[i], next[i-1])
		}
		costRow, next = next, costRow
	}

	return costRow[len(costRow)-1].cost, costRow[len(costRow)-1].count
}

// editCell represents cost & number of edits
type editCell struct {
	cost  int
	count EditCounts
}

func (c *editCell) inc(t EditType) {
	if t != NONE {
		c.cost++
	}
	c.count[t]++
}

// cost returns current cost & number of edits
func cost(aRune, bRune rune, diagonal, above, left editCell) editCell {
	ins := above.cost + 1 - above.count.Get(NONE)
	del := left.cost + 1 - left.count.Get(NONE)
	rep := diagonal.cost - diagonal.count.Get(NONE)
	minEdit := NONE
	if aRune != bRune {
		rep++
		minEdit = REPLACE
	}

	minCell := diagonal
	if ins < rep {
		minCell = above
		minEdit = INSERT
	}
	if del < ins {
		minCell = left
		minEdit = DELETE
	}

	minCell.inc(minEdit)
	return minCell
}
