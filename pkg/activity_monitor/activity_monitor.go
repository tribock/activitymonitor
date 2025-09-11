package activitymonitor

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
)

type activityMonitor struct {
	statsEnabled bool
	stats        Stats
	timeOut      time.Duration
}

func NewActivityMonitor() *activityMonitor {
	return &activityMonitor{
		timeOut: time.Minute, // default value of 1 minute
	}
}

type Stats struct {
	StartTime time.Time
	Idle      bool
	IdleSince time.Time
	IdleTime  time.Duration
}

func (m *activityMonitor) KeepOnMoving() {

	// Setup signal handling to intercept Ctrl+C and prevent ^C from being displayed
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Stdout.WriteString("\r") // Clear the ^C character
		if m.statsEnabled {
			showStats(m.stats)
		}
		os.Exit(0)
	}()

	currentX, currentY := robotgo.Location()
	robotgo.MouseSleep = 10 // milliseconds
	for {
		currentX, currentY = m.moveBackAndForth(currentX, currentY)
		time.Sleep(m.timeOut)
	}
}

func showStats(stats Stats) {
	slog.Info("Stats",
		"RunningTime", time.Since(stats.StartTime),
		"IdleTime", stats.IdleTime,
		"ActiveTime", time.Since(stats.StartTime)-stats.IdleTime,
	)
}

func (m *activityMonitor) WithStats() *activityMonitor {
	m.statsEnabled = true

	m.stats = Stats{
		StartTime: time.Now(),
		Idle:      false,
	}

	return m
}

func (m *activityMonitor) WithTimeout(timeout time.Duration) *activityMonitor {
	m.timeOut = timeout
	return m
}

func (m *activityMonitor) moveBackAndForth(startX, startY int) (int, int) {
	currentX, currentY := robotgo.Location()
	if currentX != startX || currentY != startY {
		if m.statsEnabled {
			m.handleActive()
		}
		slog.Debug("Mouse was moved manually - skipping this cycle.")
		return currentX, currentY
	}
	if m.statsEnabled {
		m.handleIdle()
	}
	robotgo.Move(currentX+1, currentY+1)
	robotgo.Move(currentX, currentY)
	afterX, afterY := robotgo.Location()
	if afterX != currentX || afterY != currentY {
		slog.Warn("Mouse did not return to original position")
		return afterX, afterY
	}
	return currentX, currentY
}

func (m *activityMonitor) handleIdle() {
	if !m.stats.Idle {
		m.stats.Idle = true
		m.stats.IdleSince = time.Now()

		return
	}
	if time.Since(m.stats.IdleSince) < m.timeOut*5 {
		return
	}
	m.stats.IdleTime += m.timeOut
}

func (m *activityMonitor) handleActive() {
	if m.stats.Idle {
		m.stats.Idle = false
		m.stats.IdleSince = time.Time{}
	}
}
