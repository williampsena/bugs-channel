package web

import (
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
)

func BuildRateLimitMiddleware(rateLimit int64) *limiter.Limiter {
	rate := float64(rateLimit) / float64(time.Minute.Seconds())

	lmt := tollbooth.NewLimiter(rate, nil)
	lmt.SetTokenBucketExpirationTTL(time.Minute)
	lmt.SetHeaderEntryExpirationTTL(time.Minute)
	lmt.SetMessage("😥 Wow, so many bugs. 🐜")

	return lmt
}
