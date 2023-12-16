package main

import (
	game "black-jack/card"
	"black-jack/config"
	"black-jack/repository"
	"black-jack/ws"
	"flag"
	"log"
	"math/rand"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("hello"))
}

func main() {
	db := config.InitDB()
	userRepository := repository.NewUserRepository(db)
	center := ws.NewGameCenter(
		userRepository,
		ws.NewRoom("一茗", game.NewCardDealer(rand.New(rand.NewSource(12345)))),
		ws.NewRoom("二穴", game.NewCardDealer(rand.New(rand.NewSource(11111)))),
		ws.NewRoom("三井", game.NewCardDealer(rand.New(rand.NewSource(78945)))),
	)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(center, w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
