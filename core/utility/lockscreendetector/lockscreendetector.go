package lockscreendetector

import (
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
				time.Sleep(10 * time.Millisecond)
				if isScreenLocked() {
					detectorChan <- true
				}
			}
		}
	}()

	return &detectorChan
}
