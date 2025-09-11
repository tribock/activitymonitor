# ActivityMonitor

ActivityMonitor is a Go program focused on monitoring user activity and collecting statistics about computer usage. It helps keep your screen active to avoid session timeouts or screen lock by periodically moving the mouse cursor by 1 pixel and then moving it back. If manual mouse movement is detected, the program records this as user activity and skips the automated movement for that cycle. This ensures accurate activity tracking while preventing unwanted sleep or lock events.

## Features
- Monitors and collects statistics about user activity (mouse movement detection)
- Prevents session timeouts and screen lock by keeping the screen active
- Periodically moves the mouse cursor by 1px and returns it to the original position if no manual activity is detected
- Skips automated movement if manual mouse movement is detected, ensuring accurate activity stats
- Logs if the mouse does not return to its original position
- Runs continuously with a 1-minute interval between checks

## Requirements
- Go 1.18 or newer
- [robotgo](https://github.com/go-vgo/robotgo#requirements) package requirements
- On macOS: Accessibility permissions for your terminal (System Settings → Privacy & Security → Accessibility)

## Build

```
make build
```

## Run

**Default behavior (with stats and 1-minute timeout):**
```
make run
```

Or directly:
```
go run main.go
```

**With CLI options:**
```
# Run with stats enabled
./activitymonitor --stats

# Run with custom timeout (30 seconds)
./activitymonitor --timeout 30

# Run with both stats and custom timeout
./activitymonitor --stats --timeout 45

# Show help
./activitymonitor --help
```

## CLI Options

- `-s, --stats`: Enable statistics collection and reporting about user activity
- `-t, --timeout <seconds>`: Set the interval between activity checks in seconds
- `-h, --help`: Show help information

**Note:** When no flags are provided, the program runs with stats enabled and a 1-minute timeout (default behavior).

## Install
To install to your $GOPATH/bin directory:

```
make install
```

## Clean

```
make clean
```

## Troubleshooting
- If you see log messages about the mouse not moving, ensure you have granted Accessibility permissions to your terminal.
- If you encounter build errors about missing C headers, make sure you have Xcode Command Line Tools and required libraries (`libpng`, `zlib`) installed.

## License
Apache License - Version 2.0, January 2004
http://www.apache.org/licenses/
