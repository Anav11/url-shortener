package util

import (
	"fmt"
	"strconv"
)

var db = make(map[string]string)

func UrlShortener(url string, host string) string {
	path :=  strconv.FormatInt(int64(len(db)), 10)
	shortUrl := fmt.Sprintf("%s/%s", host, path)
	db[shortUrl] = url

	return shortUrl
}