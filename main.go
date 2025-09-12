package main

import (
	"github.com/tribock/activitymonitor/pkg/cmd"
	"golang.org/x/exp/slog"
)

var (
	version = "latest"
	date    = "unknown"
)

func main() {

	slog.Debug("Build-Info", slog.String("date", date), slog.String("version", version))
	cmd.Execute()
}
