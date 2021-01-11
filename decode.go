package main

import (
	"encoding/json"
	"sync"

	"log"

	"github.com/gorilla/websocket"
)

func wssReceiver(c *websocket.Conn, orderBook *OrderBook, isClosed *atomicBool, wg *sync.WaitGroup) {
	defer wg.Done()
	for !isClosed.isTrue() {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			//return
		}
		// log.Printf("recv: %s", message)

		order := Order{}
		if err := json.Unmarshal([]byte(message), &order); err == nil {
			orderBook.HandleOrder(order)
		} else {
			log.Println("order, err", order, err)
		}
	}
}

// FIXME: need to care about graceful shutdown
func printMidPrice(orderBook *OrderBook) {
	for {
		sleepToNearest15sec()

		midPrice, ok := orderBook.DisplayMidPrice()
		log.Println(`midPrice, ok`, midPrice, ok)
	}
}
