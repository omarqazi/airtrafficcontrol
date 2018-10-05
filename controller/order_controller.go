package controller

import (
	"encoding/json"
	"fmt"
	"github.com/omarqazi/airtrafficcontrol/model"
	"io/ioutil"
	"net/http"
	"strings"
)

type OrderController struct {
}

func (oc OrderController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.URL.Path

	w.Header().Add("Content-Type", "application/json")

	if len(requestUrl) == 0 && r.Method == "GET" {
		oc.OrderQueueLength(w, r)
	} else if r.Method == "POST" {
		oc.TakeOrder(w, r)
	} else if r.Method == "GET" && requestUrl == "list" {
		oc.GetAllOrders(w, r)
	} else if r.Method == "GET" && requestUrl == "next" {
		oc.GetNextOrder(w, r)
	} else if r.Method == "PUT" {
		oc.PopTopOrder(w, r)
	} else if r.Method == "PATCH" {
		oc.UpdateOrderStatus(w, r)
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
	requestOrder.GenerateId()

	ordersWaiting, err := redisClient.RPush("smick-drone-orders", requestOrder.ToJSON()).Result()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "error saving your order")
	} else {
		fmt.Fprintln(w, "your order is coming right up", ordersWaiting)
	}
}

func (oc OrderController) GetNextOrder(w http.ResponseWriter, r *http.Request) {
	topItem, err := redisClient.LIndex("smick-drone-orders", 0).Result()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "error getting next order")
		return
	}

	var topOrder model.Order
	json.Unmarshal([]byte(topItem), &topOrder)

	encoder := json.NewEncoder(w)
	encoder.Encode(topOrder)
}

func (oc OrderController) GetAllOrders(w http.ResponseWriter, r *http.Request) {
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

func (oc OrderController) PopTopOrder(w http.ResponseWriter, r *http.Request) {
	topOrder, err := redisClient.LPop("smick-drone-orders").Result()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Error popping top order")
		return
	}

	var parsedOrder model.Order
	json.Unmarshal([]byte(topOrder), &parsedOrder)

	encoder := json.NewEncoder(w)
	encoder.Encode(parsedOrder)
}

func (oc OrderController) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	requestUrl := r.URL.Path
	comps := strings.Split(requestUrl, "/")
	orderId := comps[0]
	pubsubChannel := "drone-order-updates-" + orderId

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "error reading request body")
		return
	}

	messageToPublish := string(requestBody)
	redisClient.Publish(pubsubChannel, messageToPublish)

	fmt.Fprintln(w, messageToPublish)

}
