package models

import "time"

type TickerTimer struct {
	Gid    string
	TaskId string
	Timer  *time.Timer
	Ticker *time.Ticker
}
