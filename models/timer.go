package models

import "time"

type ChanTimer struct {
	Gid    int
	TaskId int
	Timer  *time.Timer
	Ticker *time.Ticker
}
