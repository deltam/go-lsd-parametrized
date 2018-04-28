package lsd_parametrized

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
