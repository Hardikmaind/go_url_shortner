package routes

import (
	"fmt"
	"os"
	"time"

	"github.com/Hardikmaind/go_url_shortner/db"
	"github.com/Hardikmaind/go_url_shortner/helpers"
	"github.com/Hardikmaind/go_url_shortner/types"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	
)

// ? below is the handler which will be called when the user hits the /api/v1 endpoint
func UrlToQrcode(c *fiber.Ctx) error {
	var reqBody types.Request
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	//!WE CAN ALSO USE THE SAME SAME DB FOR THE QR CODE STORING AND URL SHORTENING. BUT FOR SEPARACTION OF CONCERNS WE WILL USE A DIFFERENT DB FOR QR CODE STORING.also diff client. ALSO WE WILL USE A DIFFERENT go context package for each client
	r2 := db.CreateClient2

	//TODO IMPLEMENT RATE LIMITING

	_, err := r2.Get(db.Ctx2, c.IP()).Result()
	if err == redis.Nil { //In Redis, when you attempt to GET a key that doesn't exist, it returns the error "redis.Nil", which is a special error used to indicate that the key was not found.
		//! Set rate limit and expiration for new IP.
		//HERE THE c.IP() is the key and api_quota is the value(defined in .env which is 10. we decrement this value by 1 (--1) each time the user hits the endpoint)
		err := r2.Set(db.Ctx2, c.IP(), os.Getenv("api_quota"), 2*time.Minute).Err() // Set key with 1-minute expiration
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to set rate limit",
			})
		}
	} else {
		valInt, _ := r2.Get(db.Ctx2, c.IP()).Int()
		if valInt <= 0 {
			limit, _ := r2.TTL(db.Ctx2, c.IP()).Result() // Check TTL for the current IP

			if limit <= 0 { // If TTL has expired (or not set)
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"error":      "Rate limit exceeded",
					"retryAfter": 1, // Retry after 1 minute
				})
			}

			// If TTL is still valid, return retry time in seconds
			retryAfter := fmt.Sprintf("%d seconds", int(limit.Seconds()))
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":      "Rate limit exceeded",
				"retryAfter": retryAfter,
			})
		}
	}





		//TODO if the key exists then we will return the value of the key(here that is the imaage(png)) else we will create the qr code and store it in the redis db and then return the value of the key(image) 

	if exists, _ := r2.Exists(db.Ctx2, reqBody.URL).Result(); exists > 0 {
		png, _ := r2.Get(db.Ctx2, reqBody.URL).Bytes()

		//? here we are decrementing the ip counter for rate limiting if the key exists in the redis db
		remainingQuota, _ := r2.Decr(db.Ctx2, c.IP()).Result() //decrement the count of the ip address by 1
		if remainingQuota < 0 {
			ttl, _ := r2.TTL(db.Ctx2, c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":      "Rate limit exceeded",
				"retryAfter": ttl / time.Second,
			})
		}


		ttl, err := r2.TTL(db.Ctx2, reqBody.URL).Result() //? Fixed the expiry..now expiry is consistent, in cache hit or miss
		if err != nil {
			fmt.Println("error in getting the ttl of the key")
		}

		resp := &types.QrResponse{
			URL:            reqBody.URL,
			QrCode:         png,
			Expiry:         ttl/time.Minute,
			XRateRemaining: int(remainingQuota),
			XRateLimitReset: r2.TTL(db.Ctx2, c.IP()).Val() / time.Second,
		}

		return c.Status(fiber.StatusOK).JSON(resp) // if the key exists then we will return the value of the key(i.e the reqBody struct which contains ) in the redis db.
	} else {

		if reqBody.ExpiryDate == 0 {
			reqBody.ExpiryDate = 24 * time.Hour //if the user does not provide the expiry date then we will set the expiry date to 24 hours.
		}

		//! Since  the Qr does not exist we create the qr code and store it in the redis db
		png, err := helpers.CreateQRCode(reqBody.URL)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "cannot generate QR code",
			})
		}

		//TODO after creating the qr code we will store it in the redis db
		//? if the key does not exist then we will set the key in the redis db. here key the url and the value is the reqBody struct. Yes we can set the struct as the value in the redis db.
		err=r2.Set(db.Ctx2, reqBody.URL, png, time.Hour*24).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save qrcode image in the redis db",
			})
		}

		//TODO after storing the qr code in the redis db we will decrement the ip counter for rate limiting
		remainingQuota, _ := r2.Decr(db.Ctx2, c.IP()).Result() //decrement the count of the ip address by 1
		if remainingQuota < 0 {
			ttl, _ := r2.TTL(db.Ctx2, c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":      "Rate limit exceeded",
				"retryAfter": ttl / time.Second,
			})
		}

		resp := &types.QrResponse{
			URL:            reqBody.URL,
			QrCode:         png,
			Expiry:         reqBody.ExpiryDate,
			XRateRemaining: int(remainingQuota),
			XRateLimitReset: r2.TTL(db.Ctx2, c.IP()).Val() / time.Second,
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}

}
