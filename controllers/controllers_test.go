package controllers

import (
	. "lampros-support/models"
	"testing"
)

func TestStartUrgentTimer(t *testing.T) {
	tt, err := startUrgentTimer("0", "1", 0)
	if err != nil {
		t.Errorf("Error starting timer: %v", err)
	}
	timerStop := tt.Timer.Stop()
	tt.Ticker.Stop()
	if tt.Gid != "0" {
		t.Errorf("Timer Gid should be 0: got %s", tt.Gid)
	}
	if tt.TaskId != "1" {
		t.Errorf("Timer taskId should be 1: got %s", tt.TaskId)
	}
	if !timerStop {
		t.Errorf("Timer returned false on stop: %t", timerStop)
	}
}

func TestSendTwilioMessage(t *testing.T) {
	twilioRes, err := sendTwilioMessage("+18592402898", "Successful test")
	if err != nil {
		t.Error(err)
	}
	if twilioRes.Status != "queued" {
		t.Errorf("Message not queued!  Status: %s", twilioRes.Status)
	}
}

func TestDeleteFromTimers(t *testing.T) {
	var timerArray []TickerTimer
	t0, err := startUrgentTimer("0", "1", 0)
	if err != nil {
		t.Errorf("Error starting timer 0: %v", err)
	}
	timerArray = append(timerArray, t0)
	t1, err := startUrgentTimer("1", "1", 0)
	if err != nil {
		t.Errorf("Error starting timer 1: %v", err)
	}
	timerArray = append(timerArray, t1)
	t2, err := startUrgentTimer("2", "1", 0)
	if err != nil {
		t.Errorf("Error starting timer 2: %v", err)
	}
	timerArray = append(timerArray, t2)
	timerArray = deleteFromTimers(timerArray, t2)
	for _, timer := range timerArray {
		if timer.Gid == "2" {
			t.Errorf("Found timer with ID: %s.  It should have been deleted", timer.Gid)
		}
	}
	timerArray = deleteFromTimers(timerArray, t0)
	for _, timer := range timerArray {
		if timer.Gid == "0" {
			t.Errorf("Found timer with ID: %s.  It should have been deleted", timer.Gid)
		}
	}
	if len(timerArray) == 1 {
		if timerArray[0].Gid != "1" {
			t.Errorf("Remaining timer should be 1. Got: %s", timerArray[0].Gid)
		}
	} else {
		t.Errorf("Timer array should be length 1.  Length is: %d", len(timerArray))
	}
}
