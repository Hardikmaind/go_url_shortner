package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.resolveUrl)
	app.Post("/:api/v1", routes.shortenUrl)
}
func main() {

	/*func godotenv.Load(filenames ...string) (err error)
	Load will read your env file(s) and load them into ENV for this process.

	Call this function as close as possible to the start of your program (ideally in main).

	If you call Load without any args it will default to loading .env in the current path.

	You can otherwise tell it which files to load (there can be more than one) like:

	godotenv.Load("fileone", "filetwo")
	It's important to note that it WILL NOT OVERRIDE an env variable that already exists - consider the .env file to set dev vars or sensible defaults.*/

	
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
	}
	app := fiber.New() //this returns a pointer to a new instance of the Fiber struct
	//! i can also define routes like this ...or in another function. Like in above function
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })
	app.Use(logger.New())
	setupRoutes(app)

	//? I could also write below as log.Fatal(app.Listen(":3000"))
	// run:=app.Listen(os.Getenv("app_port"))
	// log.Fatal(run)
	log.Fatal(app.Listen(os.Getenv("app_port")))

}
