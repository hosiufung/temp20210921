package main

import (
	"sort"
	"sync"
)

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

const maxOrderStorage = 10

// though the requirement need to store TOP 10 order of each side
// we need to cope with the order removal
type OrderBook struct {
	// map of ccseq to Order

	bidOrder []Order
	askOrder []Order

	rwLock *sync.RWMutex
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		bidOrder: []Order{},
		askOrder: []Order{},
		rwLock:   &sync.RWMutex{},
	}
}

func (ob *OrderBook) HandleOrder(order Order) {
	ob.rwLock.Lock()
	defer ob.rwLock.Unlock()

	if order.Type != "8" {
		return
	}

	//The side is 0 for BID and 1 for ASK
	//Side int `json:"SIDE"`
	var target []Order
	var sortFunc func(i, j int) bool
	if order.Side == 0 {
		target = ob.bidOrder
		// The bid price refers to the highest price a buyer will pay for a security.
		// thus sort decending
		sortFunc = func(i, j int) bool { return target[i].Price > target[j].Price }
	} else {
		target = ob.askOrder
		// The ask price refers to the lowest price a seller will accept for a security.
		// thus sort ascending
		sortFunc = func(i, j int) bool { return target[i].Price < target[j].Price }
	}

	switch order.Action {
	case 1:
		target = append(target, order)
	case 2:
		for i := range target {
			if target[i].CcSeq == order.CcSeq {
				target = append(target[:i], target[1+i:]...)
				break
			}
		}
	case 3:
	// NOACTION
	case 4:
		for i := range target {
			if target[i].CcSeq == order.CcSeq {
				target[i] = order
				break
			}
		}
	}
	sort.Slice(target, sortFunc)
	if len(target) > maxOrderStorage {
		target = target[:maxOrderStorage]
	}

	if order.Side == 0 {
		ob.bidOrder = target
	} else {
		ob.askOrder = target
	}
}

func (ob *OrderBook) DisplayMidPrice() (midPrice float64, ok bool) {
	ob.rwLock.RLock()
	defer ob.rwLock.RUnlock()

	if len(ob.bidOrder) == 0 || len(ob.askOrder) == 0 {
		return 0, false
	}

	return (ob.bidOrder[0].Price + ob.askOrder[0].Price) / 2, true
}
