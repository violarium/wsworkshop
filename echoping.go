package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func wsEchoPingHandler(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer conn.Close()

	msgToSend := make(chan wsutil.Message)
	var pinger int32
	pingTicker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			msgList, err := wsutil.ReadClientMessage(conn, []wsutil.Message{})
			if err != nil {
				log.Println(err)
				return
			}

			for _, msg := range msgList {
				if msg.OpCode.IsControl() {
					wsutil.HandleControlMessage(conn, ws.StateServerSide, msg)
				}

				if msg.OpCode == ws.OpPong {
					log.Println("pong")
					atomic.AddInt32(&pinger, -1)
				}

				if msg.OpCode.IsData() {
					log.Println("data", msg.Payload)
					msgToSend <- msg
				}
			}
		}
	}()

	for {
		select {
		case <-pingTicker.C:
			if atomic.LoadInt32(&pinger) != 0 {
				return
			}
			atomic.AddInt32(&pinger, 1)

			log.Println("ping")
			if err := wsutil.WriteServerMessage(conn, ws.OpPing, nil); err != nil {
				return
			}
		case msg := <-msgToSend:
			err = wsutil.WriteServerMessage(conn, msg.OpCode, msg.Payload)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
