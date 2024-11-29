package helpers

import (
	"os"
	"strings"
)

func EnforceHTTPS(url string) string {
	if url[:4] != "http" {
		return "https" + url
	}
	return url

}

func IsRequestURLAllowed(url string) bool {
	if url == os.Getenv("domain") {
		return false
	}
	newUrl := strings.Replace(url, "https://", "", 1)			//:= is only used to declare and initialize a variable. If you want to reassign a value to a variable, you should use =. like in below cases.
	newUrl = strings.Replace(newUrl, "http://", "", 1)
	newUrl = strings.Replace(newUrl, "www.", "", 1)
	newUrl = strings.Split(newUrl, "/")[0]		//if the uri is like ->"localhost/abc" then it will split the uri and return "localhost" only.but since localhost is our domain name in this exmpale our fuinction will return the false value


	/*
	? Instead of below code you can also write "return newUrl != os.Getenv("domain")"
	if newUrl == os.Getenv("domain") {
		return false
	}
	return true
	*/
	return newUrl != os.Getenv("domain")

}
