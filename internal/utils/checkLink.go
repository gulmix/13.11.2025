package utils

import "net/http"

func CheckLink(url string) string {
	if len(url) < 7 || url[:7] != "http://" && url[:8] != "https://" {
		url = "http://" + url
	}
	_, err := http.Get(url)
	if err != nil {
		return "not available"
	}
	return "available"
}
