package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mobocrat/trade/sampling"
	"github.com/mobocrat/trade/stock"
	"github.com/mobocrat/trade/trade"
)

const (
	samplingInterval = time.Second
	samplesPerUpdate = 40
)

func main() {
	market := stock.New(samplingInterval / 5)
	sampler := sampling.New(samplingInterval, market)
	traders := []*trade.Trader{}
	for i := 0; i < 1000; i++ {
		traders = append(traders, trade.New(market))
	}
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	http.HandleFunc("/price", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("upgrade: %v\n", err)
			return
		}
		go func(conn *websocket.Conn) {
			ticker := time.NewTicker(samplingInterval)
			for range ticker.C {
				err := conn.WriteJSON(sampler.Sample(samplesPerUpdate))
				if err != nil {
					fmt.Printf("write: %v\n", err)
					break
				}
			}
			ticker.Stop()
			conn.Close()
		}(conn)
	})
	http.HandleFunc("/confidence", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("upgrade: %v\n", err)
			return
		}
		go func(conn *websocket.Conn) {
			ticker := time.NewTicker(samplingInterval)
			for range ticker.C {
				h := make([]int, 101)
				for _, t := range traders {
					if c := t.Confidence(); c >= 0 && c <= 100 {
						h[c]++
					}
				}
				err := conn.WriteJSON(h)
				if err != nil {
					fmt.Printf("write: %v\n", err)
					break
				}
			}
			ticker.Stop()
			conn.Close()
		}(conn)
	})
	http.HandleFunc("/asset", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("upgrade: %v\n", err)
			return
		}
		go func(conn *websocket.Conn) {
			ticker := time.NewTicker(samplingInterval)
			for range ticker.C {
				h := make([]int, 21)
				for _, t := range traders {
					if c := t.Assets(); c >= 0 && c <= 20 {
						h[c]++
					}
				}
				err := conn.WriteJSON(h)
				if err != nil {
					fmt.Printf("write: %v\n", err)
					break
				}
			}
			ticker.Stop()
			conn.Close()
		}(conn)
	})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
