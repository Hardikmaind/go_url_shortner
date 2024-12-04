package helpers

import (
	"os"
	"strings"
	"time"

	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/exp/rand"
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
	newUrl := strings.Replace(url, "https://", "", 1) //:= is only used to declare and initialize a variable. If you want to reassign a value to a variable, you should use =. like in below cases.
	newUrl = strings.Replace(newUrl, "http://", "", 1)
	newUrl = strings.Replace(newUrl, "www.", "", 1)
	newUrl = strings.Split(newUrl, "/")[0] //if the uri is like ->"localhost/abc" then it will split the uri and return "localhost" only.but since localhost is our domain name in this exmpale our fuinction will return the false value

	/*
		? Instead of below code you can also write "return newUrl != os.Getenv("domain")"
		if newUrl == os.Getenv("domain") {
			return false
		}
		return true
	*/
	return newUrl != os.Getenv("domain")

}

//! Atual logic behind url shortening

// Base62 character set
const base62Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// GenerateRandomString generates a random Base62 string of a given length
func GenerateRandomString(length int) string {
	// Seed the random number generator to ensure randomness
	rand.Seed(uint64(time.Now().UnixNano()))

	// Create a byte slice to store the generated characters
	result := make([]byte, length)

	// Populate the result slice with random characters from the Base62 charset
	for i := 0; i < length; i++ {
		result[i] = base62Charset[rand.Intn(len(base62Charset))]
	}

	return string(result)
}

//Below is the function to create the QR code
func CreateQRCode(url string) ([]byte, error) {
	// Generate the QR code
	qr, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return qr, nil
}
