# MouseMover

MouseMover is a Go program that keeps your computer awake by periodically moving the mouse cursor by 1 pixel and then moving it back. If the mouse is moved manually, the program detects this and skips the movement for that cycle. This helps prevent sleep or screen lock without interfering with manual mouse use.

## Features
- Periodically moves the mouse cursor by 1px and returns it to the original position
- Detects manual mouse movement and skips movement if detected
- Logs if the mouse does not return to its original position
- Runs continuously with a 1-minute interval between movements

## Requirements
- Go 1.18 or newer
- [robotgo](https://github.com/go-vgo/robotgo) library
- On macOS: Accessibility permissions for your terminal (System Settings → Privacy & Security → Accessibility)

## Build

```
make build
```

## Run

```
make run
```

Or directly:

```
go run main.go
```

## Clean

```
make clean
```

## Troubleshooting
- If you see log messages about the mouse not moving, ensure you have granted Accessibility permissions to your terminal.
- If you encounter build errors about missing C headers, make sure you have Xcode Command Line Tools and required libraries (`libpng`, `zlib`) installed.

## License
MIT
