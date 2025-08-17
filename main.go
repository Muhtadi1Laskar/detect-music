package main

import (
	"fmt"
	common "shazam/Common"
	database "shazam/Database"
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

	db := database.OpenBadger("./fingerprints.db")
	defer db.Close()

	for _, item := range fpsSong {
		database.AddFingerprint(
			db,
			item.Hash,
			database.DBEntry{
				SongID:     "song1",
				TimeOffset: item.TimeIndex, // ✅ use real offset
			},
		)
	}

	fmt.Println("Generated", len(fpsSong), "fingerprints")
}
