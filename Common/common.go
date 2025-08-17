package common

import (
	"path/filepath"
	"runtime"
)

func GetFilePath(songName string) string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, "../audio/" + songName + ".wav")
}