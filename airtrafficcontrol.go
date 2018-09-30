package main

import (
	"github.com/omarqazi/airtrafficcontrol/controller"
	"log"
	"net/http"
)

func main() {
	log.Println("Drone Air traffic control started")

	controller.RegisterHandlers()
	log.Println("Registered controllers... starting HTTP server")

	if err := http.ListenAndServe(controller.DefaultHost, nil); err != nil {
		log.Println("Error starting server:", err)
	}
}
