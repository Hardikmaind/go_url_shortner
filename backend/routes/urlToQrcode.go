package routes

import (
	"time"

	"github.com/Hardikmaind/go_url_shortner/db"
	"github.com/Hardikmaind/go_url_shortner/types"
	"github.com/gofiber/fiber/v2"
	qrcode "github.com/skip2/go-qrcode"
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

	if exists, _ := r2.Exists(db.Ctx2, reqBody.URL).Result(); exists > 0 {
		png, _ := r2.Get(db.Ctx2, reqBody.URL).Bytes()

		resp := &types.QrResponse{
			URL:            reqBody.URL,
			QrCode:         png,
			Expiry:         time.Hour * 24,
			XRateRemaining: 1000,
			XRateLimitRest: time.Hour * 24,
		}

		return c.Status(fiber.StatusOK).JSON(resp) // if the key exists then we will return the value of the key(i.e the reqBody struct which contains ) in the redis db.
	} else {

		//! Since  the Qr does not exist we create the qr code and store it in the redis db
		png, err := qrcode.Encode(reqBody.URL, qrcode.Medium, 256)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "cannot generate QR code",
			})
		}
		//? if the key does not exist then we will set the key in the redis db. here key the url and the value is the reqBody struct. Yes we can set the struct as the value in the redis db.
		r2.Set(db.Ctx2, reqBody.URL, png, time.Hour*24)



		resp := &types.QrResponse{
			URL:            reqBody.URL,
			QrCode:         png,
			Expiry:         time.Hour * 24,
			XRateRemaining: 1000,
			XRateLimitRest: time.Hour * 24,
		}
		return c.Status(fiber.StatusOK).JSON(resp)
	}

}
