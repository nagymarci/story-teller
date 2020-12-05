package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/nagymarci/story-teller/events"
	"github.com/nagymarci/story-teller/store"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func EventsWSHandler(router *mux.Router, sub *events.InApp, s *store.Default) {
	router.HandleFunc("/{gameId}/ws", func(w http.ResponseWriter, r *http.Request) {
		gameID := mux.Vars(r)["gameId"]

		if _, err := s.Load(gameID); err != nil {
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			if _, ok := err.(websocket.HandshakeError); !ok {
				http.Error(w, "Unknown error", http.StatusInternalServerError)
			}
			return
		}

		eventChannel, err := sub.Subscribe(gameID, ws)
		if err != nil {
			http.Error(w, "Unable to subscribe", http.StatusInternalServerError)
			return
		}

		go writer(ws, eventChannel, sub, gameID)
		reader(ws, sub, gameID)
	})
}

const (
	pongWait   = 30 * time.Second
	pingPeriod = (pongWait * 8) / 10
)

func writer(ws *websocket.Conn, events <-chan interface{}, sub *events.InApp, gameID string) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		sub.Unsubscribe(gameID, ws)
		pingTicker.Stop()
		ws.Close()
	}()

	for {
		select {
		case e := <-events:
			if err := ws.WriteJSON(e); err != nil {
				return
			}
		case <-pingTicker.C:
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func reader(ws *websocket.Conn, sub *events.InApp, gameID string) {
	defer func() {
		sub.Unsubscribe(gameID, ws)
		ws.Close()
	}()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
