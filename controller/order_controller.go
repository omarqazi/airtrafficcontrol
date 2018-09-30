package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/omarqazi/airtrafficcontrol/model"
	"log"
	"net/http"
	"os"
)

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
	decoder := json.NewDecoder(r.Body)
	var requestOrder model.Order
	if err := decoder.Decode(&requestOrder); err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "cant parse your shitty json... check it for mistakes")
		return
	}

	ordersWaiting, err := redisClient.RPush("smick-drone-orders", requestOrder.ToJSON()).Result()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "error saving your order")
	} else {
		fmt.Fprintln(w, "your order is coming right up", ordersWaiting)
	}
}

func (oc OrderController) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	listResults, err := redisClient.LRange("smick-drone-orders", 0, -1).Result()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "error getting order status")
		return
	}

	parsedOrders := make([]model.Order, 0)

	for i := range listResults {
		var parsedOrder model.Order
		json.Unmarshal([]byte(listResults[i]), &parsedOrder)
		parsedOrders = append(parsedOrders, parsedOrder)
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(parsedOrders)
}
