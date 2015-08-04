package main

import "math"

// closestRatio returns whether you should pick a, if you wonder whether
// adding 1 to a or 1 to b will result in a a/(a+b) ratio closer to ratio.
func closestRatio(ratio, a, b float64) bool {
	aPlusOne := (a + 1) / (a + b + 1)
	bPlusOne := a / (a + b + 1)
	return math.Abs(aPlusOne-ratio) < math.Abs(bPlusOne-ratio)
}
