package main

import (
	"fmt"
	common "shazam/Common"
	"shazam/transformation"
)



func main() {
	windowSize := 4096
	overlap := 0.5
	topNPeaks := 5
	songName := "sample"
	path := common.GetFilePath(songName)

	peaksSong, _ := transformation.GetPeaksWindow(path, windowSize, overlap, topNPeaks)
	fpsSong := transformation.BuildFingerprints(peaksSong, "song1")

	fmt.Println("Generated", len(fpsSong), "fingerprints")
}