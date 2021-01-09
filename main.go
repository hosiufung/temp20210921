package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
	//"time"
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

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	go func() {
		// FIXME, need waitgroup here
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	// received termination signal, perform graceful shutdown here
	<-signalChan
	log.Println("interrupt")

	// FIXME:  fire the termination signal to all worker wait for the death
	if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Println("write close:", err)
		return
	}

	/*
		// FIXME: it is not good for graceful shutdown
		for {
			select {
			case <-done:
				return
			case <-interrupt:

				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	*/
}
