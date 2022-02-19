package util

import (
	"strconv"
)

var db = make(map[string]string)

func URLShortener(url string) string {
	shortPath := strconv.FormatInt(int64(len(db)), 10)
	db[shortPath] = url

	return shortPath
}