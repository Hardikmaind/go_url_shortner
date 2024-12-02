package routes

import (
	"fmt"
	"os"
	"time"

	"github.com/Hardikmaind/go_url_shortner/db"
	"github.com/Hardikmaind/go_url_shortner/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"customShort"`
	ExpiryDate  time.Duration `json:"expiryDate"`
}
type Response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"customShort"`
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"xRateRemaining"`
	XRateLimitRest time.Duration `json:"xRateLimitRest"`
}

func ShortenUrl(c *fiber.Ctx) error {
	//* we can also use "var reqBody RequestBody"..but this does not allocate memory to the struct. and does not give pointer to the struct. so we use new keyword to allocate memory to the struct.
	reqbody := new(Request) //this new keyword is used to create a new instance of the Request struct.This was introduced in Go 1.4. It allocates zeroed storage for a new item and returns a pointer to it.

	//now we are binding the request body to the Request struct.
	//? c.BodyParser binds the request body to a struct. It supports JSON, form, query, and multipart requests.
	//also we cannot direcly do reqbody:=c.Body() as it returns a byte slice. to get body in json format we use c.BodyParser. similarly like in express we do app.use(body.parser())

	//?* now instead of passing the "reqbody" to it we pass the address of the reqbody. so in this way we avoid the copying of the "reqbody" struct. This is a good practice to pass the address of the struct to the function.
	if err := c.BodyParser(&reqbody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ //to retun the response in json format we send it in a fiber.Map. enclosing in a Json method will automatically set the content type to application/json
			"error": "cannot parse JSON",
		})
	}

	//!TODO: IF reqbody.url is empty then return an alternative response
	if reqbody.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL cannot be empty",
		})
	}

	//!TODO: Implement Rate Limiting here

	//r2 := db.CreateClient(1) //this is the client for the rate limiting. we are creating a new client for the rate limiting. we can also use the same client which we used for the url shortening. but it is a good practice to use a different client for the rate limiting.

	r := db.CreateClient //? this is the  client to same redis database which was initialised in main.go file. we are using the same client for the rate limiting. we can also use a different client for the rate limiting. but it is a good practice to use the same client for the rate limiting and the url shortening. and why create the redis client for each request to resolve function.
	//defer r.Close()		//SINCE THIS IS THE SAME CLIENT IN THE MAIN.GO..WE WILL TERMINATE THIS CLIENT WHEN THE APPLICATION TERMINATES. SO WE WILL WRITE THIS LINE IN MAIN.GO INSTED OF HERE.

	//! ALSO WE CAN USE THE SAME DB FOR COUNTER AND KEY VALUUE. ALSO WE CAN USE MULTIPLE COUNTER IN THE SAME REDIS DB

	_, err := r.Get(db.Ctx, c.IP()).Result()
	if err == redis.Nil {
		//! Set rate limit and expiration for new IP
		err := r.Set(db.Ctx, c.IP(), os.Getenv("api_quota"), 2*time.Minute).Err() // Set key with 1-minute expiration
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to set rate limit",
			})
		}
	} else {
		valInt, _ := r.Get(db.Ctx, c.IP()).Int()
		if valInt <= 0 {
			limit, _ := r.TTL(db.Ctx, c.IP()).Result() // Check TTL for the current IP

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

	//!TODO: check whether the input is an actual URl

	if !govalidator.IsURL(reqbody.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid URL",
		})
	}

	//!TODO:  check for domain error...(check whether the use is not using the localhost or any other domain which is not allowed)
	if !helpers.IsRequestURLAllowed(reqbody.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "This URL is not allowed",
		})
	}

	//!tODO: enfore https/SSL
	//* method 1. this we can do to validate
	/*
		if reqbody.URL[:5] != "https" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "URL must be HTTPS",
				})
				}
	*/

	//* method 2. this we can do to check and replace the http with https
	reqbody.URL = helpers.EnforceHTTPS(reqbody.URL)

	//!TODO: check if the custom short URL already exists.which means check if we have already created the custom short URL in the redis database.

	//! WE HAVE ALREADY CREATED THE CLIENT FOR IT SEE ABOVE.
	//r := db.CreateClient			//? this is the  client to same redis database which was initialised in main.go file. we are using the same client for the rate limiting. we can also use a different client for the rate limiting. but it is a good practice to use the same client for the rate limiting and the url shortening. and why create the redis client for each request to resolve function.

	exists, _ := r.Exists(db.Ctx, reqbody.URL).Result()
	if exists > 0 {
		id, _ := r.Get(db.Ctx, reqbody.URL).Result()

		remainingQuota, _ := r.Decr(db.Ctx, c.IP()).Result() //decrement the count of the ip address by 1
		if remainingQuota < 0 {
			ttl, _ := r.TTL(db.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":      "Rate limit exceeded",
				"retryAfter": ttl / time.Second,
			})
		}
		domain := os.Getenv("domain")

		//!TODO: return the response
		//? this is method 1
		// resp:=new(Response)
		// resp.URL=reqbody.URL
		// resp.CustomShort=domain + "/" + id,
		// resp.Expiry=reqbody.ExpiryDate
		// resp.XRateRemaining,_=int(remainingQuota),
		// resp.XRateLimitRest,_=r2.TTL(db.Ctx, c.IP()).Val() / time.Second

		//? this is method 2 to do above

		ttl, err := r.TTL(db.Ctx, reqbody.URL).Result() //? Fixed the expiry..now expiry is consistent, in cache hit or miss
		if err != nil {
			fmt.Println("error in getting the ttl of the key")
		}

		resp := &Response{
			URL:            reqbody.URL,
			CustomShort:    domain + "/" + id,
			Expiry:         ttl / time.Minute,
			XRateRemaining: int(remainingQuota),
			XRateLimitRest: r.TTL(db.Ctx, c.IP()).Val() / time.Second,
		}

		return c.Status(fiber.StatusCreated).JSON(resp)
	} else {
		id := helpers.GenerateRandomString(7) //*declare the id string for generating the short url HASH

		fmt.Println("id:", id)

		//!TODO: save the URL in Redis
		if reqbody.ExpiryDate == 0 {
			reqbody.ExpiryDate = 24 * time.Hour //if the user does not provide the expiry date then we will set the expiry date to 24 hours.
		}

		//! THIS BELOW I HAVE DONE THE DUAL MAPPING.URL -> ID AND ID -> URL
		if err := r.Set(db.Ctx, reqbody.URL, id, reqbody.ExpiryDate).Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save URL in Redis",
			})
		}

		err = r.Set(db.Ctx, id, reqbody.URL, reqbody.ExpiryDate).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to save reverse mapping in Redis",
			})
		}

		remainingQuota, _ := r.Decr(db.Ctx, c.IP()).Result() //decrement the count of the ip address by 1
		if remainingQuota < 0 {
			ttl, _ := r.TTL(db.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":      "Rate limit exceeded",
				"retryAfter": ttl / time.Second,
			})
		}
		domain := os.Getenv("domain")

		//!TODO: return the response
		//? this is method 1
		// resp:=new(Response)
		// resp.URL=reqbody.URL
		// resp.CustomShort=domain + "/" + id,
		// resp.Expiry=reqbody.ExpiryDate
		// resp.XRateRemaining,_=int(remainingQuota),
		// resp.XRateLimitRest,_=r2.TTL(db.Ctx, c.IP()).Val() / time.Second

		//? this is method 2 to do above
		resp := &Response{
			URL:            reqbody.URL,
			CustomShort:    domain + "/" + id,
			Expiry:         reqbody.ExpiryDate,
			XRateRemaining: int(remainingQuota),
			XRateLimitRest: r.TTL(db.Ctx, c.IP()).Val() / time.Second,
		}

		return c.Status(fiber.StatusCreated).JSON(resp)

	}

}
