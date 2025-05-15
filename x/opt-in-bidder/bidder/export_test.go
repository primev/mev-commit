package optinbidder

import "time"

func SetNowFunc(f func() time.Time) {
	nowFunc = f
}
