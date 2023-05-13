package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func eventSourceHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ctx := r.Context()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("closed")
			return
		case t := <-ticker.C:
			log.Println(t)

			data := "message " + t.String()
			if _, err := fmt.Fprintf(w, "data: %v\n\n", data); err != nil {
				log.Println(err)
				return
			}

			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			} else {
				log.Println("cant' flush")
				return
			}
		}
	}
}
