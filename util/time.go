package util

import (
	"math/rand"
	"time"
)

const DEFAULT_MIN_CACHE_TIME time.Duration = 60 * time.Minute
const DEFAULT_MAX_CACHE_TIME time.Duration = 3 * 60 * time.Minute

func JitteredCacheDuration(minMins, maxMins time.Duration) time.Duration {
	maxJitter := int((maxMins - minMins) / time.Minute)
	randomDurationMinute := time.Duration(rand.Intn(maxJitter))
	return minMins + randomDurationMinute*time.Minute
}
