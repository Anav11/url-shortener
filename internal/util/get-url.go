package util

import "fmt"

func GetURL(url string, host string) string{
	_url := fmt.Sprintf("%s%s", host, url)
	return db[_url]
}
