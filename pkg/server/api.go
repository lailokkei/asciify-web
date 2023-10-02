package server

import (
	"bytes"
	"encoding/json"
	"image"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/toodemhard/asciify/pkg/asciify-lib"
)

var upgrader = websocket.Upgrader{}

func updateImage(message []byte) image.Image {
	img, _, err := image.Decode(bytes.NewReader(message))
	if err != nil {
		log.Println(err)
	}
	return img
}

func generateText(response chan string, message []byte, img image.Image) {
	go func() {
		options := asciify.Options{}
		err := json.Unmarshal(message, &options)
		if err != nil {
			response <- "Invalid json"
		}

		text, err := asciify.ImageToText(img, options)
		if err != nil {
			response <- "Invalid request options"
		}

		log.Println("img2text done")
		response <- text
	}()
}

func serveTextQueue(conn *websocket.Conn, quit chan int) chan chan string {
	newResponse := make(chan chan string)

	go func() {
		var responseQueue []chan string
		for {
			var firstResponse chan string
			if len(responseQueue) > 0 {
				firstResponse = responseQueue[0]
			}

			select {
			case response := <-newResponse:
				responseQueue = append(responseQueue, response)
			case message := <-firstResponse:
				log.Println("begin write")
				conn.WriteMessage(websocket.TextMessage, []byte(message))
				log.Println("end write")
				responseQueue = responseQueue[1:]
			case <-quit:
				log.Println("ded")
				return
			}
		}
	}()

	return newResponse

}

func connect(w http.ResponseWriter, r *http.Request) {
	img, err := asciify.DecodeImageFile("images/default.png")
	if err != nil {
		log.Println(err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	quit := make(chan int)
	responseChan := serveTextQueue(conn, quit)

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			quit <- 0
			log.Println(err)
			return
		}

		if mt == websocket.TextMessage {
			log.Println("request received")
			response := make(chan string)
			responseChan <- response
			generateText(response, message, img)
		}

		if mt == websocket.BinaryMessage {
			img = updateImage(message)
		}
	}
}
