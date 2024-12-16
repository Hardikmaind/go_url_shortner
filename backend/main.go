package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/Hardikmaind/go_url_shortner/db"
	"github.com/Hardikmaind/go_url_shortner/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(api fiber.Router) {
	// Define routes relative to the group prefix
	api.Get("/:url", routes.ResolveUrl) // e.g., GET /api/v1/some-url
	api.Post("/", routes.ShortenUrl)    // e.g., POST /api/v1/shorten
	api.Post("/qr", routes.UrlToQrcode) // e.g., POST /api/v1/qr
}

func InitRedisClients() {

	db.InitRedisClient() //! HERE WE INITIALIZE THE NEW REDIS INSTANCE

	db.InitRedisClient2() //! HERE WE INITIALIZE THE NEW REDIS INSTANCE FOR THE QR CODE STORING. WE CAN ALSO USE THE SAME BUT TO KEEP THE SEPARATION OF CONCERNS WE ARE USING A DIFFERENT DB FOR QR CODE STORING

}

var embeddedFiles embed.FS

func main() {
	err := godotenv.Load()
	/*func godotenv.Load(filenames ...string) (err error)
	Load will read your env file(s) and load them into ENV for this process.

	Call this function as close as possible to the start of your program (ideally in main).

	If you call Load without any args it will default to loading .env in the current path.

	You can otherwise tell it which files to load (there can be more than one) like:

	godotenv.Load("fileone", "filetwo")
	It's important to note that it WILL NOT OVERRIDE an env variable that already exists - consider the .env file to set dev vars or sensible defaults.*/

	InitRedisClients() //? this function here will intialize the redis client
	//DONT WRITE THE DEFER STATEMENT IN THE "InitRedisClients" FUNCTION AS IT WILL CLOSE THE CLIENTS IMMEDIATELY AFTER THE FUNCTION EXECUTION. SO WE WILL WRITE THE DEFER STATEMENT HERE IN MAIN FUNCTION
	defer db.CreateClient.Close()
	defer db.CreateClient2.Close()

	if err != nil {
		fmt.Println("Error loading .env file")
	}
	app := fiber.New() //this returns a pointer to a new instance of the Fiber struct
	//! i can also define routes like this ...or in another function. Like in above function
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })
	app.Use(logger.New())
	// Enable CORS for all origins
	// CORS middleware with configuration to allow your frontend URL
	// Enable CORS for all origins
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",                                           // Allows all origins
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",                 // Allows these HTTP methods
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", // Allows these headers
		// AllowCredentials: true, // Include credentials if needed
	}))

	//! CREATED A ROUTE GROUP SO THAT I CAN RUN THE REST API ON /api/v1/ ROUTE
	api := app.Group("/api/v1")
	setupRoutes(api)

// !========================Serve static files from the embedded frontend build folder==========================================

app.Static("/", "./frontend/urlshortner/dist.index.html")

// Fallback route for React client-side routing
app.Use(func(c *fiber.Ctx) error {
	content, err := embeddedFiles.ReadFile("frontend/urlshortner/dist/index.html")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Index file not found")
	}
	return c.Type("html").SendString(string(content))
})

// !=================================================================================================================

	//? I could also write below as log.Fatal(app.Listen(":3000"))
	// run:=app.Listen(os.Getenv("app_port"))
	// log.Fatal(run)
	log.Fatal(app.Listen(os.Getenv("app_port")))

}
