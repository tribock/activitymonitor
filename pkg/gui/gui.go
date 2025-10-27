package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	activitymonitor "github.com/tribock/activitymonitor/pkg/activity_monitor"
)

func main() {
	a := app.New()
	w := a.NewWindow("Activity Monitor")
	w.Resize(fyne.NewSize(300, 200))

	go handleWindowContent(w)

	w.ShowAndRun()
}

func handleWindowContent(w fyne.Window) {
	monitor := activitymonitor.NewActivityMonitor().WithStats().WithoutMoving()

	// Start the activity monitor in a separate goroutine
	go monitor.KeepOnMoving()

	// Create labels that will be updated
	runningTimeLabel := widget.NewLabel("Running Time: 0s")
	idleTimeLabel := widget.NewLabel("Idle Time: 0s")
	activeTimeLabel := widget.NewLabel("Active Time: 0s")

	// Set initial content
	content := container.NewVBox(
		widget.NewLabel("Activity Monitor"),
		runningTimeLabel,
		idleTimeLabel,
		activeTimeLabel,
	)

	// Use fyne.Do to update UI on main thread
	fyne.Do(func() {
		w.SetContent(content)
	})

	// Update stats every second
	for {
		time.Sleep(time.Second)
		stats := monitor.GetStats()

		// Update labels on main UI thread
		fyne.Do(func() {
			runningTimeLabel.SetText("Running Time: " + time.Since(stats.StartTime).Round(time.Second).String())
			idleTimeLabel.SetText("Idle Time: " + stats.IdleTime.Round(time.Second).String())
			activeTimeLabel.SetText("Active Time: " + fmt.Sprintf("%v", time.Since(stats.StartTime).Round(time.Second)-stats.IdleTime.Round(time.Second)))
		})
	}
}
