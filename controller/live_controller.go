package controller

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type LiveController struct {
}

func (lc LiveController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading websocket: ", err)
		w.WriteHeader(500)
		fmt.Fprintln(w, "error upgrading connection")
		return
	}

	go lc.ReadMessagesAndPostToChannel(conn, r)

	for {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("hi doggie")); err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (lc LiveController) ReadMessagesAndPostToChannel(conn *websocket.Conn, r *http.Request) {
	for {
		orderId := r.URL.Path
		pubsubChannel := "drone-order-updates-" + orderId

		_, socketBytes, err := conn.ReadMessage()
		if err != nil {
			return
		}

		redisClient.Publish(pubsubChannel, string(socketBytes))
	}
}
