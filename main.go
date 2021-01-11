package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	//"time"

	"github.com/gorilla/websocket"
)

// FIXME: instead of hardcode, use viper to store such info
const apiKey = `6dcfa17346295ab13be431dba51993666601f924a326707296808acef673c247`

type JSONParams struct {
	Action string   `json:"action"`
	Subs   []string `json:"subs"`
}

// such bool burrowed from go std lib
type atomicBool int32

func (b *atomicBool) isTrue() bool { return atomic.LoadInt32((*int32)(b)) != 0 }
func (b *atomicBool) setTrue()     { atomic.StoreInt32((*int32)(b), 1) }
func (b *atomicBool) setFalse()    { atomic.StoreInt32((*int32)(b), 0) }

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

	temp := atomicBool(0)
	isClosed := &temp

	orderBook := NewOrderBook()
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go wssReceiver(c, orderBook, isClosed, wg)
	go printMidPrice(orderBook)

	// received termination signal, perform graceful shutdown here
	<-signalChan
	log.Println("interrupt")

	// fire the termination signal
	isClosed.setTrue()

	// wait for worker die first
	wg.Wait()

	// fire the termination message in wss
	if err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Println("write close:", err)
		return
	}

}
