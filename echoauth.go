package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func wsEchoAuthHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer conn.Close()

	for {
		msgData, msgOp, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Println(err)
			break
		}

		log.Println("data", msgData)

		err = wsutil.WriteServerMessage(conn, msgOp, msgData)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
