package lsd_parametrized

type LevenshteinParam struct {
	Insert  float64
	Delete  float64
	Replace float64
}

type LsdFunc func(string, string) float64

func generateCostFunc(p LevenshteinParam) func(rune, rune, float64, float64, float64) float64 {
	return func(aRune, bRune rune, diagonal, above, left float64) float64 {
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
}

func GenerateLsdFunc(p LevenshteinParam) LsdFunc {
	costFunc := generateCostFunc(p)
	return func(a, b string) float64 {
		ar, br := []rune(a), []rune(b)
		costRow := make([]float64, len(ar)+1)
		for i, _ := range costRow {
			costRow[i] = float64(i)
		}

		next := make([]float64, len(costRow))
		for bc := 1; bc < len(br)+1; bc++ {
			next[0] = float64(bc)
			for i := 1; i < len(next); i++ {
				next[i] = costFunc(ar[i-1], br[bc-1], costRow[i-1], costRow[i], next[i-1])
			}
			costRow, next = next, costRow
		}

		return costRow[len(costRow)-1]
	}
}

func Lsd(a, b string) float64 {
	f := GenerateLsdFunc(LevenshteinParam{1, 1, 1})
	return f(a, b)
}
