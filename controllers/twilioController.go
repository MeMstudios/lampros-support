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

//IDea is to call this function which starts a big timer for one hour (the maximum response time for urgent tickets)
//This would start tickers to fire off the various levels of urgent texts possibly culminating in a call if not stopped
//Timers could move to their own controller, but I figured their function mostly has to do with twilio.
func StartUrgentTimer() *time.Timer {
	//For 45 minutes send a text every five minutes
	bigTimer := time.NewTimer(time.Minute * 6)
	bigTicker := time.NewTicker(time.Minute * 5)
	go func() {
		for t := range bigTicker.C {
			SendTwilioMessage("+18592402898", "You have an urgent support ticket that hasn't been responded to.  Please check your email and respond!")
			fmt.Println("Sent semi-urgent text at:", t)
		}
	}()
	//start a go process so we can continue on and return the timer
	go func() {
		//After 45 minutes start sending texts every minute
		<-bigTimer.C
		bigTicker.Stop()
		littleTimer := time.NewTimer(time.Minute * 14)
		littleTicker := time.NewTicker(time.Minute)
		go func() {
			for t := range littleTicker.C {
				SendTwilioMessage("+18592402898", "You have an urgent support ticket that hasn't been responded to.  PLEASE RESPOND OR YOU WILL BE FINED!")
				fmt.Println("Sent urgent text at:", t)
			}
		}()
		<-littleTimer.C
		littleTicker.Stop()
	}()
	return bigTimer
}

func StopTimer(timer *time.Timer) {
	stopTime := timer.Stop()
	if !stopTime {
		<-timer.C
	}
}
