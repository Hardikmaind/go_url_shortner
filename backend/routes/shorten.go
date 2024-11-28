package routes

import (
	"time"

	govalidator "github.com/Hardikmaind/go_url_shortner/helpers"
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
	if reqbody.URL[:5] != "https" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "URL must be HTTPS",
		})
	}
	



	return c.JSON(fiber.Map{
		"message": "shortenUrl",
	})

}
