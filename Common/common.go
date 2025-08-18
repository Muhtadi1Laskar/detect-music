package common

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func GetFilePath(songName string) string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "../audio/"+songName)
}

func GetFullPath(path string) string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, path)
}

func GetFileNames() []string {
	var result []string
	// dirPath := "C:/Users/laska/OneDrive/Documents/Coding/Work/shazam-clone/audio"
	dirPath := GetFullPath("../audio/")
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() { // Check if it's a file
			result = append(result, entry.Name())
		}
	}
	return result
}
