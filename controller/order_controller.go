package controller

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"net/http"
)

var redisClient *redis.Client = nil

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := redisClient.Ping().Err(); err != nil {
		log.Fatalln("Error connecting to redis: ", err)
	}
}

type OrderController struct {
}

func (oc OrderController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ordersWaiting, err := redisClient.LLen("smick-drone-orders").Result()
	if err != nil {
		fmt.Fprintln(w, "error talking to redis")
		return
	}

	fmt.Fprintln(w, "orders in line:", ordersWaiting)
}
