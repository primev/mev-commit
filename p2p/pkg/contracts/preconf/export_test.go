package preconfcontract

import "time"

var PreConfABI = preconfABI

func SetDefaultWaitTimeout(timeout time.Duration) func() {
	oldTimeout := defaultWaitTimeout
	defaultWaitTimeout = timeout
	return func() {
		defaultWaitTimeout = oldTimeout
	}
}
