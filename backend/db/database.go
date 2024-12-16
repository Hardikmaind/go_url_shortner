package db

import (
	"context"
	"os"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// ? this below context will be exported in the resolve.go file and will be used to create a context for the redis client
// ! this is a global variable which is used to create a context for the redis client
var (
	Ctx           = context.Background()
	Ctx2          = context.Background()
	CreateClient  *redis.Client
	CreateClient2 *redis.Client
)

//dummy:="abc"		//?this kind of declaration can only be used inside the function. This is a shorthand declaration. This is not allowed outside the function.

// first letter of function name should be capital to make it public and to export it to other packages
func InitRedisClient() {
	CreateClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("db_addr"),
		DB:       0,
		Password: os.Getenv("dp_pass"),
	})
	// Check if the Redis client can connect
	_, err := CreateClient.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Redis client: %v", err))
	}

}

// first letter of function name should be capital to make it public and to export it to other packages
// ! THIS FUNCTION CREATES A DIFFERENT REDIS CLIENT to DATABASE "1" FOR QR CODE STORING
func InitRedisClient2() {
	CreateClient2 = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("db_addr"),
		DB:       1,
		Password: os.Getenv("dp_pass"),
	})
	// Check if the Redis client can connect
	_, err := CreateClient2.Ping(Ctx2).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Redis client: %v", err))
	}

}

// ! THIS WE CAN USE IF WE NEED TO CREATE A NEW DB WITH .
// ? Yes, by default, Redis provides 16 databases (numbered 0 to 15) in a single instance. .
// Utility function for specific DB numbers, if needed
func GetClientForDB(dbNo int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("db_addr"),
		Password: os.Getenv("db_pass"),
		DB:       dbNo,
	})
}
