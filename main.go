package main

import (
	"fmt"
	common "shazam/Common"
	database "shazam/Database"
	"shazam/transformation"

	"github.com/dgraph-io/badger/v4"
)


func matchFingerprints(db *badger.DB, queryFingerprints []transformation.FingerPrints) {
	score := make(map[string]int)

	for _, q := range queryFingerprints {
		if entries, err := database.LookupFingerprint(db, q.Hash); err == nil {
			for _, e := range entries {
				key := fmt.Sprintf("%s_%d", e.SongID, e.TimeOffset-q.TimeIndex)
				score[key]++
			}
		}
	}

	bestKey := ""
	bestCount := 0
	for k, n := range score {
		if n > bestCount {
			bestCount = n
			bestKey = k
		}
	}
	fmt.Printf("\n\n\n")
	fmt.Println("Best match:", bestKey, "count =", bestCount)
	fmt.Printf("\n\n\n")
}

func SaveNewSong() {
	windowSize := 4096
	overlap := 0.5
	topNPeaks := 5
	fileName := common.GetFileNames()

	db := database.OpenBadger("./fingerprints.db")
	defer db.Close()

	for _, name := range fileName {
		fmt.Println("Processing the song: ", name)
		path := common.GetFilePath(name)
		peaksSong, _ := transformation.GetPeaksWindow(path, windowSize, overlap, topNPeaks)
		fpsSong := transformation.BuildFingerprints(peaksSong, name)


		for _, item := range fpsSong {
			database.AddFingerprint(
				db,
				item.Hash,
				database.DBEntry{
					SongID:     name,
					TimeOffset: item.TimeIndex, 
				},
			)
		}

		fmt.Printf("\n\n\n")
		fmt.Println("Generated", len(fpsSong), "fingerprints")
		fmt.Printf("\n\n\n")
	}
}

func MatchSongs(songName string) {
	windowSize := 4096
	overlap := 0.5
	topNPeaks := 5

	db := database.OpenBadger("./fingerprints.db")
	defer db.Close()

	path := common.GetFilePath(songName)
	peakSong, _ := transformation.GetPeaksWindow(path, windowSize, overlap, topNPeaks)
	fpsSong := transformation.BuildFingerprints(peakSong, songName)

	matchFingerprints(db, fpsSong)
}

func main() {
	MatchSongs("sample3.mp3")
	// SaveNewSong()
}

