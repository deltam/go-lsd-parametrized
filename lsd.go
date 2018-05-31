// Calculate Levestein Distance by specific parameters written in Go.
package lsd_parametrized

// DistanceMeasurer provides measurement of the distance between 2 strings
type DistanceMeasurer interface {
	Distance(string, string) float64
}

// LevenshteinParam represents normal & weighted Levenshtein distance parameters
type LevenshteinParam struct {
	Insert  float64
	Delete  float64
	Replace float64
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

// Lsd returns normal Levenshtein distance
func Lsd(a, b string) int {
	d, _ := CountEdit(a, b)
	return d
}

// Distance returns weighted Levenshtein distance
func (p LevenshteinParam) Distance(a, b string) float64 {
	_, cnt := CountEdit(a, b)
	return float64(cnt.Get(INSERT))*p.Insert + float64(cnt.Get(DELETE))*p.Delete + float64(cnt.Get(REPLACE))*p.Replace
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

	return costRow[len(costRow)-1].Cost, costRow[len(costRow)-1].Counts
}

// Get the number of specified edit
func (ec EditCounts) Get(t EditType) int {
	return ec[t]
}

// editCell represents cost & number of edits
type editCell struct {
	Cost   int
	Counts EditCounts
}

func (c *editCell) inc(t EditType) {
	if t != NONE {
		c.Cost++
	}
	c.Counts[t]++
}

// cost returns current cost & number of edits
func cost(aRune, bRune rune, diagonal, above, left editCell) editCell {
	ins := above.Cost + 1 - above.Counts.Get(NONE)
	del := left.Cost + 1 - left.Counts.Get(NONE)
	rep := diagonal.Cost - diagonal.Counts.Get(NONE)
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
