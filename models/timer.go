package models

import "time"

type TickerTimer struct {
	Gid    int
	TaskId int
	Timer  *time.Timer
	Ticker *time.Ticker
}
