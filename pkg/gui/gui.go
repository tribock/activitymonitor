package main

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	activitymonitor "github.com/tribock/activitymonitor/pkg/activity_monitor"
)

type ControlPanel struct {
	Window           fyne.Window
	StatsWindow      fyne.Window
	RunningTimeLabel *widget.Label
	IdleTimeLabel    *widget.Label
	ActiveTimeLabel  *widget.Label
	cancel           context.CancelFunc
	monitor          *activitymonitor.ActivityMonitor
	ctx              context.Context
}

func main() {

	// setup Control Panel

	cp, _ := NewControlPanel(context.Background())
	cp.setupControlPanel()
	cp.Window.ShowAndRun()
}

func NewControlPanel(ctx context.Context) (*ControlPanel, error) {
	a := app.New()
	// Placeholder for future control panel setup
	w := a.NewWindow("Activity Monitor")
	statsWindow := a.NewWindow("Stats")

	monitor := activitymonitor.NewActivityMonitor().WithStats()

	ctx, cancel := context.WithCancel(ctx)
	// Override the close behavior for the stats window
	statsWindow.SetCloseIntercept(func() {
		statsWindow.Hide() // Hide instead of close when "X" is clicked
	})
	return &ControlPanel{
		Window:      w,
		StatsWindow: statsWindow,
		cancel:      cancel,
		monitor:     monitor,
		ctx:         ctx,
	}, nil
}

func (cp *ControlPanel) setupControlPanel() error {

	// Create labels that will be updated
	cp.RunningTimeLabel = widget.NewLabel("Running Time: 0s")
	cp.IdleTimeLabel = widget.NewLabel("Idle Time: 0s")
	cp.ActiveTimeLabel = widget.NewLabel("Active Time: 0s")

	cp.StatsWindow.SetContent(container.NewVBox(
		cp.RunningTimeLabel,
		cp.ActiveTimeLabel,
		cp.IdleTimeLabel,
	))

	startButton := widget.NewButton("Start Monitoring", func() {

		go cp.monitor.RunWithCancel(cp.ctx)

	})

	stopButton := widget.NewButton("Stop Monitoring", func() {
		cp.cancel()
	})

	statsButton := widget.NewButton("Show Stats", func() {
		cp.handleWindowContent()
		cp.StatsWindow.Show()
	})

	close := widget.NewButton("Close", func() {
		cp.Window.Close()
		cp.StatsWindow.Close()
	})

	// Set initial content
	content := container.NewVBox(
		widget.NewLabel("Activity Monitor"),
		close,
		statsButton,
		startButton,
		stopButton,
	)

	// Use fyne.Do to update UI on main thread
	fyne.Do(func() {
		cp.Window.SetContent(content)
	})

	return nil

}

func (cp *ControlPanel) handleWindowContent() {

	time.Sleep(time.Second)
	stats := cp.monitor.GetStats()

	// Update labels on main UI thread
	fyne.Do(func() {
		cp.RunningTimeLabel.SetText("Running Time: " + time.Since(stats.StartTime).Round(time.Second).String())
		cp.IdleTimeLabel.SetText("Idle Time: " + stats.IdleTime.Round(time.Second).String())
		cp.ActiveTimeLabel.SetText("Active Time: " + fmt.Sprintf("%v", time.Since(stats.StartTime).Round(time.Second)-stats.IdleTime.Round(time.Second)))
	})

}
