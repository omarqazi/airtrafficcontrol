package controller

import "net/http"

const DefaultHost = ":8080"

var Order = OrderController{}
var Home = HomeController{}

func RegisterHandlers() error {
	http.Handle("/", Home)
	http.Handle("/order/", http.StripPrefix("/order/", Order))
	return nil
}
