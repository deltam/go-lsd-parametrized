package lsdp

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

// DistanceAll returns slice of distance orig to each strs
func DistanceAll(dm DistanceMeasurer, orig string, strs []string) []float64 {
	dists := make([]float64, len(strs))
	done := make(chan struct{}, len(strs))
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
