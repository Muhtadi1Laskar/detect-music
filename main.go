package main

import (
	"fmt"
	fourier "shazam/Fourier"
)

func main() {
	samples := []float64{1, 2, 1, 2, 0, 0, 0, 0}
	specturm := fourier.DFT(samples)

	for i, j := range(specturm) {
		fmt.Println(i, j)
	}
}