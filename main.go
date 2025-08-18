package main

import (
	"fmt"
	common "shazam/Common"
	database "shazam/Database"
	"shazam/transformation"
)


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

	transformation.MatchFingerprints(db, fpsSong)
}


func main() {
	MatchSongs("sample3.mp3")
	fmt.Println(common.GetFileNames())
	// SaveNewSong()
}

