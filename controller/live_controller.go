package controller

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const pubsubDelimeter = "//&/"

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

	clientUniqueId, _ := uuid.NewRandom()
	clientId := clientUniqueId.String()

	go lc.ReadMessagesAndPostToChannel(conn, r, clientId)
	go lc.WriteMessagesFromChannel(conn, r, clientId)
}

func (lc LiveController) ReadMessagesAndPostToChannel(conn *websocket.Conn, r *http.Request, clientId string) {
	redisChannel := lc.channelNameForRequest(r)

	for {
		_, sb, err := conn.ReadMessage()
		if err != nil {
			return
		}

		clientPayload := string(sb)
		serializedPayload := clientId + pubsubDelimeter + clientPayload

		redisClient.Publish(redisChannel, serializedPayload)
	}
}

func (lc LiveController) WriteMessagesFromChannel(conn *websocket.Conn, r *http.Request, clientId string) {
	redisChannel := lc.channelNameForRequest(r)
	pubsub := redisClient.Subscribe(redisChannel)
	defer pubsub.Close()

	if _, err := pubsub.Receive(); err != nil {
		log.Println("Error connecting to pubsub:", err)
		return
	}

	ch := pubsub.Channel()

	for msg := range ch {
		messagePayload := msg.Payload

		comps := strings.Split(messagePayload, pubsubDelimeter)
		messageSender := comps[0]
		messageBody := strings.Join(comps[1:], pubsubDelimeter)

		if messageSender != clientId {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(messageBody)); err != nil {
				return
			}
		}
	}
}

func (lc LiveController) channelNameForRequest(r *http.Request) string {
	orderId := r.URL.Path
	channelName := "drone-order-updates-" + orderId
	return channelName
}
