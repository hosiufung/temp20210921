package main

import (
	"encoding/json"
	"sync"

	"log"

	"github.com/gorilla/websocket"
)

func wssReceiver(c *websocket.Conn, isClosed *atomicBool, wg *sync.WaitGroup) {
	defer wg.Done()
	for !isClosed.isTrue() {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			//return
		}
		log.Printf("recv: %s", message)

		bm := Order{}
		err = json.Unmarshal([]byte(message), &bm)
		log.Println("bm, err", bm, err)
	}
}

type Order struct {
	Type       string `json:"TYPE"`
	M          string `json:"M"`
	FromSymbol string `json:"FSYM"`
	ToSymbol   string `json:"TSYM"`
	//The side is 0 for BID and 1 for ASK
	Side int `json:"SIDE"`
	// 1 for ADD (add this position to your ordebook),
	// 2 for REMOVE (take this position out of your orderbook, REMOVE orders also have a quantity of 0),
	// 3 for NOACTION (you should not see these messages they represent updates that we receive from the exchange but have no impact on the orderbook) and
	// 4 for CHANGE/UPDATE (update the available quantity for this position)
	Action int `json:"ACTION"`
	CcSeq  int `json:"CCSEQ"`
	// FIXME: float64 is BAD for price, we need decimal calculation
	Price      float64 `json:"P"`
	Amount     float64 `json:"Q"`
	Seq        int     `json:"SEQ"`
	ReportedNs int64   `json:"REPORTEDNS"`
	DelayNs    int64   `json:"DELAYNS"`
}
