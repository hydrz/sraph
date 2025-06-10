package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	sa := []float64{
		0,
		1,
		-1,
		-0,
		1e-3,
		1e+3,
		1e-6,
		1e+6,
		1.5,
		-1.5,
		math.MaxFloat64,
		-math.MaxFloat64,
		math.Inf(1),
		math.Inf(-1),
		math.NaN(),
	}
	for _, s := range sa {
		fmt.Printf("%s\n", strconv.FormatFloat(s, 'g', -1, 64))

		// fmt.Printf("%.3f\n", s)
	}

}
