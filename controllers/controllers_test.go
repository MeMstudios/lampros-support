package controllers

import "testing"

func TestStartUrgentTimer(t *testing.T) {
	tt := startUrgentTimer("0", "1", 0)
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
