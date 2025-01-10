package utils

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/websocket"
)

func SocketConnection(WSHandler http.HandlerFunc) (*websocket.Conn, *http.Response, error, func()) {
	ts := httptest.NewServer(http.HandlerFunc(WSHandler))
	defer ts.Close()

	wsURL := "ws" + ts.URL[4:]

	dialer := websocket.DefaultDialer

	conn, resp, err := dialer.Dial(wsURL, nil)
	deferFunc := func() {
		if conn != nil {
			conn.Close()
		}
		ts.Close()
	}
	return conn, resp, err, deferFunc
}
