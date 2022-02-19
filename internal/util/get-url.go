package util

import (
	"strings"
)

func GetURL(url string) string{
	url = strings.TrimPrefix(url, "/")
	return db[url]
}
