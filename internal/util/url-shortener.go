package util

import (
	"fmt"
	"strconv"
)

var db = make(map[string]string)

func URLShortener(url string, host string) string {
	path :=  strconv.FormatInt(int64(len(db)), 10)
	shortURL := fmt.Sprintf("%s/%s", host, path)
	db[shortURL] = url

	return shortURL
}