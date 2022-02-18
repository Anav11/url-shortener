package util

import "fmt"

func GetUrl(url string, host string) string{
	_url := fmt.Sprintf("%s%s", host, url)
	return db[_url]
}
