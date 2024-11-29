package routes

import (
	"github.com/Hardikmaind/go_url_shortner/db"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveUrl(c *fiber.Ctx) error {
	url := c.Params("url")

	r := db.CreateClient(0)				//now here everytime a resolveUrl function is called, a new redis client is created and then closed. This is not a good practice. So we will create a single redis client(shared rdb) and use it everywhere.
	defer r.Close()

	//? Check in the database. if url is not in the database, return 404
	value, err := r.Get(db.Ctx, url).Result()
    if err == redis.Nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{		//! we can direcly code the status code like this=>"return c.Status(404).JSON(fiber.Map{" or we can use the fiber package constants like this=>"return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"
            "error": "URL not found",
        })
    } else if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Internal Server Error",
        })
    }

	//* this rInr tracks the usage of the shorten url

	rInr:=db.CreateClient(1)			//now here everytime a resolveUrl function is called, a new redis client is created and then closed. This is not a good practice. So we will create a single redis client(shared rdb) and use it everywhere.
	defer rInr.Close()

	_=rInr.Incr(db.Ctx,"counter")		//increment the count of the url by 1

	return c.Redirect(value, 301)			//redirect with status code 301. 301 means that the resource has been permanently moved to a new location, and future references should use a new URI with their requests.

	//and if sent 302 then it means that the resource has been temporarily moved to a new location.

}
