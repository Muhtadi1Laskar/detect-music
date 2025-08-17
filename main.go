package main

import (
	"fmt"
	"log"
	"os"
	fourier "shazam/Fourier"

	"github.com/faiface/beep/wav"
)

type Peak struct {
	Freq float64
	Mag float64
}

func getPeakWindow(filePath string, windowSize int, overlapRatio float64, topN int) ([][]Peak, int) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	var samples []float64
	buf := make([][2]float64, 1)
	for {
		sample, ok := streamer.Stream(buf)
		if !ok || sample == 0 {
			break
		}
		samples = append(samples, buf[0][0])
	}

	step := int(float64(windowSize) * (1 - overlapRatio))
	var peakPerWindow [][]Peak
	
	for start := 0; start+windowSize <= len(samples); start += step {
		window := samples[start : start+windowSize]
		coeffs := fourier.DFT(window)
		peaks := extractTopNPeaks(coeffs, format.SampleRate.N(1), windowSize, topN)
		peakPerWindow = append(peakPerWindow, peaks)
	}

	return peakPerWindow, step
}

func main() {
	samples := []float64{1, 2, 1, 2, 0, 0, 0, 0}
	specturm := fourier.DFT(samples)

	for i, j := range(specturm) {
		fmt.Println(i, j)
	}
}