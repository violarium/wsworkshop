package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, http.StatusOK, struct {
			Message string `json:"message"`
		}{Message: "Websocket Workshop"})
	})

	files := []string{"client-echo.html", "client-echo-auth.html", "client-echo-ping.html", "event-source.html", "favicon.ico"}
	for _, file := range files {
		f := file
		r.Get("/"+file, func(w http.ResponseWriter, r *http.Request) {
			fileContent, err := os.ReadFile(f)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if _, err := w.Write(fileContent); err != nil {
				log.Println(err)
			}
		})
	}

	r.Get("/ws-echo", wsEchoHandler)
	r.Get("/ws-echo-auth", wsEchoAuthHandler)
	r.Get("/ws-echo-ping", wsEchoPingHandler)
	r.Get("/event-source", eventSourceHandler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}

func writeJson(w http.ResponseWriter, status int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Println(message)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Println(err)
	}
}
