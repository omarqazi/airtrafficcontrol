package controller

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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
	go lc.WriteMessagesFromChannel(conn, r)
}

func (lc LiveController) ReadMessagesAndPostToChannel(conn *websocket.Conn, r *http.Request) {
	redisChannel := lc.channelNameForRequest(r)

	for {
		_, sb, err := conn.ReadMessage()
		if err != nil {
			return
		}

		redisClient.Publish(redisChannel, string(sb))
	}
}

func (lc LiveController) WriteMessagesFromChannel(conn *websocket.Conn, r *http.Request) {
	redisChannel := lc.channelNameForRequest(r)
	pubsub := redisClient.Subscribe(redisChannel)
	defer pubsub.Close()

	if _, err := pubsub.Receive(); err != nil {
		log.Println("Error connecting to pubsub:", err)
		return
	}

	ch := pubsub.Channel()

	for msg := range ch {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			return
		}
	}
}

func (lc LiveController) channelNameForRequest(r *http.Request) string {
	orderId := r.URL.Path
	channelName := "drone-order-updates-" + orderId
	return channelName
}
