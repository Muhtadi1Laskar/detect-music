package main

import (
	"fmt"
	"log"
	"math"
	"os"
	fourier "shazam/Fourier"
	"sort"

	"github.com/faiface/beep/wav"
)

type Peak struct {
	Freq float64
	Mag  float64
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
		peaks := extractTopNPeaks(coeffs, float64(format.SampleRate.N(1)), windowSize, topN)
		peakPerWindow = append(peakPerWindow, peaks)
	}

	return peakPerWindow, step
}

func extractTopNPeaks(coeffs []complex128, sampleRate float64, windowSize int, topN int) []Peak {
	peaks := make([]Peak, 0, len(coeffs))
	for i := range coeffs {
		re := real(coeffs[i])
		im := imag(coeffs[i])
		mag := math.Sqrt(re*re + im*im)
		freq := float64(i) * sampleRate / float64(windowSize)
		peaks = append(peaks, Peak{Freq: freq, Mag: mag})
	}

	sort.Slice(peaks, func(i, j int) bool {
		return peaks[i].Mag > peaks[j].Mag
	})

	if len(peaks) > topN {
		return peaks[:topN]
	}
	return peaks
}

func main() {
	samples := []float64{1, 2, 1, 2, 0, 0, 0, 0}
	specturm := fourier.DFT(samples)

	for i, j := range specturm {
		fmt.Println(i, j)
	}
}
