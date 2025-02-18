package util

import (
	"os"
)

// IsFileExist checks if a file exists.
func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// If returns then if cond is true, otherwise els.
func If[T any](cond bool, then T, els T) T {
	if cond {
		return then
	}
	return els
}
