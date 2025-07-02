package constants

import "time"

const (
	Backoff1stAttempt = "base.backoff-1-attempt"
	Backoff2ndAttempt = "base.backoff-2-attempt"
	Backoff3rdAttempt = "base.backoff-3-attempt"
)

var RetrySleepMap = map[int]time.Duration{
	1: 1 * time.Second,
	2: 2 * time.Second,
	3: 3 * time.Second,
	4: 5 * time.Second,
	5: 10 * time.Second,
	// fallback for any N > 5
}

const DefaultRetrySleep = 15 * time.Second

const (
	Consumer1stRetryAttemptDelay = 1
	Consumer2ndRetryAttemptDelay = 5
	Consumer3rdRetryAttemptDelay = 10
)

func GetBackoffAttemptTopicName() map[int]string {
	return map[int]string{
		1: Backoff1stAttempt,
		2: Backoff2ndAttempt,
		3: Backoff3rdAttempt,
	}
}

func GetBackoffAttemptDelayMinute() map[string]int {
	return map[string]int{
		Backoff1stAttempt: Consumer1stRetryAttemptDelay,
		Backoff2ndAttempt: Consumer2ndRetryAttemptDelay,
		Backoff3rdAttempt: Consumer3rdRetryAttemptDelay,
	}
}
