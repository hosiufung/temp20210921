package main

import (
	"time"
)

func sleepToNearest15sec() {
	now := time.Now()

	// round down the second to the nearest 00, 15, 30, 45
	temp := now.Second() / 15 * 15

	// the second need to sleep for:
	secondDiff := temp + 15 - now.Second()
	time.Sleep(time.Duration(secondDiff)*time.Second - time.Duration(now.Nanosecond())*time.Nanosecond)
}
