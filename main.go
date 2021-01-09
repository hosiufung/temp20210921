package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	//"net/url"
	"flag"
	"os"
	"os/signal"
	"time"
)

// FIXME: instead of hardcode, use viper to store such info
const apiKey = `6dcfa17346295ab13be431dba51993666601f924a326707296808acef673c247`

type JSONParams struct {
	Action string   `json:"action"`
	Subs   []string `json:"subs"`
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// this is where you paste your api key
	c, _, err := websocket.DefaultDialer.Dial("wss://streamer.cryptocompare.com/v2?api_key="+apiKey, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	jsonObj := JSONParams{Action: "SubAdd", Subs: []string{"8~Binance~BTC~USDT"}}
	s, _ := json.Marshal(jsonObj)
	fmt.Println(string(s))
	err = c.WriteMessage(websocket.TextMessage, []byte(string(s)))
	if err != nil {
		log.Fatal("message:", err)
	}

	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	// FIXME: it is not good for graceful shutdown
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
