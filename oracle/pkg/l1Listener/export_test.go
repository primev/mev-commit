package l1Listener

import "time"

func SetCheckInterval(interval time.Duration) func() {
	oldInterval := checkInterval
	checkInterval = interval
	return func() {
		checkInterval = oldInterval
	}
}
