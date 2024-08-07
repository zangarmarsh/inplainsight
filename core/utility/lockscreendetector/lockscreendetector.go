package lockscreendetector

import (
	"log"
	"time"
)

func Analyze(killSignal *chan bool) *chan bool {
	detectorChan := make(chan bool)

	go func() {
		for {
			select {
			case <-*killSignal:
				return
			default:
				log.Printf("Sleep")
				time.Sleep(250 * time.Millisecond)
				if isScreenLocked() {
					detectorChan <- true
				}
			}
		}
	}()

	return &detectorChan
}
