package models

import "time"

//TickerTimer combines Go's Timer and Ticker and holds a reference to an Asana task id.
type TickerTimer struct {
	Gid    string
	TaskId string
	Timer  *time.Timer
	Ticker *time.Ticker
}
