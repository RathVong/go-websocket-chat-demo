package main

import (
	"encoding/json"
	"io"
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// message sent to us by the javascript client
type message struct {
	Handle string `json:"handle"`
	Text   string `json:"text"`
}

// validateMessage so that we know it's valid JSON and contains a Handle and
// Text
func validateMessage(data []byte) (message, error) {
	var msg message

	if err := json.Unmarshal(data, &msg); err != nil {
		return msg, errors.Wrap(err, "Unmarshaling message")
	}

	if msg.Handle == "" && msg.Text == "" {
		return msg, errors.New("Message has no Handle or Text")
	}

	return msg, nil
}

// handleWebsocket connection.
func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		m := "Unable to upgrade to websockets"
        log.Printf("err: %v msg -> %v", err, m)
		http.Error(w, m, http.StatusBadRequest)
		return
	}

	rr.register(ws)

	for {
		mt, data, err := ws.ReadMessage()
    
        	if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) || err == io.EOF {
				log.Println("Websocket closed!")
				break
			}
		
			log.Println("Error reading websocket message")
		}
	
	switch mt {
		case websocket.TextMessage:
			msg, err := validateMessage(data)
			if err != nil {
				log.Printf("msg %v err %v msg -> Invalid Message", msg, err)
				break
			}
			rw.publish(data)
		default:
		    log.Println("Unknown Message!")
        }
	}

	rr.deRegister(ws)

	ws.WriteMessage(websocket.CloseMessage, []byte{})
}
