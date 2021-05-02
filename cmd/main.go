package main

import (
	"github.com/gorilla/websocket"
	"github.com/yigitsadic/minigame/internal/game"
	"log"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return r.Host == "localhost:9090"
		},
	}
	g *game.Game
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)

		return
	}

	p := game.NewPlayer(conn.RemoteAddr().String(), conn)

	go g.JoinPlayer(p)
	go g.HandleGame()
}

func main() {
	g = game.NewGame()

	http.HandleFunc("/ws", serveWs)

	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("Error occurred : %s", err)
	}
}
