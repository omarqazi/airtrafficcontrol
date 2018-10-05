package controller

import (
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"os"
)

const DefaultHost = ":8080"

var Order = OrderController{}
var Home = HomeController{}
var Live = LiveController{}

var redisClient *redis.Client = nil

func init() {
	redisAddress := "localhost:6379"
	redisOverride := os.Getenv("REDIS_ADDR")
	if len(redisOverride) > 0 {
		redisAddress = redisOverride
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})

	if err := redisClient.Ping().Err(); err != nil {
		log.Fatalln("Error connecting to redis: ", err)
	}
}

func RegisterHandlers() error {
	http.Handle("/", Home)
	http.Handle("/order/", http.StripPrefix("/order/", Order))
	http.Handle("/live/", http.StripPrefix("/live/", Live))
	return nil
}
