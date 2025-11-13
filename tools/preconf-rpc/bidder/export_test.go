package bidder

import "time"

func SetNowFunc(f func() time.Time) {
	nowFunc = f
}
