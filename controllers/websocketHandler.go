package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type sessionWrapper struct {
	upgrader   *websocket.Upgrader
	connection *websocket.Conn
}

func (s *sessionWrapper) upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) error {
	s.upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connection, err := s.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return err
	} else {
		s.connection = connection
		return nil
	}
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	var session = &sessionWrapper{}
	err := session.upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	log.Printf("Connected with %v", session.connection.RemoteAddr())
	defer session.connection.Close()

	for {
		_, _, err := session.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Disconnected: %v", session.connection.RemoteAddr())
			}
			break
		}
	}
}
