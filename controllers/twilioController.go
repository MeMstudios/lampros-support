package controllers

import (
	"encoding/json"
	"fmt"
	. "lampros-support/models"
	"strconv"
	"time"
)

func SendTwilioMessage(toNumber, message string) TwilioMessageResponse {
	params := make(map[string]string)
	params["To"] = toNumber
	params["From"] = TwilioNumber
	params["Body"] = message
	respData := postTwilioRequest(params, parseUrl(TwilioBase+"/Messages.json"))
	var resp TwilioMessageResponse
	json.Unmarshal(respData, &resp)
	if resp.ErrorCode > 0 {
		fmt.Println("Twilio API Error! Code: " + strconv.Itoa(resp.ErrorCode) + " Message: " + resp.ErrorMessage)
	}
	return resp
}

//This function starts a big timer for one hour (the maximum response time for urgent tickets)
//This starts tickers to fire off the various levels of urgent texts.
func StartUrgentTimer(gid, taskId, urgency int) TickerTimer {
	var timer *time.Timer
	var ticker *time.Ticker
	switch urgency {
	case 0:
		timer = time.NewTimer(time.Minute * 20)
		ticker = time.NewTicker(time.Minute * 10)
	case 1:
		timer = time.NewTimer(time.Minute * 45)
		ticker = time.NewTicker(time.Minute * 5)
	case 2:
		timer = time.NewTimer(time.Minute * 14)
		ticker = time.NewTicker(time.Minute)
	}

	//This is how I send a message to a channel associated with the timer.
	var channelTimer TickerTimer
	channelTimer.Gid = gid
	channelTimer.TaskId = taskId
	channelTimer.Timer = timer
	channelTimer.Ticker = ticker
	// go func() {
	// 	for t := range ticker.C {
	// 		SendTwilioMessage("+18592402898", "You have an urgent support ticket that hasn't been responded to.  Please check your email and respond!")
	// 		fmt.Println("Sent semi-urgent text at:", t)
	// 	}
	// }()
	// //start a go process so we can continue on and return the timer
	// go func() {
	// 	//After 45 minutes start sending texts every minute
	// 	<-timer.C
	// 	ticker.Stop()
	// 	fmt.Println("big timer stopped.")
	// 	//If the stop channel has not been set to true it means the timer ran out on its own
	// 	if !<-*channelTimer.Channel {
	// 		fmt.Println("Starting short timer.")
	// 		shortTimer := time.NewTimer(time.Minute * 14)
	// 		shortTicker := time.NewTicker(time.Minute)
	// 		go func() {
	// 			for t := range shortTicker.C {
	// 				SendTwilioMessage("+18592402898", "You have an urgent support ticket that hasn't been responded to.  PLEASE RESPOND OR YOU WILL BE FINED!")
	// 				fmt.Println("Sent urgent text at:", t)
	// 			}
	// 		}()
	// 		<-shortTimer.C
	// 		shortTicker.Stop()
	// 	}
	// }()
	return channelTimer
}

func StopTimer(timer TickerTimer) {
	stopTime := timer.Timer.Stop()
	timer.Ticker.Stop()
	fmt.Println("Stop Tick and Time")
	if !stopTime {
		fmt.Println("Not stopped.  Flushing toilet.")
		<-timer.Timer.C
	}
	return
}

func DeleteFromTimers(timers []TickerTimer, timer TickerTimer) []TickerTimer {
	if len(timers) > 1 {
		for i, t := range timers {
			if timer.Gid == t.Gid {
				timers[i] = timers[len(timers)-1]
			}
		}
	}
	return timers[:len(timers)-1]
}
