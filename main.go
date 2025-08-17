package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"os"
	common "shazam/Common"
	"sort"

	"github.com/faiface/beep/wav"
	"gonum.org/v1/gonum/dsp/fourier"
)

type Peak struct {
	Freq float64
	Mag  float64
}

type FingerPrints struct {
	Hash      string
	TimeIndex int
}

func getPeaksWindow(filePath string, windowSize int, overlapRatio float64, topN int) ([][]Peak, int) {
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
		samples = append(samples, float64(buf[0][0]))
	}

	step := int(float64(windowSize) * (1 - overlapRatio))
	var peakPerWindow [][]Peak
	fft := fourier.NewFFT(windowSize)

	for start := 0; start+windowSize <= len(samples); start += step {
		window := samples[start : start+windowSize]

		if len(window) != windowSize {
			continue
		}

		// Compute FFT coefficients
		coeffs := fft.Coefficients(nil, window)

		peaks := extractTopNPeaks(coeffs, float64(format.SampleRate), windowSize, topN)
		peakPerWindow = append(peakPerWindow, peaks)
	}

	return peakPerWindow, step
}

func extractTopNPeaks(coeffs []complex128, sampleRate float64, windowSize int, topN int) []Peak {
	peaks := make([]Peak, 0, len(coeffs)/2) // Only need first half (real FFT symmetry)

	for i := 1; i < len(coeffs)/2; i++ { // Skip DC component (i=0)
		re := real(coeffs[i])
		im := imag(coeffs[i])
		// mag := math.Sqrt(re*re + im*im)
		mag := 10 * math.Log10(re*re+im*im+1e-12)
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

func buildFingerprints(peaksPerWindow [][]Peak, songID string) []FingerPrints {
	var fps []FingerPrints
	for winIdx, peaks := range peaksPerWindow {
		for i := 0; i < len(peaks); i++ {
			for j := i + 1; j < len(peaks); j++ {
				key := fmt.Sprintf("%d|%d|0", int(peaks[i].Freq), int(peaks[j].Freq))
				h := sha1.Sum([]byte(key))
				fp := FingerPrints{
					Hash:      hex.EncodeToString(h[:]),
					TimeIndex: winIdx,
				}
				fps = append(fps, fp)
			}
		}
	}
	return fps
}

func main() {
	windowSize := 4096
	overlap := 0.5
	topNPeaks := 5
	songName := "sample"
	path := common.GetFilePath(songName) // Replace with actual file path

	peaksSong, _ := getPeaksWindow(path, windowSize, overlap, topNPeaks)
	fpsSong := buildFingerprints(peaksSong, "song1")

	fmt.Println("Generated", len(fpsSong), "fingerprints")
	// fmt.Println(fpsSong)
}
