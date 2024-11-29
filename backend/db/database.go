package db

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

//? this below context will be exported in the resolve.go file and will be used to create a context for the redis client
var Ctx = context.Background()			//! this is a global variable which is used to create a context for the redis client
//first letter of function name should be capital to make it public and to export it to other packages
func CreateClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("db_addr"),
		DB:       dbNo,
		Password: os.Getenv("dp_pass"),
	})

	// Use ctx if needed for further operations
	return rdb
}
