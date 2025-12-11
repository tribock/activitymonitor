package gui

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
	StateLabel       *widget.Label
	RunningTimeLabel *widget.Label
	IdleTimeLabel    *widget.Label
	ActiveTimeLabel  *widget.Label
	cancel           context.CancelFunc
	monitor          *activitymonitor.ActivityMonitor
	ctx              context.Context
	isRunning        bool
}

func Gui() {

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

	ctx, cancel := context.WithCancel(ctx)
	// Override the close behavior for the stats window
	statsWindow.SetCloseIntercept(func() {
		statsWindow.Hide() // Hide instead of close when "X" is clicked
	})

	// Override the close behavior for the main window
	w.SetCloseIntercept(func() {
		if cancel != nil {
			cancel() // Cancel the context to stop all goroutines
		}
		a.Quit() // Force quit the application
	})
	return &ControlPanel{
		Window:      w,
		StatsWindow: statsWindow,
		cancel:      cancel,
		ctx:         ctx,
	}, nil
}

func (cp *ControlPanel) setupControlPanel() error {

	// Create labels that will be updated
	cp.RunningTimeLabel = widget.NewLabel("Running Time: Not running")
	cp.IdleTimeLabel = widget.NewLabel("Idle Time: Not running")
	cp.ActiveTimeLabel = widget.NewLabel("Active Time: Not running")
	cp.StateLabel = widget.NewLabel("State: Not running")
	cp.StateLabel.Importance = widget.MediumImportance

	cp.StatsWindow.SetContent(container.NewVBox(
		cp.StateLabel,
		cp.RunningTimeLabel,
		cp.ActiveTimeLabel,
		cp.IdleTimeLabel,
	))
	cp.handleWindowContent()

	startButton := widget.NewButton("Start Monitoring", func() {
		if !cp.isRunning {
			go cp.startStatsUpdater()
			// Create a new context each time we start
			cp.monitor = activitymonitor.NewActivityMonitor().WithStats()
			cp.ctx, cp.cancel = context.WithCancel(context.Background())
			cp.isRunning = true
			go cp.monitor.RunWithCancel(cp.ctx)
		}
	})

	stopButton := widget.NewButton("Stop Monitoring", func() {
		if cp.isRunning {
			cp.cancel()
			cp.isRunning = false
			cp.monitor = nil
			cp.handleWindowContent()
		}
	})

	statsButton := widget.NewButton("Show Stats", func() {
		cp.StatsWindow.Show()
	})

	close := widget.NewButton("Close", func() {
		cp.cancel()
		cp.Window.Close()
		cp.StatsWindow.Close()
	})

	// Set initial content
	content := container.NewVBox(
		cp.StateLabel,
		statsButton,
		startButton,
		stopButton,
		close,
	)

	// Use fyne.Do to update UI on main thread
	fyne.Do(func() {
		cp.Window.SetContent(content)
	})
	return nil

}

func (cp *ControlPanel) startStatsUpdater() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cp.handleWindowContent()
		case <-cp.ctx.Done():
			return
		}
	}
}

func (cp *ControlPanel) handleWindowContent() {
	if !cp.isRunning {
		fyne.Do(func() {
			cp.StateLabel.Importance = widget.MediumImportance
			cp.StateLabel.SetText("State: Not running")

		})
		return
	}

	stats := cp.monitor.GetStats()

	// Update labels on main UI thread
	fyne.Do(func() {
		cp.StateLabel.Importance = widget.HighImportance // Use Importance for emphasis
		cp.StateLabel.SetText("State: Running")
		cp.RunningTimeLabel.SetText("Running Time: " + time.Since(stats.StartTime).Round(time.Second).String())
		cp.IdleTimeLabel.SetText("Idle Time: " + stats.IdleTime.Round(time.Second).String())
		cp.ActiveTimeLabel.SetText("Active Time: " + fmt.Sprintf("%v", time.Since(stats.StartTime).Round(time.Second)-stats.IdleTime.Round(time.Second)))
	})

}
