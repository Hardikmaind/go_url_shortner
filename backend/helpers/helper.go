package helpers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
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

// Below is the function to create the QR code
func CreateQRCode(url string) ([]byte, error) {
	// Generate the QR code
	qr, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	return qr, nil
}

// ? I CAN USE THIS FUNCTION TO RATE LIMIT THE API REQUESTS DIRECLTY IN THE SHORTEN.GO FILE. INSTEAD OF ALL THE CLUTTER CODE 
func RateLimit(ctx *fiber.Ctx, r *redis.Client, c context.Context, IP string) error {
	_, err := r.Get(c, IP).Result()
	if err == redis.Nil { //In Redis, when you attempt to GET a key that doesn't exist, it returns the error "redis.Nil", which is a special error used to indicate that the key was not found.
		//! Set rate limit and expiration for new IP.
		//HERE THE IP is the key and api_quota is the value(defined in .env which is 10. we decrement this value by 1 (--1) each time the user hits the endpoint)
		err := r.Set(c, IP, os.Getenv("api_quota"), 2*time.Minute).Err() // Set key with 1-minute expiration
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to set rate limit",
			})
		}
	} else {
		valInt, _ := r.Get(c, IP).Int()
		if valInt <= 0 {
			limit, _ := r.TTL(c, IP).Result() // Check TTL for the current IP

			if limit <= 0 { // If TTL has expired (or not set)
				return ctx.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error":      "Rate limit exceeded",
					"retryAfter": 1, // Retry after 1 minute
				})
			}

			// If TTL is still valid, return retry time in seconds
			retryAfter := fmt.Sprintf("%d seconds", int(limit.Seconds()))
			return ctx.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":      "Rate limit exceeded",
				"retryAfter": retryAfter,
			})
		}
	}
	return nil // Return nil if no rate limit is exceeded
}
