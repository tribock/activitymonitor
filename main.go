package main

import (
	"log"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	keepOnMoving()
}

func keepOnMoving() {
	currentX, currentY := robotgo.Location()
	robotgo.MouseSleep = 10 // milliseconds
	for {
		currentX, currentY = moveBackAndForth(currentX, currentY)
		time.Sleep(1 * time.Minute)
	}
}

func moveBackAndForth(startX, startY int) (int, int) {
	currentX, currentY := robotgo.Location()
	if currentX != startX || currentY != startY {
		log.Println("Mouse was moved manually. skipping this cycle.")
		return currentX, currentY
	}
	robotgo.Move(currentX+1, currentY+1)
	robotgo.Move(currentX, currentY)
	afterX, afterY := robotgo.Location()
	if afterX != currentX || afterY != currentY {
		log.Println("Mouse did not return to original position")
		return afterX, afterY
	}
	return currentX, currentY
}
