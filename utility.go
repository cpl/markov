package markov

func clamp(val, min, max int) int {
	if val > max {
		val = max
	} else if val < min {
		val = min
	}

	return val
}
