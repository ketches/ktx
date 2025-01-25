package util

import (
	"os"
)

func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func If[T any](cond bool, then T, els T) T {
	if cond {
		return then
	}
	return els
}
