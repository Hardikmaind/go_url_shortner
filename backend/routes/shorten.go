package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/Hardikmaind/go_url_shortner/db"
	"github.com/Hardikmaind/go_url_shortner/helpers"
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

func shortenUrl(c *fiber.Ctx) error {
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

	//! Implement Rate Limiting here

	r2 := db.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(db.Ctx, c.IP()).Result()
	if err == redis.Nil {
		r2.Set(db.Ctx, c.IP(), os.Getenv("api_quota"), 30*time.Minute).Err()
	} else {
		// val, _ := r2.Get(db.Ctx, c.IP()).Result()			//instead of this we can dirctly get the val in int format by val,_:=r2.Get(db.Ctx,c.IP()).Int()
		// valInt, _ := strconv.Atoi(val)
		valInt,_:=r2.Get(db.Ctx,c.IP()).Int()
		if valInt <= 0 {
			limit,_:=r2.TTL(db.Ctx,c.IP()).Result()		//this will return the time left for the key to expire.in this case the key is the ip address of the user.
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
				"retryAfter": limit/time.Nanosecond/time.Minute,
			})
		}
	}

	// ! check whether the input is an actual URl

	if !govalidator.IsURL(reqbody.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid URL",
		})
	}

	//! check for domain error...(check whether the use is not using the localhost or any other domain which is not allowed)
	if !helpers.IsRequestURLAllowed(reqbody.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "This URL is not allowed",
		})
	}

	//! enfore https/SSL
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


	r2.Decr(db.Ctx, c.IP()) //decrement the count of the ip address by 1
	return c.JSON(fiber.Map{
		"message": "shortenUrl",
	})

}