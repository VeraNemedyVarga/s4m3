package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type modelresponse struct {
	Board    Board `json:"board"`
	GameOver bool  `json:"gameOver"`
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func handleBoard(w http.ResponseWriter, _ *http.Request) {
	enableCors(&w)
	m := getWebModel()

	resp := modelresponse{
		Board:    m.board,
		GameOver: m.gameOver,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func initApi(config Config) {
	http.HandleFunc("/api/board", handleBoard)
	// http.Handle("/ws", wshandler())

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Printf("Listening on %s", config.Addr)
	log.Fatal(http.ListenAndServe(config.Addr, nil))
}
