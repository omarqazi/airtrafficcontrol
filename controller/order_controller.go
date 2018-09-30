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
	requestUrl := r.URL.Path

	if len(requestUrl) == 0 && r.Method == "GET" {
		oc.OrderQueueLength(w, r)
	} else if r.Method == "POST" {
		oc.TakeOrder(w, r)
	} else if r.Method == "GET" {
		oc.GetOrderStatus(w, r)
	}
}

func (oc OrderController) OrderQueueLength(w http.ResponseWriter, r *http.Request) {
	ordersWaiting, err := redisClient.LLen("smick-drone-orders").Result()
	if err != nil {
		fmt.Fprintln(w, "error talking to redis")
		return
	}

	fmt.Fprintln(w, "orders in queue:", ordersWaiting)
}

func (oc OrderController) TakeOrder(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "your order is coming right up")
}

func (oc OrderController) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "uhhh haha yeah we are working on it")
}
