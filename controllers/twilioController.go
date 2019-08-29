package controllers

import (
	"encoding/json"
	"fmt"
	. "lampros-support/models"
	"strconv"
	"time"
)

//Send a message to any valid phone number
func sendTwilioMessage(toNumber, message string) TwilioMessageResponse {
	params := make(map[string]string)
	params["To"] = toNumber
	params["From"] = TwilioNumber
	params["Body"] = message
	respData := postTwilioRequest(params, parseURL(TwilioBase+"/Messages.json"))
	var resp TwilioMessageResponse
	json.Unmarshal(respData, &resp)
	if resp.ErrorCode > 0 {
		fmt.Println("Twilio API Error! Code: " + strconv.Itoa(resp.ErrorCode) + " Message: " + resp.ErrorMessage)
	}
	return resp
}

//This function starts a big timer for one hour (the maximum response time for urgent tickets)
//This starts tickers to fire off the various levels of urgent texts.
//Urgency: 1 = 5 minute ticks, 2 = 1 minute ticks, 0 = 10 minute ticks for 20 minutes (change for testing)
//From the agreement:
//If reported issues are marked as emergency or high-priority when reported,
//the company will provide a response within 1 hour during normal business hours as defined above,
//3 hours if the report is made outside of normal business hours,
//and 6 hours if during the holiday.
func startUrgentTimer(gid, taskId string, urgency int) TickerTimer {
	var timer *time.Timer
	var ticker *time.Ticker
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println("Error with timezone!")
	}
	now := time.Now().In(loc)
	hour := now.Hour()
	fmt.Printf("Hour of the day: %d\n", hour)
	switch urgency {
	case 0:
		timer = time.NewTimer(time.Minute * 20)
		ticker = time.NewTicker(time.Minute * 10)
	case 1:
		if hour < 9 || hour > 17 {
			timer = time.NewTimer(time.Minute * 165)
		} else {
			timer = time.NewTimer(time.Minute * 45)
		}
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
	return channelTimer
}

//Stops a timer.  If it fails to stop we flush the channel in case there is still a thread
func stopTimer(timer TickerTimer) {
	stopTime := timer.Timer.Stop()
	timer.Ticker.Stop()
	fmt.Println("Stop Tick and Time")
	if !stopTime {
		fmt.Println("Not stopped.  Flushing toilet.")
		<-timer.Timer.C
	}
	return
}

//Accepts parameters TickerTimer array and TickerTimer to delete from the array.
//Returns the modified array of TickerTimers
func deleteFromTimers(timers []TickerTimer, timer TickerTimer) []TickerTimer {
	if len(timers) > 1 {
		for i, t := range timers {
			if timer.Gid == t.Gid {
				timers[i] = timers[len(timers)-1]
			}
		}
	}
	return timers[:len(timers)-1]
}
