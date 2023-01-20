package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type modelresponse struct {
	Board    Board `json:"board"`
	GameOver bool  `json:"gameOver"`
	Points   int   `json:"points"`
}

type WebHit struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	resp chan model
}

func (w WebHit) getCoords() (int, int) {
	return w.X, w.Y
}

func (w WebHit) getResp() chan model {
	return w.resp
}

type WebGet struct {
	resp chan model
}

func (w WebGet) getResp() chan model {
	return w.resp
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func handleBoard(sub chan webHitMsg) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		ch := make(chan model)

		if r.Method == "POST" {
			decoder := json.NewDecoder(r.Body)
			var hit WebHit
			err := decoder.Decode(&hit)
			hit.resp = ch
			if err != nil {
				log.Println(err)
			}
			defer r.Body.Close()

			sub <- hit
		} else {
			sub <- WebGet{resp: ch}
		}

		m := <-ch
		resp := modelresponse{
			Board:    m.board,
			GameOver: m.gameOver,
			Points:   m.points,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func initApi(config Config, sub chan webHitMsg) {
	http.HandleFunc("/api/board", handleBoard(sub))
	// http.Handle("/ws", wshandler())

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Printf("Listening on %s", config.Addr)
	log.Fatal(http.ListenAndServe(config.Addr, nil))
}
