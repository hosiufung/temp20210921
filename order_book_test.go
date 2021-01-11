package main

import (
	"log"

	"testing"
)

func TestHandleOrderAdd(t *testing.T) {
	ob := NewOrderBook()

	// add 4 orders
	order1 := Order{
		Type:       "8",
		FromSymbol: "BTC",
		ToSymbol:   "USD",
		Side:       0,
		Action:     1,
		Price:      30000,
		Amount:     100,
	}
	order2 := Order{
		Type:       "8",
		FromSymbol: "BTC",
		ToSymbol:   "USD",
		Side:       0,
		Action:     1,
		Price:      31000,
		Amount:     100,
	}

	order3 := Order{
		Type:       "8",
		FromSymbol: "BTC",
		ToSymbol:   "USD",
		Side:       1,
		Action:     1,
		Price:      32000,
		Amount:     100,
	}
	order4 := Order{
		Type:       "8",
		FromSymbol: "BTC",
		ToSymbol:   "USD",
		Side:       1,
		Action:     1,
		Price:      33000,
		Amount:     100,
	}

	ob.HandleOrder(order1)
	ob.HandleOrder(order2)
	ob.HandleOrder(order3)
	ob.HandleOrder(order4)

	midPrice, ok := ob.DisplayMidPrice()
	if !ok || midPrice != 31500.0 {
		log.Println(`failed on HandleOrder and DisplayMidPrice`, midPrice, ok)
	}

}
