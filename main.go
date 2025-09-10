package main

import (
	"time"

	activityMonitor "github.com/tribock/activitymonitor/pkg/activity_monitor"
)

func main() {
	mover := activityMonitor.NewActivityMonitor().WithStats().WithTimeout(1 * time.Minute)
	mover.KeepOnMoving()
}
