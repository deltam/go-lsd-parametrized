package lsdp

// Nearest returns the nearest string in the specified distance measurer
func Nearest(dm DistanceMeasurer, orig string, strs []string) (nearest string, distance float64) {
	type result struct {
		str  string
		dist float64
	}

	ch := make(chan result)
	for _, s := range strs {
		go func(s string) {
			ch <- result{s, dm.Distance(orig, s)}
		}(s)
	}

	initFlag := true
	for range strs {
		r := <-ch
		if initFlag || r.dist < distance {
			distance = r.dist
			nearest = r.str
			initFlag = false
		}
	}
	return
}

// DistanceAll returns slice of distance orig to each strs
func DistanceAll(dm DistanceMeasurer, orig string, strs []string) []float64 {
	dists := make([]float64, len(strs))
	done := make(chan struct{})
	for i, s := range strs {
		go func(i int, s string) {
			dists[i] = dm.Distance(orig, s)
			done <- struct{}{}
		}(i, s)
	}
	for range strs {
		<-done
	}
	return dists
}
