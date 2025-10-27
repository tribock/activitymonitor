package activitymonitor

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
)

type ActivityMonitor struct {
	statsEnabled  bool
	movingEnabled bool
	stats         Stats
	timeOut       time.Duration
}

func NewActivityMonitor() *ActivityMonitor {
	return &ActivityMonitor{
		movingEnabled: true,
		timeOut:       time.Minute, // default value of 1 minute
	}
}

type Stats struct {
	StartTime time.Time
	Idle      bool
	IdleSince time.Time
	IdleTime  time.Duration
}

func (m *ActivityMonitor) KeepOnMoving() {

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

	m.Run()

}

// Run the activity monitor loop
// Blocking and not catching interrupts
// Use KeepOnMoving to handle interrupts
func (m *ActivityMonitor) Run() {
	currentX, currentY := robotgo.Location()
	robotgo.MouseSleep = 10 // milliseconds
	for {
		currentX, currentY = m.moveBackAndForth(currentX, currentY)
		time.Sleep(m.timeOut)
	}
}

func (m *ActivityMonitor) RunWithCancel(ctx context.Context) {
	currentX, currentY := robotgo.Location()
	robotgo.MouseSleep = 10 // milliseconds
	for {
		select {
		case <-ctx.Done():
			return
		default:
			currentX, currentY = m.moveBackAndForth(currentX, currentY)
			time.Sleep(m.timeOut)
		}
	}
}

func (m *ActivityMonitor) GetStats() Stats {
	return m.stats
}

func showStats(stats Stats) {
	slog.Info("Stats",
		"RunningTime", time.Since(stats.StartTime),
		"IdleTime", stats.IdleTime,
		"ActiveTime", time.Since(stats.StartTime)-stats.IdleTime,
	)
}

func (m *ActivityMonitor) WithStats() *ActivityMonitor {
	m.statsEnabled = true

	m.stats = Stats{
		StartTime: time.Now(),
		Idle:      false,
	}

	return m
}

func (m *ActivityMonitor) WithoutMoving() *ActivityMonitor {
	m.movingEnabled = false

	return m
}

func (m *ActivityMonitor) WithTimeout(timeout time.Duration) *ActivityMonitor {
	m.timeOut = timeout
	return m
}

func (m *ActivityMonitor) moveBackAndForth(startX, startY int) (int, int) {
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
	if !m.movingEnabled {
		return currentX, currentY
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

func (m *ActivityMonitor) handleIdle() {
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

func (m *ActivityMonitor) handleActive() {
	if m.stats.Idle {
		m.stats.Idle = false
		m.stats.IdleSince = time.Time{}
	}
}
