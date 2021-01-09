package main

import (
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
			return
		}
		log.Printf("recv: %s", message)
	}

}
